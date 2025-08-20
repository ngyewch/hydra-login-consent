package basic

import (
	"context"
	"net/http"

	"github.com/ory/client-go"
)

type LoginValidator func(ctx context.Context, email string, password string) (bool, error)

type ClaimsPopulator func(ctx context.Context, consentRequest *client.OAuth2ConsentRequest, claims map[string]interface{}) error

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
		success, err := h.loginValidator(r.Context(), email, password)
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

func (h *Handler) PopulateClaims(ctx context.Context, consentRequest *client.OAuth2ConsentRequest, claims map[string]interface{}) error {
	cp := h.claimsPopulator
	if cp == nil {
		cp = DefaultClaimsPopulator
	}
	return cp(ctx, consentRequest, claims)
}

func DefaultClaimsPopulator(ctx context.Context, consentRequest *client.OAuth2ConsentRequest, claims map[string]interface{}) error {
	for _, scope := range consentRequest.RequestedScope {
		switch scope {
		case "email":
			if consentRequest.HasSubject() {
				claims["email"] = consentRequest.GetSubject()
				claims["email_verified"] = true
			}
		}
	}
	return nil
}
