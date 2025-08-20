package adaptor

import (
	"net/http"

	ory "github.com/ory/client-go"
)

type Renderer interface {
	RenderLoginPage(w http.ResponseWriter, r *http.Request, request *ory.OAuth2LoginRequest, errorMessage string) error
	RenderConsentPage(w http.ResponseWriter, r *http.Request, request *ory.OAuth2ConsentRequest) error
	RenderLogoutPage(w http.ResponseWriter, r *http.Request, request *ory.OAuth2LogoutRequest) error
	RenderErrorPage(w http.ResponseWriter, errorName string, errorDescription string, errorHint string, errorDebug string) error
	RenderError(w http.ResponseWriter, err error) error
}
