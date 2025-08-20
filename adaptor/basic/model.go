package basic

import (
	"html/template"

	ory "github.com/ory/client-go"
)

type LoginPageTemplateData struct {
	Config            *Config
	Request           *ory.OAuth2LoginRequest
	CSRFToken         string
	CSRFTemplateField template.HTML
	ErrorMessage      string
}

type ConsentPageTemplateData struct {
	Config            *Config
	Request           *ory.OAuth2ConsentRequest
	CSRFToken         string
	CSRFTemplateField template.HTML
}

type LogoutPageTemplateData struct {
	Config            *Config
	Request           *ory.OAuth2LogoutRequest
	CSRFToken         string
	CSRFTemplateField template.HTML
}

type ErrorPageTemplateData struct {
	Config           *Config
	Error            string
	ErrorDescription string
	ErrorHint        string
	ErrorDebug       string
}

type ErrorTemplateData struct {
	Config     *Config
	StatusCode int
	Error      error
}
