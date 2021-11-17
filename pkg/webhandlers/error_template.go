// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
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

package webhandlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/errors"
)

// Data contains data to render templates.
type Data struct {
	ErrorTitle          string
	ErrorMessage        string
	ErrorID             string
	CorrelationID       string
	BackendErrorDetails string
	Year                int
	IsGenericNotFound   bool
}

const errorHTML = `
<!doctype html>
<html lang="en">
  <head>
    <title>Error - The Things Stack Enterprise</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, minimum-scale=1">
    <meta name="theme-color" content="#0D83D0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge" >
    <meta name="description" content="The Things Stack Error Page">
    <meta property="og:image" content="/assets/console-og-image.png">
    <meta property="og:image:secure_url" content="/assets/console-og-image.png">
    <meta property="og:image:width" content="1200">
    <meta property="og:image:height" content="630">
    <link rel="alternate icon" href="/assets/console-favicon.ico">
    <link rel="alternate icon" type="image/png" href="/assets/console-favicon.png">
    <link rel="icon" type="image/svg+xml" href="/assets/console-favicon.svg">
    <link rel="apple-touch-icon" sizes="180x180" href="/assets/console-touch-icon.png">
    <link href="/assets/error/error.css" rel="stylesheet">
  </head>
  <body>
    <div class="wrapper">
      <header class="tti-header">
        <div class="tti-container__full">
          <div class="img-container">
            <img class="tti-header__logo" src="/assets/error/logo.svg" alt="The Things Stack Logo">
          </div>
        </div>
      </header>
      <div class="flex-wrapper">
        <div class="full-view-error">
          <div class="container">
            <div class="row">
              <h1>
                <span class="icon logo">error_outline</span>
                {{.ErrorTitle}}
              </h1>
              <div class="full-view-error-sub">
                <span>
                  {{.ErrorMessage}}
                </span>
                {{ if not .IsGenericNotFound }}
                <span>
                  If the error persists please contact support or administrator.
                </span>
                <br/>
                <span>
                  We're sorry for the inconvenience.
                </span>
                {{ end }}
              </div>
              <div class="error-actions">
                <a href="https://www.thethingsindustries.com/docs" target="_blank" class="button">
                <span class="logo">contact_support</span>
                Documentation
                </a>
                {{ if not .IsGenericNotFound }}
                <span class="error-actions-message">
                  Please attach technical details below to support inquiries.
                </span>
                {{ end }}
              </div>
              {{ if not .IsGenericNotFound }}
              <hr/>
              <div class="detail-colophon">
                <span>
                  Error ID: <code>{{.ErrorID}}</code>
                </span>
                <span>
                  Correlation ID: <code>{{.CorrelationID}}</code>
                </span>
              </div>
              <hr/>
              <details>
                <summary>
                  Technical details
                </summary>
                <pre>{{.BackendErrorDetails}}</pre>
                <button id="copy-button" class="button action-button" data-clipboard-text="{{.BackendErrorDetails}}">
                  <span class="logo">file_copy</span>
                  Copy to clipboard
                </button>
              </details>
              {{ end }}
            </div>
          </div>
        </div>
      </div>
      <footer class="tti-footer">
        <div class="left">
          <div>
            © {{.Year}}
            <a class="link" href="https://www.thethingsindustries.com/docs">The Things Stack</a>
            <span>
              by
              <a class="link" href="https://www.thethingsnetwork.org">The Things Network</a>
              and
              <a class="link" href="https://www.thethingsindustries.com">The Things Industries</a>
            </span>
          </div>
        </div>
      </footer>
    </div>
  </body>
  <script>
    var button = document.getElementById('copy-button');
    var text = button.getAttribute('data-clipboard-text');
    var icon = button.firstChild;
    button.addEventListener("click", function(e) {
      e.preventDefault();
      navigator.clipboard.writeText(text).then(function() {
        button.innerHTML = '<span class="logo">done</span>Copied to clipboard!';
        setTimeout(() => {
          button.innerHTML = '<span class="logo">file_copy</span>Copy to clipboard!';
        }, 3000);
      }, function(err) {
        console.error('Could not copy text: ', err);
      });
    });
  </script>
</html>
`

// Template for rendering the static error.
var Template = func() *ErrorTemplate {
	return NewErrorTemplate(template.Must(template.New("error").Parse(errorHTML)))
}()

// ErrorTemplate wraps the error template for the static error route.
type ErrorTemplate struct {
	template *template.Template
}

// NewErrorTemplate instantiates a new error template for non-frontend handled routes.
func NewErrorTemplate(t *template.Template) *ErrorTemplate {
	return &ErrorTemplate{template: t}
}

// RenderError implements web.ErrorRenderer.
func (t *ErrorTemplate) RenderError(w http.ResponseWriter, err error, code int) error {
	errMsg, _ := json.MarshalIndent(err, "", " ")
	errorID := "n/a"
	errorCorrelationID := "n/a"
	if ttnErr, ok := errors.From(err); ok {
		errorID = ttnErr.FullName()
		errorCorrelationID = ttnErr.CorrelationID()
	}
	var errorTitle string
	var errorMessage string
	switch code {
	case http.StatusNotFound:
		errorTitle = "Page not found"
		errorMessage = "The page you requested cannot be found."
	case http.StatusUnauthorized:
		errorTitle = "Unauthorized"
		errorMessage = "You are not allowed to perform this action."
	default:
		errorTitle = "Unknown error"
		errorMessage = "An unknown error occurred."
	}
	return t.template.Execute(w, Data{
		ErrorTitle:          errorTitle,
		ErrorMessage:        errorMessage,
		ErrorID:             errorID,
		CorrelationID:       errorCorrelationID,
		BackendErrorDetails: string(errMsg),
		Year:                time.Now().Year(),
		IsGenericNotFound:   strings.Contains(errorID, "route_not_found"),
	})
}
