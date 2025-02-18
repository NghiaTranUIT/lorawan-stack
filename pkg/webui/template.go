// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
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

package webui

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"

	echo "github.com/labstack/echo/v4"
	"go.thethings.network/lorawan-stack/v3/pkg/experimental"
)

// Data contains data to render templates.
type Data struct {
	TemplateData
	AppConfig            interface{}
	ExperimentalFeatures map[string]bool
	PageData             interface{}
	CSPNonce             string
}

// TemplateData contains data to use in the App template.
type TemplateData struct {
	SiteName        string   `name:"site-name" description:"The site name"`
	Title           string   `name:"title" description:"The page title"`
	SubTitle        string   `name:"sub-title" description:"The page sub-title"`
	Description     string   `name:"descriptions" description:"The page description"`
	Language        string   `name:"language" description:"The page language"`
	ThemeColor      string   `name:"theme-color" description:"The page theme color"`
	CanonicalURL    string   `name:"canonical-url" description:"The page canonical URL"`
	AssetsBaseURL   string   `name:"assets-base-url" description:"The base URL to the page assets"`
	BrandingBaseURL string   `name:"branding-base-url" description:"The base URL to the branding assets"`
	IconPrefix      string   `name:"icon-prefix" description:"The prefix to put before the page icons (favicon.ico, touch-icon.png, og-image.png)"`
	CSSFiles        []string `name:"css-file" description:"The names of the CSS files"`
	JSFiles         []string `name:"js-file" description:"The names of the JS files"`
	SentryDSN       string   `name:"sentry-dsn" description:"The Sentry DSN"`
	CSRFToken       string   `name:"-"`
}

// MountPath derives the mount path from the canonical URL of the config.
func (t TemplateData) MountPath() string {
	if url, err := url.Parse(t.CanonicalURL); err == nil {
		if url.Path == "" {
			return "/"
		}
		return url.Path
	}
	return ""
}

const appHTML = `
{{- $assetsBaseURL := .AssetsBaseURL -}}
{{- $brandingBaseURL := or .BrandingBaseURL .AssetsBaseURL -}}
{{- $cspNonce := .CSPNonce -}}
<!doctype html>
<html lang="{{with .Language}}{{.}}{{else}}en{{end}}">
  <head>
    <title>{{.SiteName}}{{with .Title}} {{.}}{{end}}</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, minimum-scale=1">
    <meta name="theme-color" content="{{with .ThemeColor}}{{.}}{{else}}#0D83D0{{end}}">
    <meta http-equiv="X-UA-Compatible" content="IE=edge" >
    {{with .Description}}<meta name="description" content="{{.}}">{{end}}
    <meta property="og:url" content="{{.CanonicalURL}}">
    <meta property="og:site_name" content="{{.SiteName}}{{with .Title}} {{.}}{{end}}">
    {{with .SubTitle}}<meta property="og:title" content="{{.}}">{{end}}
    {{with .Description}}<meta property="og:description" content="{{.}}">{{end}}
    <meta property="og:image" content="{{$brandingBaseURL}}/{{.IconPrefix}}og-image.png">
    <meta property="og:image:secure_url" content="{{$brandingBaseURL}}/{{.IconPrefix}}og-image.png">
    <meta property="og:image:width" content="1200">
    <meta property="og:image:height" content="630">
    <link rel="alternate icon" href="{{$brandingBaseURL}}/{{.IconPrefix}}favicon.ico">
    <link rel="alternate icon" type="image/png" href="{{$brandingBaseURL}}/{{.IconPrefix}}favicon.png">
    <link rel="icon" type="image/svg+xml" href="{{$brandingBaseURL}}/{{.IconPrefix}}favicon.svg">
    <link rel="apple-touch-icon" sizes="180x180" href="{{$brandingBaseURL}}/{{.IconPrefix}}touch-icon.png">
    {{range .CSSFiles}}<link href="{{$assetsBaseURL}}/{{.}}" rel="stylesheet">{{end}}
  </head>
  <body>
    <div id="app"></div>
		<script nonce="{{$cspNonce}}">
		(function (win) {
			var config = {
				APP_ROOT:{{.MountPath}},
				ASSETS_ROOT:{{$assetsBaseURL}},
				BRANDING_ROOT:{{$brandingBaseURL}},
				APP_CONFIG:{{.AppConfig}},
				EXPERIMENTAL_FEATURES:{{.ExperimentalFeatures}},
				SITE_NAME:{{.SiteName}},
				SITE_TITLE:{{.Title}},
				SITE_SUB_TITLE:{{.SubTitle}},
				SENTRY_DSN:{{.SentryDSN}},
				{{with .CSRFToken}}CSRF_TOKEN:{{.}},{{end}}
				{{with .PageData}}PAGE_DATA:{{.}}{{end}}
			};
			win.__ttn_config__ = config;
			if (win.Cypress && win.__initStackConfig) {
				win.__initStackConfig(config);
			}
		})(window);
    </script>
    {{range .JSFiles}}<script nonce="{{$cspNonce}}" type="text/javascript" src="{{$assetsBaseURL}}/{{.}}"></script>{{end}}
  </body>
</html>
`

