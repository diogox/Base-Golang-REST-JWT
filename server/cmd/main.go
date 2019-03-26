package main

import (
	"github.com/diogox/Calendoer/server/cmd/app"
	"github.com/diogox/Calendoer/server/pkg/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {

	// Create app
	serverApp := app.NewCli("Calendoer")

	// Run it
	err := serverApp.Run(func(c *cli.Context) error {

		e := echo.New()

		// Middleware
		e.Use(middleware.Logger(),
			middleware.Recover(),
			//middleware.HTTPSRedirect(),
		)

		// Routes
		routes.SetupRoutes(e, serverApp.Opts)

		return e.Start(serverApp.Opts.Port)
	})

	// Something went wrong!
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
