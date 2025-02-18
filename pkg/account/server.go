// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package account

import (
	"context"

	echo "github.com/labstack/echo/v4"
	sess "go.thethings.network/lorawan-stack/v3/pkg/account/session"
	account_store "go.thethings.network/lorawan-stack/v3/pkg/account/store"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	web_errors "go.thethings.network/lorawan-stack/v3/pkg/errors/web"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/oauth"
	"go.thethings.network/lorawan-stack/v3/pkg/ratelimit"
	"go.thethings.network/lorawan-stack/v3/pkg/web"
	"go.thethings.network/lorawan-stack/v3/pkg/web/middleware"
	"go.thethings.network/lorawan-stack/v3/pkg/webui"
)

// Server is the interface for the account app server.
type Server interface {
	web.Registerer

	Login(c echo.Context) error
	CurrentUser(c echo.Context) error
	Logout(c echo.Context) error
}

// Component represents the Component to the Account app.
type Component interface {
	Context() context.Context
	RateLimiter() ratelimit.Interface
}

type server struct {
	c           Component
	config      oauth.Config
	store       account_store.Interface
	session     sess.Session
	generateCSP func(config *oauth.Config, nonce string) string
}

// NewServer returns a new account app on top of the given store.
func NewServer(c *component.Component, store account_store.Interface, config oauth.Config, cspFunc func(config *oauth.Config, nonce string) string) (Server, error) {
	s := &server{
		c:           c,
		config:      config,
		store:       store,
		session:     sess.Session{Store: store},
		generateCSP: cspFunc,
	}

	if s.config.Mount == "" {
		s.config.Mount = s.config.UI.MountPath()
	}

	return s, nil
}

type ctxKeyType struct{}

var ctxKey ctxKeyType

func (s *server) configFromContext(ctx context.Context) *oauth.Config {
	if config, ok := ctx.Value(ctxKey).(*oauth.Config); ok {
		return config
	}
	return &s.config
}

func (s *server) Printf(format string, v ...interface{}) {
	log.FromContext(s.c.Context()).Warnf(format, v...)
}

func (s *server) RegisterRoutes(server *web.Server) {
	csrfMiddleware := middleware.CSRF("_csrf", "/", s.config.CSRFAuthKey)
	root := server.Group(
		s.config.Mount,
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if webui.CSPFeatureFlag.GetValue(c.Request().Context()) {
					nonce := webui.GenerateNonce()
					c.Set("csp_nonce", nonce)
					cspString := s.generateCSP(s.configFromContext(c.Request().Context()), nonce)
					c.Response().Header().Set("Content-Security-Policy", cspString)
				}
				return next(c)
			}
		},
		ratelimit.EchoMiddleware(s.c.RateLimiter(), "http:account"),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				config := s.configFromContext(c.Request().Context())
				c.Set("template_data", config.UI.TemplateData)
				frontendConfig := config.UI.FrontendConfig
				frontendConfig.Language = config.UI.TemplateData.Language
				c.Set("app_config", struct {
					oauth.FrontendConfig
				}{
					FrontendConfig: frontendConfig,
				})
				return next(c)
			}
		},
		web_errors.ErrorMiddleware(map[string]web_errors.ErrorRenderer{
			"text/html": webui.Template,
		}),
		csrfMiddleware,
	)

	api := root.Group("/api")
	api.POST("/auth/login", s.Login)
	api.POST("/auth/token-login", s.TokenLogin)
	api.POST("/auth/logout", s.Logout, s.requireLogin)
	api.GET("/me", s.CurrentUser, s.requireLogin)

	page := root.Group("")
	page.GET("/login", webui.Template.Handler, s.redirectToNext)
	page.GET("/token-login", webui.Template.Handler, s.redirectToNext)
	page.GET("/*", webui.Template.Handler)
}