// Template for rendering the web UI.
// The context is expected to contain TemplateData as "template_data".
// The "app_config" will be rendered into the environment.
var Template *AppTemplate

func init() {
	appHTML := appHTML
	appHTMLLines := strings.Split(appHTML, "\n")
	for i, line := range appHTMLLines {
		appHTMLLines[i] = strings.TrimSpace(line)
	}
	Template = NewAppTemplate(template.Must(template.New("app").Parse(strings.Join(appHTMLLines, ""))))
}

// AppTemplate wraps the application template for the web UI.
type AppTemplate struct {
	template *template.Template
}

// NewAppTemplate instantiates a new application template for the web UI.
func NewAppTemplate(t *template.Template) *AppTemplate {
	return &AppTemplate{template: t}
}

var hashedFiles = map[string]string{}

// RegisterHashedFile maps filenames to webpack generated hashed filenames
func RegisterHashedFile(original, hashed string) {
	hashedFiles[original] = hashed
}

// Render is the echo.Renderer that renders the web UI.
func (t *AppTemplate) Render(w io.Writer, _ string, pageData interface{}, c echo.Context) error {
	templateData := c.Get("template_data").(TemplateData)
	var cspNonce string
	if CSPFeatureFlag.GetValue(c.Request().Context()) {
		if v, ok := c.Get("csp_nonce").(string); ok {
			cspNonce = v
		}
	}
	cssFiles := make([]string, len(templateData.CSSFiles))
	for i, cssFile := range templateData.CSSFiles {
		if hashedFile, ok := hashedFiles[cssFile]; ok {
			cssFiles[i] = hashedFile
		} else {
			cssFiles[i] = cssFile
		}
	}
	templateData.CSSFiles = cssFiles
	jsFiles := make([]string, len(templateData.JSFiles))
	for i, jsFile := range templateData.JSFiles {
		if hashedFile, ok := hashedFiles[jsFile]; ok {
			jsFiles[i] = hashedFile
		} else {
			jsFiles[i] = jsFile
		}
	}
	templateData.JSFiles = jsFiles
	return t.template.Execute(w, Data{
		TemplateData:         templateData,
		AppConfig:            c.Get("app_config"),
		ExperimentalFeatures: experimental.AllFeatures(c.Request().Context()),
		PageData:             pageData,
		CSPNonce:             cspNonce,
	})
}

// Handler is the echo.HandlerFunc that renders the web UI.
// The context is expected to contain TemplateData as "template_data".
// The "app_config" and "page_data" will be rendered into the environment.
func (t *AppTemplate) Handler(c echo.Context) error {
	buf := new(bytes.Buffer)
	if err := Template.Render(buf, "", c.Get("page_data"), c); err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}

// RenderError implements web.ErrorRenderer.
func (t *AppTemplate) RenderError(c echo.Context, statusCode int, err error) error {
	buf := new(bytes.Buffer)
	if err := Template.Render(buf, "", map[string]interface{}{"error": err}, c); err != nil {
		return err
	}
	return c.HTMLBlob(statusCode, buf.Bytes())
}
