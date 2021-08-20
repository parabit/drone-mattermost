package main

import (
	"os"

	"github.com/drone-plugins/drone-plugin-lib/errors"
	"github.com/joho/godotenv"
	"github.com/kenshaw/drone-mattermost/plugin"
	"github.com/urfave/cli/v2"
)

// version is the app version.
var version = "0.0.0-dev"

func main() {
	// overload environment
	if _, err := os.Stat("/run/drone/env"); err == nil {
		godotenv.Overload("/run/drone/env")
	}
	// create plugin
	p := plugin.New()
	app := &cli.App{
		Name:    "drone-mattermost",
		Usage:   "build notifications for mattermost",
		Version: version,
		Flags:   p.Flags(),
		Action:  p.Run,
	}
	// execute
	if err := app.Run(os.Args); err != nil {
		errors.HandleExit(err)
	}
}
