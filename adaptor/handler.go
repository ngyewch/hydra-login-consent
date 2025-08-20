package adaptor

import (
	"context"
	"net/http"

	"github.com/ory/client-go"
)

type Handler interface {
	HandleLogin(r *http.Request) (string, error)
	PopulateClaims(ctx context.Context, consentRequest *client.OAuth2ConsentRequest, claims map[string]interface{}) error
}
