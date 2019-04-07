package main

import (
	"fmt"
	"os"

	"github.com/diogox/Calendoer/server/cmd/app"
	"github.com/diogox/Calendoer/server/pkg/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	echoLog "github.com/neko-neko/echo-logrus/log"
	"github.com/urfave/cli"
)

func main() {

	// Instantiate logger
	logger := echoLog.Logger()
	logger.SetLevel(log.DEBUG)

	// Create app
	serverApp := app.NewCli("Calendoer")

	// Run it
	err := serverApp.Run(func(c *cli.Context) error {

		e := echo.New()
		e.Logger = logger

		// TODO: Cleaner logs
		loggerMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				req := c.Request()
				res := c.Response()

				logger.Debug(fmt.Sprintf("[%s] %s - (%d)", req.Method, c.Path(), res.Status))
				return next(c)
			}
		}

		// Middleware
		e.Use(
			//middleware.Logger(),
			loggerMiddleware,
			middleware.Recover(),
			//middleware.HTTPSRedirect(),
		)

		// Routes
		routes.SetupRoutes(e, serverApp.Opts)

		return e.Start(serverApp.Opts.Port)
	})

	// Something went wrong!
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}
