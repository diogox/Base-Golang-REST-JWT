package app

import (
	"os"

	"github.com/urfave/cli"
)

func NewCli(name string) *CliApp {
	// Instantiate app
	app := cli.NewApp()
	app.Name = name

	// Configure it
	opts := NewAppOptions()
	opts.AddFlags(app)

	return &CliApp{
		App:  app,
		Opts: opts,
	}
}

type CliApp struct {
	App  *cli.App
	Opts *AppOptions
}

func (sa *CliApp) Run(action func(c *cli.Context) error) error {

	// Set action
	sa.App.Action = action

	return sa.App.Run(os.Args)
}
