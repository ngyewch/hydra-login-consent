package middleware

import "github.com/ory/client-go"

type DummyProvider struct {
}

func NewDummyProvider() *DummyProvider {
	return &DummyProvider{}
}

func (provider *DummyProvider) Validate(email string, password string) (bool, error) {
	return password == "password", nil
}

func (provider *DummyProvider) PopulateClaims(consentRequest *client.OAuth2ConsentRequest, idToken map[string]interface{}) error {
	for _, scope := range consentRequest.RequestedScope {
		switch scope {
		case "email":
			idToken["email"] = consentRequest.Subject
			idToken["email_verified"] = true
		}
	}
	return nil
}
