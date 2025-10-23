package basic

import (
	ory "github.com/ory/client-go"
)

type LoginPageTemplateData struct {
	Config       *Config
	Request      *ory.OAuth2LoginRequest
	ErrorMessage string
}

type ConsentPageTemplateData struct {
	Config  *Config
	Request *ory.OAuth2ConsentRequest
}

type LogoutPageTemplateData struct {
	Config  *Config
	Request *ory.OAuth2LogoutRequest
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
