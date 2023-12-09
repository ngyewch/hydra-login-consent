package basic

import (
	"github.com/ory/client-go"
	"net/http"
)

type LoginValidator func(email string, password string) (bool, error)

type ClaimsPopulator func(subject string, scope string, claims map[string]interface{}) error

type Handler struct {
	loginValidator  LoginValidator
	claimsPopulator ClaimsPopulator
}

func NewHandler(loginValidator LoginValidator, claimsPopulator ClaimsPopulator) *Handler {
	return &Handler{
		loginValidator:  loginValidator,
		claimsPopulator: claimsPopulator,
	}
}

func (h *Handler) HandleLogin(r *http.Request) (string, error) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if h.loginValidator != nil {
		success, err := h.loginValidator(email, password)
		if err != nil {
			return "", err
		}
		if !success {
			return "", nil
		}
		return email, nil
	} else {
		return "", nil
	}
}

func (h *Handler) PopulateClaims(consentRequest *client.OAuth2ConsentRequest, claims map[string]interface{}) error {
	cp := h.claimsPopulator
	if cp == nil {
		cp = DefaultClaimsPopulator
	}
	for _, scope := range consentRequest.RequestedScope {
		err := cp(consentRequest.GetSubject(), scope, claims)
		if err != nil {
			return err
		}
	}
	return nil
}

func DefaultClaimsPopulator(subject string, scope string, claims map[string]interface{}) error {
	switch scope {
	case "email":
		if subject != "" {
			claims["email"] = subject
			claims["email_verified"] = true
		}
	}
	return nil
}
