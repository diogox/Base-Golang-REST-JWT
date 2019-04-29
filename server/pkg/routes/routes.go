package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/database"
	"github.com/diogox/REST-JWT/server/pkg/email"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/refresh_whitelist"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

var (
	AppUrl                        string
	db                            server.SqlDB
	refreshTokenWhitelist         server.InMemoryDB
	emailService                  server.EmailService
	jwtSecret                     []byte
	authTokenDurationInMinutes    int
	refreshTokenDurationInMinutes int
)

type RouteOptions struct {
	// Server Configs
	AppUrl                        string
	JWTSecret                     []byte
	AuthTokenDurationInMinutes    int
	RefreshTokenDurationInMinutes int

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
	AppUrl = opts.AppUrl

	// Instantiate Prisma client
	db = server.SqlDB(database.NewPrismaDB(opts.PrismaHost))

	// Instantiate redis client
	whitelist, err := refresh_whitelist.NewWhitelist(opts.RedisHost)
	if err != nil {
		panic("Failed to connect to redis database: " + err.Error())
	}
	refreshTokenWhitelist = whitelist

	// Instantiate email client
	emailOpts := email.EmailClientOptions{
		Host:     opts.EmailHost,
		Port:     opts.EmailPort,
		Username: opts.EmailUsername,
		Password: opts.EmailPassword,
	}
	emailService = server.EmailService(email.NewEmailClient(opts.Email, emailOpts))

	// Set vars from options
	jwtSecret = opts.JWTSecret
	authTokenDurationInMinutes = opts.AuthTokenDurationInMinutes
	refreshTokenDurationInMinutes = opts.RefreshTokenDurationInMinutes

	e.Validator = newValidator()

	// Serve website
	e.File("/*", "./web/build/index.html")
	//e.File("/favicon.ico", "../web/images/favicon.ico")
	e.Static("/static", "./web/build/static")

	// By default, the key is extracted from the header "Authorization".
	// To get it from a field named `token` in the JSON we could add `TokenLookup: "query:token"` to the JWT Configs
	requireAuth := func(next echo.HandlerFunc) echo.HandlerFunc {
		// Middleware to check the token is of type `AuthToken`
		f := func(c echo.Context) error {
			t := c.Get("user").(*jwt.Token)

			// Check if valid (the jwt middleware already does this, but we might want to do additional checks...)
			if !token.AssertAndValidate(t, token.AuthToken) {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Message: "Invalid Token!",
				})
			}

			// Set user id to context
			claims := t.Claims.(jwt.MapClaims)
			userID := claims["user_id"].(string)
			c.Set("userID", userID)

			return next(c)
		}

		// Return both middleware
		jwtMiddleware := middleware.JWT(jwtSecret)
		return jwtMiddleware(f)
	}

	// Auth
	e.POST("/api/auth/login", login)
	e.POST("/api/auth/register", register)
	e.POST("/api/auth/logout", logout, requireAuth)
	e.POST("/api/auth/verify", sendVerificationEmail)
	e.GET("/api/auth/verify/:token", verifyEmail)
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

	userID, _ := c.Get("userID").(string)

	users, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, users)
}
