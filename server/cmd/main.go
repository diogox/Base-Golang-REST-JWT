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
	"net/http"
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

		// CORS restricted
		// Allows requests from `localhost` and the specified port.
		// wth GET, PUT, POST or DELETE method.
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			//AllowOrigins: []string{"http://localhost:" + app.Port},
			//AllowOrigins: []string{app.AppUrl},
			// TODO: this is for development only!!
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		}))

		// Routes
		opts := routes.RouteOptions{
			// Server Configs
			AppUrl:                        app.AppUrl,
			JWTSecret:                     []byte(app.JWTSecret),
			AuthTokenDurationInMinutes:    app.JWTAuthDuration,
			RefreshTokenDurationInMinutes: app.JWTRefreshDuration,

			// Databases Configs
			PrismaHost: app.PrismaHost,
			RedisHost:  app.RedisHost,

			// Email Service Configs
			EmailBodyPath: app.EmailBodyPath,
			Email:         app.Email,
			EmailHost:     app.EmailHost,
			EmailPort:     app.EmailPort,
			EmailUsername: app.EmailUsername,
			EmailPassword: app.EmailPassword,

			// Account Configs
			RemoveUnverifiedAccountAfterNDays:    app.RemoveUnverifiedAccountAfterNDays,
			AccountAllowedNOfFailedLoginAttempts: app.AccountAllowedNOfFailedLoginAttempts,
			AccountLockDuration:                  app.AccountLockDuration,
		}
		routes.SetupRoutes(e, opts)

		return e.Start(fmt.Sprintf(":%s", app.Port))
	}

	if err := app.Cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
