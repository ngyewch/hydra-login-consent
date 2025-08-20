package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v3"
)

var (
	flagConfigFile = &cli.StringFlag{
		Name:  "config-file",
		Usage: "config file",
	}

	app = &cli.Command{
		Name:   "hydra-login-consent-test",
		Usage:  "hydra-login-consent test",
		Action: nil,
		Commands: []*cli.Command{
			{
				Name:   "serve",
				Usage:  "serve",
				Action: doServe,
				Flags: []cli.Flag{
					flagConfigFile,
				},
			},
		},
	}
)

func main() {
	buildInfo, _ := debug.ReadBuildInfo()
	if buildInfo != nil {
		app.Version = buildInfo.Main.Version
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
