package basic

import (
	ory "github.com/ory/client-go"
	"html/template"
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
	Error            string
	ErrorDescription string
	ErrorHint        string
	ErrorDebug       string
}

type ErrorTemplateData struct {
	StatusCode int
	Error      error
}
