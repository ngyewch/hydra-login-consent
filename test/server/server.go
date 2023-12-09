package server

import (
	"context"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	uiMiddleware "github.com/ngyewch/hydra-login-consent/middleware"
	"github.com/ngyewch/hydra-login-consent/resources"
	ory "github.com/ory/client-go"
	"html/template"
)

type Server struct {
	cfg *Config
	e   *echo.Echo
}

func New(cfg *Config, uiCfg *uiMiddleware.Config, provider uiMiddleware.Provider) (*Server, error) {
	configuration := ory.NewConfiguration()
	configuration.Servers = make([]ory.ServerConfiguration, 0)
	for _, hydraAdminApiUrl := range cfg.HydraAdminApiUrls {
		configuration.Servers = append(configuration.Servers, ory.ServerConfiguration{
			URL: hydraAdminApiUrl,
		})
	}
	oryClient := ory.NewAPIClient(configuration)

	oryContext := context.Background()

	templates, err := template.ParseFS(resources.TemplateFS, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	loginConsentMiddleware := uiMiddleware.New(uiCfg, oryClient, oryContext, templates, provider)

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
