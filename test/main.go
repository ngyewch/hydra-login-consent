package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"runtime/debug"
)

var (
	flagConfigFile = &cli.PathFlag{
		Name:     "config-file",
		Usage:    "config file",
		Required: true,
	}

	app = &cli.App{
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

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
