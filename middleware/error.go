package middleware

import (
	"net/http"
)

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
