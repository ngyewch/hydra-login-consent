package basic

import (
	"github.com/ory/client-go"
	"net/http"
)

type LoginValidator func(email string, password string) (bool, error)

type Handler struct {
	loginValidator LoginValidator
}

func NewHandler(loginValidator LoginValidator) *Handler {
	return &Handler{
		loginValidator: loginValidator,
	}
}

func (h *Handler) HandleLogin(r *http.Request) (string, error) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	success, err := h.loginValidator(email, password)
	if err != nil {
		return "", err
	}
	if !success {
		return "", nil
	}
	return email, nil
}

func (h *Handler) PopulateClaims(consentRequest *client.OAuth2ConsentRequest, idToken map[string]interface{}) error {
	for _, scope := range consentRequest.RequestedScope {
		switch scope {
		case "email":
			idToken["email"] = consentRequest.Subject
			idToken["email_verified"] = true
		}
	}
	return nil
}
