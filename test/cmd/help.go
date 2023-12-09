package cmd

import (
	"github.com/spf13/cobra"
)

func help(cmd *cobra.Command, args []string) error {
	err := cmd.Help()
	if err != nil {
		return err
	}
	return nil
}
