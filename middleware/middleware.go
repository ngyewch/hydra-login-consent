package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/fastbill/go-httperrors"
	"github.com/gorilla/mux"
	ory "github.com/ory/client-go"
	"html/template"
	"net/http"
)

type Middleware struct {
	cfg              *Config
	oryClient        *ory.APIClient
	oryAuthedContext context.Context
	templates        *template.Template
	provider         Provider
}

type ErrorTemplateData struct {
	StatusCode int
	Error      error
}

func New(cfg *Config, oryClient *ory.APIClient, oryAuthedContext context.Context, templates *template.Template, provider Provider) *Middleware {
	return &Middleware{
		cfg:              cfg,
		oryClient:        oryClient,
		oryAuthedContext: oryAuthedContext,
		templates:        templates,
		provider:         provider,
	}
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router := mux.NewRouter()

		m.handleGET(router, "/login", m.getLogin)       // TODO make configurable
		m.handlePOST(router, "/login", m.postLogin)     // TODO make configurable
		m.handleGET(router, "/consent", m.getConsent)   // TODO make configurable
		m.handlePOST(router, "/consent", m.postConsent) // TODO make configurable
		m.handleGET(router, "/logout", m.getLogout)     // TODO make configurable
		m.handlePOST(router, "/logout", m.postLogout)   // TODO make configurable
		// TODO error URL
		// TODO post_logout_redirect URL

		var routeMatch mux.RouteMatch
		if router.Match(r, &routeMatch) {
			router.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (m *Middleware) handleGET(router *mux.Router, path string, handler func(w http.ResponseWriter, r *http.Request) error) {
	router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			m.handleError(w, err)
			return
		}
	}).Methods("GET")
}

func (m *Middleware) handlePOST(router *mux.Router, path string, handler func(w http.ResponseWriter, r *http.Request) error) {
	router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			m.handleError(w, err)
			return
		}
	}).Methods("POST")
}

func (m *Middleware) handleError(w http.ResponseWriter, err error) {
	var httpError *httperrors.HTTPError
	if errors.As(err, &httpError) {
		m.handleHttpError(w, httpError.StatusCode, fmt.Errorf("%s", httpError.Message))
	} else {
		m.handleHttpError(w, http.StatusInternalServerError, err)
	}
}

func (m *Middleware) handleHttpError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	_ = m.renderPage(w, "error.gohtml", ErrorTemplateData{
		StatusCode: statusCode,
		Error:      err,
	})
}
