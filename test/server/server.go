package server

import (
	"context"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ngyewch/hydra-login-consent/adaptor"
	uiMiddleware "github.com/ngyewch/hydra-login-consent/middleware"
	ory "github.com/ory/client-go"
)

type Server struct {
	cfg *Config
	e   *echo.Echo
}

func New(cfg *Config, renderer adaptor.Renderer, handler adaptor.Handler) (*Server, error) {
	configuration := ory.NewConfiguration()
	configuration.Servers = make([]ory.ServerConfiguration, 0)
	for _, hydraAdminApiUrl := range cfg.HydraAdminApiUrls {
		configuration.Servers = append(configuration.Servers, ory.ServerConfiguration{
			URL: hydraAdminApiUrl,
		})
	}
	oryClient := ory.NewAPIClient(configuration)

	oryContext := context.Background()

	loginConsentMiddleware := uiMiddleware.New(oryClient, oryContext, renderer, handler)

	e := echo.New()
	e.HideBanner = true

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echo.WrapMiddleware(csrf.Protect([]byte(cfg.CsrfAuthKey))))
	e.Use(echo.WrapMiddleware(loginConsentMiddleware.Handler))

	return &Server{
		cfg: cfg,
		e:   e,
	}, nil
}

func (server *Server) Start() error {
	return server.e.Start(server.cfg.ListenAddr)
}
