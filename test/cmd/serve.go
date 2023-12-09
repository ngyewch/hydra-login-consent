package cmd

import (
	"context"
	"github.com/gorilla/csrf"
	"github.com/knadh/koanf/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ngyewch/hydra-login-consent/adaptor/basic"
	uiMiddleware "github.com/ngyewch/hydra-login-consent/middleware"
	ory "github.com/ory/client-go"
	"github.com/spf13/cobra"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve [flags]",
		Short: "Serve",
		RunE:  serve,
	}
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

func serve(cmd *cobra.Command, args []string) error {
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return err
	}

	k := koanf.New(".")
	err = mergeConfig(k, configFile)
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
	handler := basic.NewHandler(newLoginValidator(config.Users))

	loginConsentMiddleware := uiMiddleware.New(oryClient, oryContext, renderer, handler)

	e := echo.New()
	e.HideBanner = true

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echo.WrapMiddleware(csrf.Protect([]byte(config.CsrfAuthKey))))
	e.Use(echo.WrapMiddleware(loginConsentMiddleware.Handler))

	return e.Start(config.ListenAddr)
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

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("config-file", "", "config file")
}
