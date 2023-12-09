package basic

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fastbill/go-httperrors"
	"github.com/gorilla/csrf"
	"github.com/ngyewch/hydra-login-consent/adaptor/basic/templates"
	ory "github.com/ory/client-go"
	"html/template"
	"net/http"
)

type Renderer struct {
	config    *Config
	templates *template.Template
}

func NewRenderer(config *Config, templates *template.Template) *Renderer {
	return &Renderer{
		config:    config,
		templates: templates,
	}
}

func DefaultTemplates() (*template.Template, error) {
	return template.ParseFS(templates.TemplateFS, "*.gohtml")
}

func (renderer *Renderer) RenderLoginPage(w http.ResponseWriter, r *http.Request, request *ory.OAuth2LoginRequest, errorMessage string) error {
	return renderer.renderPage(w, "login.gohtml",
		LoginPageTemplateData{
			Config:            renderer.config,
			Request:           request,
			CSRFToken:         csrf.Token(r),
			CSRFTemplateField: csrf.TemplateField(r),
			ErrorMessage:      errorMessage,
		},
	)
}

func (renderer *Renderer) RenderConsentPage(w http.ResponseWriter, r *http.Request, request *ory.OAuth2ConsentRequest) error {
	return renderer.renderPage(w, "consent.gohtml",
		ConsentPageTemplateData{
			Config:            renderer.config,
			Request:           request,
			CSRFToken:         csrf.Token(r),
			CSRFTemplateField: csrf.TemplateField(r),
		},
	)
}

func (renderer *Renderer) RenderLogoutPage(w http.ResponseWriter, r *http.Request, request *ory.OAuth2LogoutRequest) error {
	return renderer.renderPage(w, "logout.gohtml",
		LogoutPageTemplateData{
			Config:            renderer.config,
			Request:           request,
			CSRFToken:         csrf.Token(r),
			CSRFTemplateField: csrf.TemplateField(r),
		},
	)
}

func (renderer *Renderer) RenderErrorPage(w http.ResponseWriter, errorName string, errorDescription string, errorHint string, errorDebug string) error {
	return renderer.renderPage(w, "error.gohtml",
		ErrorPageTemplateData{
			Error:            errorName,
			ErrorDescription: errorDescription,
			ErrorHint:        errorHint,
			ErrorDebug:       errorDebug,
		},
	)
}

func (renderer *Renderer) RenderError(w http.ResponseWriter, err error) error {
	var httpError *httperrors.HTTPError
	if errors.As(err, &httpError) {
		renderer.renderHttpError(w, httpError.StatusCode, fmt.Errorf("%s", httpError.Message))
	} else {
		renderer.renderHttpError(w, http.StatusInternalServerError, err)
	}
	return nil
}

func (renderer *Renderer) renderHttpError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	_ = renderer.renderPage(w, "err.gohtml", ErrorTemplateData{
		StatusCode: statusCode,
		Error:      err,
	})
}

func (renderer *Renderer) renderPage(w http.ResponseWriter, templateName string, templateData any) error {
	buf := bytes.NewBuffer(nil)
	err := renderer.templates.Lookup(templateName).
		Execute(buf, templateData)
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
