package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/ngyewch/hydra-login-consent/adaptor"
	ory "github.com/ory/client-go"
	"net/http"
)

type Middleware struct {
	oryClient        *ory.APIClient
	oryAuthedContext context.Context
	renderer         adaptor.Renderer
	handler          adaptor.Handler
}

func New(oryClient *ory.APIClient, oryAuthedContext context.Context, renderer adaptor.Renderer, handler adaptor.Handler) *Middleware {
	return &Middleware{
		oryClient:        oryClient,
		oryAuthedContext: oryAuthedContext,
		renderer:         renderer,
		handler:          handler,
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
		m.handleGET(router, "/error", m.getError)       // TODO make configurable

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
			_ = m.renderer.RenderError(w, err)
			return
		}
	}).Methods("GET")
}

func (m *Middleware) handlePOST(router *mux.Router, path string, handler func(w http.ResponseWriter, r *http.Request) error) {
	router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			_ = m.renderer.RenderError(w, err)
			return
		}
	}).Methods("POST")
}
