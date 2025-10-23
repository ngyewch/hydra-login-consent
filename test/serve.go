package main

import (
	"context"
	"net/http"

	"filippo.io/csrf"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/v2"
	"github.com/ngyewch/hydra-login-consent/adaptor/basic"
	"github.com/ngyewch/hydra-login-consent/adaptor/basic/static"
	uiMiddleware "github.com/ngyewch/hydra-login-consent/middleware"
	ory "github.com/ory/client-go"
	"github.com/urfave/cli/v3"
)

type ServeConfig struct {
	ListenAddr        string        `koanf:"listenAddr" validate:"required"`
	HydraAdminApiUrls []string      `koanf:"hydraAdminApiUrls" validate:"gt=0,dive,url"`
	UI                *basic.Config `koanf:"ui" validate:"required"`
	Users             []UserEntry   `koanf:"user" validate:"required,dive"`
}

type UserEntry struct {
	Email    string `koanf:"email" validate:"email"`
	Password string `koanf:"password" validate:"required"`
}

func doServe(ctx context.Context, cmd *cli.Command) error {
	configFile := cmd.String(flagConfigFile.Name)

	k := koanf.New(".")
	err := mergeConfig(k, configFile)
	if err != nil {
		return err
	}

	var config ServeConfig
	err = k.Unmarshal("", &config)
	if err != nil {
		return err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		return err
	}

	configuration := ory.NewConfiguration()
	configuration.Servers = make([]ory.ServerConfiguration, 0)
	for _, hydraAdminApiUrl := range config.HydraAdminApiUrls {
		configuration.Servers = append(configuration.Servers, ory.ServerConfiguration{
			URL: hydraAdminApiUrl,
		})
	}
	oryClient := ory.NewAPIClient(configuration)
	oryContext := context.Background()

	templates, err := basic.DefaultTemplates()
	if err != nil {
		return err
	}
	renderer := basic.NewRenderer(config.UI, templates)
	handler := basic.NewHandler(newLoginValidator(config.Users), nil)

	loginConsentMiddleware := uiMiddleware.New(oryClient, oryContext, renderer, handler)

	csrfMiddleware := csrf.New()

	router := chi.NewRouter()
	router.Use(csrfMiddleware.Handler)
	router.Use(loginConsentMiddleware.Handler)
	router.Mount("/", http.FileServer(http.FS(static.StaticFS)))
	return http.ListenAndServe(config.ListenAddr, router)
}

func newLoginValidator(users []UserEntry) basic.LoginValidator {
	return func(ctx context.Context, email string, password string) (bool, error) {
		for _, user := range users {
			if (user.Email == email) && (user.Password == password) {
				return true, nil
			}
		}
		return false, nil
	}
}
