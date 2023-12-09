package adaptor

import (
	"github.com/ory/client-go"
	"net/http"
)

type Handler interface {
	HandleLogin(r *http.Request) (string, error)
	PopulateClaims(consentRequest *client.OAuth2ConsentRequest, idToken map[string]interface{}) error
}
