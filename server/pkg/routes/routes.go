package routes

import (
	"github.com/diogox/Calendoer/server/cmd/app"
	"github.com/diogox/Calendoer/server/pkg/models/auth"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func SetupRoutes(e *echo.Echo, opts *app.AppOptions) {

	// Serve website
	e.File("/", "web/index.html")

	// Doesn't require Auth
	apiEndpoint := e.Group("/api")
	apiEndpoint.POST("/auth", handleAuth)

	// By default, the key is extracted from the header "Authorization".
	// To get it from a field named `token` in the JSON we could add `TokenLookup: "query:token"` to the JWT Configs
	jwtMiddleware := middleware.JWT([]byte(opts.JWTSecret))

	// Requires Auth
	apiEndpoint.GET("/users", handleUsers, jwtMiddleware)
}

func handleAuth(c echo.Context) error {
	var userLogin auth.Login

	// Get POST body
	err := c.Bind(&userLogin)
	if err != nil {
		return err
	}

	// Create response
	res := auth.AuthResponse{
		Token: "Token",
	}

	return c.JSON(http.StatusOK, res)
}

func handleUsers(c echo.Context) error {
	return nil
}
