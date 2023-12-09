package cmd

import (
	"github.com/knadh/koanf/v2"
	"github.com/ngyewch/hydra-login-consent/middleware"
	"github.com/ngyewch/hydra-login-consent/test/server"
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
	Serve *server.Config     `koanf:"serve"`
	UI    *middleware.Config `koanf:"ui"`
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

	var serveConfig ServeConfig
	err = k.Unmarshal("", &serveConfig)
	if err != nil {
		return err
	}

	provider := middleware.NewDummyProvider()

	s, err := server.New(serveConfig.Serve, serveConfig.UI, provider)
	if err != nil {
		return err
	}

	return s.Start()
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("config-file", "", "config file")
}
