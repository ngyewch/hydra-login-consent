package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/fastbill/go-httperrors"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/oauth2"
	"net/http"
)

func (m *Middleware) TokenHookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := m.postTokenHook(w, r)
		if err != nil {
			m.handleError(w, r, err)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func (m *Middleware) postTokenHook(w http.ResponseWriter, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json; charset=UTF-8" {
		return httperrors.New(http.StatusBadRequest, "invalid Content-Type")
	}

	var request oauth2.TokenHookRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return err
	}

	jsonBytes, err := json.MarshalIndent(request, "", "  ")
	fmt.Printf("request:\n%s\n", string(jsonBytes))

	response := oauth2.TokenHookResponse{
		Session: flow.AcceptOAuth2ConsentRequestSession{
			AccessToken: map[string]interface{}{
				"boo": "hoo",
			},
			IDToken: map[string]interface{}{
				"boohoo": "gjoac",
			},
		},
	}

	jsonBytes, err = json.MarshalIndent(response, "", "  ")
	fmt.Printf("response:\n%s\n", string(jsonBytes))

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = w.Write(responseBytes)
	if err != nil {
		return err
	}

	return nil
}
