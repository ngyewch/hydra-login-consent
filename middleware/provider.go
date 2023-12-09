package middleware

import "github.com/ory/client-go"

type Provider interface {
	Validate(email string, password string) (bool, error)
	PopulateClaims(consentRequest *client.OAuth2ConsentRequest, idToken map[string]interface{}) error
}
