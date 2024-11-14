package main

import (
	"log"
	"os"

	"github.com/sinkratech/codegen/actions"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Usage: "Generate boilerplate for Sinkra API",
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
