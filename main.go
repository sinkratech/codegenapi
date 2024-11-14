package main

import (
	"log"
	"os"

	"github.com/sinkratech/codegenapi/actions"
	"github.com/urfave/cli/v2"
)

var version string

func main() {
	app := cli.App{
		Usage:   "Generate boilerplate for Sinkra API",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:    "feature",
				Aliases: []string{"feat", "f"},
				Usage:   "Generate new api feature folder (deps.go and entrypoint.go)",
				Action:  actions.GenFeature,
				Args:    true,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
