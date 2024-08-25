package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/knadh/koanf/v2"
	"github.com/ngyewch/hydra-login-consent/adaptor/basic"
	"github.com/ngyewch/hydra-login-consent/adaptor/basic/static"
	uiMiddleware "github.com/ngyewch/hydra-login-consent/middleware"
	ory "github.com/ory/client-go"
	"github.com/urfave/cli/v2"
	"net/http"
)

type ServeConfig struct {
	ListenAddr        string        `koanf:"listenAddr"`
	CsrfAuthKey       string        `koanf:"csrfAuthKey"`
	HydraAdminApiUrls []string      `koanf:"hydraAdminApiUrls"`
	UI                *basic.Config `koanf:"ui"`
	Users             []UserEntry   `koanf:"user"`
}

type UserEntry struct {
	Email    string `koanf:"email"`
	Password string `koanf:"password"`
}

func doServe(cCtx *cli.Context) error {
	configFile := flagConfigFile.Get(cCtx)

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
	csrfMiddleware := csrf.Protect([]byte(config.CsrfAuthKey))

	r := chi.NewRouter()
	r.Route("/frontend", func(r chi.Router) {
		r.Use(csrfMiddleware)
		r.HandleFunc("/login", loginConsentMiddleware.LoginHandler)
		r.HandleFunc("/consent", loginConsentMiddleware.ConsentHandler)
		r.HandleFunc("/logout", loginConsentMiddleware.LogoutHandler)
		r.HandleFunc("/error", loginConsentMiddleware.ErrorHandler)
	})
	r.Route("/backend", func(r chi.Router) {
		r.HandleFunc("/token-hook", loginConsentMiddleware.TokenHookHandler)
		r.HandleFunc("/refresh-token-hook", loginConsentMiddleware.RefreshTokenHookHandler)
	})
	r.Handle("/", http.FileServerFS(static.StaticFS))

	return http.ListenAndServe(config.ListenAddr, r)
}

func newLoginValidator(users []UserEntry) basic.LoginValidator {
	return func(email string, password string) (bool, error) {
		for _, user := range users {
			if (user.Email == email) && (user.Password == password) {
				return true, nil
			}
		}
		return false, nil
	}
}
