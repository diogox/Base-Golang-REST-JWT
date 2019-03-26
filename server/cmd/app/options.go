package app

import "github.com/urfave/cli"

const (
	Port      = "PORT"
	JWTSecret = "JWT_SECRET"
)

func NewAppOptions() *AppOptions {
	return &AppOptions{}
}

type AppOptions struct {
	Port      string
	JWTSecret string
}

func (opts *AppOptions) AddFlags(app *cli.App) {

	// Define flags
	flags := []cli.Flag{
		// PORT
		cli.StringFlag{
			Name:        "port",
			Value:       ":8000",
			EnvVar:      Port,
			Destination: &opts.Port,
		},
		// JWT SECRET
		cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "",
			EnvVar:      JWTSecret,
			Destination: &opts.JWTSecret,
		},
	}

	app.Flags = append(app.Flags, flags...)
}
