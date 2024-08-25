package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/fastbill/go-httperrors"
	"github.com/ory/hydra/v2/oauth2"
	"net/http"
)

func (m *Middleware) RefreshTokenHookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := m.postRefreshTokenHook(w, r)
		if err != nil {
			m.handleError(w, r, err)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func (m *Middleware) postRefreshTokenHook(w http.ResponseWriter, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json; charset=UTF-8" {
		return httperrors.New(http.StatusBadRequest, "invalid Content-Type")
	}

	var request oauth2.RefreshTokenHookRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return err
	}

	jsonBytes, err := json.MarshalIndent(request, "", "  ")
	fmt.Printf("request:\n%s\n", string(jsonBytes))

	return nil
}
