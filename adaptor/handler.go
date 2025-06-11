package adaptor

import (
	"context"
	"github.com/ory/client-go"
	"net/http"
)

type Handler interface {
	HandleLogin(r *http.Request) (string, error)
	PopulateClaims(ctx context.Context, consentRequest *client.OAuth2ConsentRequest, claims map[string]interface{}) error
}
