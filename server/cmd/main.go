package main

import (
	"fmt"
	"github.com/diogox/REST-JWT/server/cmd/app"
	"github.com/diogox/REST-JWT/server/pkg/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	echoLog "github.com/neko-neko/echo-logrus/log"
	"github.com/spf13/cobra"
	"os"
)

func main() {

	// Instantiate logger
	logger := echoLog.Logger()
	logger.SetLevel(log.DEBUG)

	// Run app
	app.Cmd.RunE = func(cmd *cobra.Command, args []string) error {

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
		opts := routes.RouteOptions{
			// Server Configs
			JWTSecret: []byte(app.JWTSecret),
			TokenDurationInMinutes: app.JWTDuration,

			// Databases Configs
			PrismaHost: app.PrismaHost,
			RedisHost: app.RedisHost,

			// Email Service Configs
			Email: app.Email,
			EmailHost: app.EmailHost,
			EmailPort: app.EmailPort,
			EmailUsername: app.EmailUsername,
			EmailPassword: app.EmailPassword,
		}
		routes.SetupRoutes(e, opts)

		return e.Start( fmt.Sprintf(":%s", app.Port) )
	}

	if err := app.Cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
