package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/REST-JWT/server/pkg/email"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/prisma-client"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

var client *prisma.Client
var emailClient *email.EmailClient
var jwtSecret []byte
var tokenDurationInMinutes int

type RouteOptions struct {
	JWTSecret              []byte
	TokenDurationInMinutes int
	Email                  string
	EmailHost              string
	EmailPort              int
	EmailUsername          string
	EmailPassword          string
}

func SetupRoutes(e *echo.Echo, opts RouteOptions) {

	// Instantiate Prisma client
	client = prisma.New(nil)

	// Instantiate email client
	emailOpts := email.EmailClientOptions{
		Host:     opts.EmailHost,
		Port:     opts.EmailPort,
		Username: opts.EmailUsername,
		Password: opts.EmailPassword,
	}
	emailClient = email.NewEmailClient(opts.Email, emailOpts)

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
	e.POST("/api/auth/register", register)
	e.POST("/api/auth/login", login)
	e.POST("/api/auth/refresh", refreshToken, requireAuth) // Refreshes the JWT token

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
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, users)
}
