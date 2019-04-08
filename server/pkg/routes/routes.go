package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/Calendoer/generated/prisma-client"
	"github.com/diogox/Calendoer/server/pkg/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

var client *prisma.Client
var jwtSecret []byte
var tokenDurationInMinutes int

type RouteOptions struct {
	JWTSecret              []byte
	TokenDurationInMinutes int
}

func SetupRoutes(e *echo.Echo, opts RouteOptions) {

	// Instantiate Prisma client
	client = prisma.New(nil)

	// Set vars from options
	jwtSecret = opts.JWTSecret
	tokenDurationInMinutes = opts.TokenDurationInMinutes

	e.Validator = newValidator()

	// Serve website
	e.File("/*", "../web/index.html")
	e.File("/favicon.ico", "../web/images/favicon.ico")

	// By default, the key is extracted from the header "Authorization".
	// To get it from a field named `token` in the JSON we could add `TokenLookup: "query:token"` to the JWT Configs
	requireAuth := middleware.JWT(jwtSecret)

	// Auth
	e.POST("/register", register)
	e.POST("/login", login)
	e.POST("/refresh", refreshToken, requireAuth) // Refreshes the JWT token

	// Endpoint that requires authentication
	apiEndpoint := e.Group("/api", requireAuth)

	// API endpoints (TODO: Add endpoints here!)
	apiEndpoint.GET("/users", handleGetUsers)
}

func handleGetUsers(c echo.Context) error {
	ctx := c.Request().Context()
	//logger := c.Logger()

	token := c.Get("user").(*jwt.Token)

	// Check if valid
	if !token.Valid {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Invalid Token!",
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	users, err := client.Users(&prisma.UsersParams{
		Where: &prisma.UserWhereInput{
			Username: &username,
		},
	}).Exec(ctx)
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, users)
}
