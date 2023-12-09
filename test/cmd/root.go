package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   fmt.Sprintf("hydra-login-consent-test [flags]"),
		Short: "hydra-login-consent test app",
		RunE:  help,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// do nothing
}
