package routes

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/Calendoer/server/cmd/app"
	"github.com/diogox/Calendoer/server/pkg/models/auth"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var JWTSecret []byte

func SetupRoutes(e *echo.Echo, opts *app.AppOptions) {

	JWTSecret = []byte(opts.JWTSecret)

	// Serve website
	e.File("/*", "../web/index.html")
	e.File("/favicon.ico", "../web/images/favicon.ico")

	// Login
	e.POST("/login", handleLogin)

	// By default, the key is extracted from the header "Authorization".
	// To get it from a field named `token` in the JSON we could add `TokenLookup: "query:token"` to the JWT Configs
	requireAuth := middleware.JWT(JWTSecret)

	// Doesn't require Auth
	apiEndpoint := e.Group("/api", requireAuth)

	// Requires Auth
	apiEndpoint.GET("/users", handleGetUsers)
}

func handleLogin(c echo.Context) error {
	logger := c.Logger()
	logger.Info("Login request received!")

	var userCreds auth.Credentials

	// Get POST body
	err := c.Bind(&userCreds)
	if err != nil {
		return err
	}

	// TODO: Perform login check here
	/*
		if userLogin.Username != "diogox" && userLogin.Password != "Diogox" {
			// Throws unauthorized error
			return echo.ErrUnauthorized
		}
	*/

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = userCreds.Username
	claims["admin"] = false

	// TODO: Change expiry to 5 minutes, when done with development
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString(JWTSecret)
	if err != nil {
		return err
	}

	// Create response
	res := auth.AuthResponse{
		Token: t,
	}

	return c.JSON(http.StatusOK, res)
}

func handleGetUsers(c echo.Context) error {
	//logger := c.Logger()

	token := c.Get("user").(*jwt.Token)

	// Check if valid
	if !token.Valid {
		return echo.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, name)
}
