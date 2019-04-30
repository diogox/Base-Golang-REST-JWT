package routes

import (
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/database"
	"github.com/diogox/REST-JWT/server/pkg/email"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/refresh_whitelist"
	"github.com/diogox/REST-JWT/server/pkg/routes/custom_middleware/authentication"
	"github.com/labstack/echo"
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

	// Permissions
	allowAllUsers := authentication.RequireAuth(jwtSecret)
	//allowOnlyPremiumUsers := authentication.RequireAuth(jwtSecret, authentication.PREMIUM_USER_ROLE)

	// Auth
	e.POST("/api/auth/login", login)
	e.POST("/api/auth/register", register)
	e.POST("/api/auth/logout", logout, allowAllUsers)
	e.POST("/api/auth/verify", sendVerificationEmail)
	e.GET("/api/auth/verify/:token", verifyEmail)
	e.POST("/api/auth/reset-password", sendPasswordResetEmail)
	e.POST("/api/auth/reset-password/:token", resetPassword)
	e.POST("/api/auth/refresh", refreshToken) // Refreshes the JWT token

	// Endpoint that requires authentication
	apiEndpoint := e.Group("/api", allowAllUsers)

	// API endpoints (TODO: Add endpoints here!)
	apiEndpoint.GET("/users", handleGetUsers)
}

func handleGetUsers(c echo.Context) error {
	ctx := c.Request().Context()
	//logger := c.Logger()

	userID, _ := c.Get(authentication.USER_ID_PARAM).(string)

	users, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, users)
}
