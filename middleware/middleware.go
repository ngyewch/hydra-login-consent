package middleware

import (
	"context"
	"github.com/ngyewch/hydra-login-consent/adaptor"
	ory "github.com/ory/client-go"
	"log/slog"
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

func (m *Middleware) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		renderErr := m.renderer.RenderError(w, err)
		if renderErr != nil {
			slog.LogAttrs(context.Background(), slog.LevelError, "render error",
				slog.Any("err", err),
				slog.Any("renderErr", renderErr),
			)
		}
		return
	}
}
