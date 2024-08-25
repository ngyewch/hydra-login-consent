package middleware

import (
	"fmt"
	"net/http"
)

func (m *Middleware) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := m.getError(w, r)
		if err != nil {
			m.handleError(w, r, err)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func (m *Middleware) getError(w http.ResponseWriter, r *http.Request) error {
	errorName := r.URL.Query().Get("error")
	errorDescription := r.URL.Query().Get("error_description")
	errorHint := r.URL.Query().Get("error_hint")
	errorDebug := r.URL.Query().Get("error_debug")

	err := m.renderer.RenderErrorPage(w, errorName, errorDescription, errorHint, errorDebug)
	if err != nil {
		return err
	}

	return nil
}
