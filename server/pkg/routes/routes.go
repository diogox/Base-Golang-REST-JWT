package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/database"
	"github.com/diogox/REST-JWT/server/pkg/email"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/refresh_whitelist"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

var db server.SqlDB
var refreshTokenWhitelist *redis.Client
var emailClient *email.EmailClient
var jwtSecret []byte
var tokenDurationInMinutes int

type RouteOptions struct {
	// Server Configs
	JWTSecret              []byte
	TokenDurationInMinutes int

	// Databases Configs
	PrismaHost string
	RedisHost  string

	// Email Service Configs
	Email         string
	EmailHost     string
	EmailPort     int
	EmailUsername string
	EmailPassword string
}

func SetupRoutes(e *echo.Echo, opts RouteOptions) {

	// Instantiate Prisma client
	db = server.SqlDB( database.NewPrismaDB(opts.PrismaHost) )

	// Instantiate redis client
	redisClient, err := refresh_whitelist.NewWhitelist(opts.RedisHost)
	if err != nil {
		panic("Failed to connect to redis database: " + err.Error())
	}
	refreshTokenWhitelist = redisClient

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
	requireAuth := func(next echo.HandlerFunc) echo.HandlerFunc {
		// Middleware to check the token is of type `AuthToken`
		f := func(c echo.Context) error {
			if err := next(c); err != nil {
				c.Error(err)
			}

			tokn := c.Get("user").(*jwt.Token)

			// Check if valid
			if !token.AssertAndValidate(tokn, token.AuthToken) {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Message: "Invalid Token!",
				})
			}

			return nil
		}

		// Return both middlewares
		jwtMiddleware := middleware.JWT(jwtSecret)
		return jwtMiddleware(f)
	}

	// Auth
	e.POST("/api/auth/register", register)
	e.POST("/api/auth/verify", sendVerificationEmail)
	e.GET("/api/auth/verify/:token", verifyEmail)
	e.POST("/api/auth/login", login)
	e.POST("/api/auth/reset-password", sendPasswordResetEmail)
	e.POST("/api/auth/reset-password/:token", resetPassword)
	e.POST("/api/auth/refresh", refreshToken) // Refreshes the JWT token

	// Endpoint that requires authentication
	apiEndpoint := e.Group("/api", requireAuth)

	// API endpoints (TODO: Add endpoints here!)
	apiEndpoint.GET("/users", handleGetUsers)
}

func handleGetUsers(c echo.Context) error {
	ctx := c.Request().Context()
	//logger := c.Logger()

	tokn := c.Get("user").(*jwt.Token)

	claims := tokn.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	users, err := db.GetUserByUsername(ctx, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, users)
}
