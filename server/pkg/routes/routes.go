package routes

import (
	"github.com/diogox/REST-JWT/server/pkg/blacklist"
	"net/http"

	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/database"
	"github.com/diogox/REST-JWT/server/pkg/email"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/routes/custom_middleware/authentication"
	"github.com/diogox/REST-JWT/server/pkg/whitelist"
	"github.com/labstack/echo"
)

var (
	// Server Configs
	AppUrl                        string
	jwtSecret                     []byte
	authTokenDurationInMinutes    int
	refreshTokenDurationInMinutes int

	// Databases
	db             server.DB
	tokenWhitelist server.Whitelist
	loginBlacklist server.Blacklist

	// Email Service
	emailService server.EmailService

	// Account Configs
	removeUnverifiedAccountAfterNDays    int
	accountAllowedNOfFailedLoginAttempts int
	accountLockDuration                  int
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

	// Account Configs
	RemoveUnverifiedAccountAfterNDays    int
	AccountAllowedNOfFailedLoginAttempts int
	AccountLockDuration                  int
}

func SetupRoutes(e *echo.Echo, opts RouteOptions) {
	AppUrl = opts.AppUrl

	/* Init Databases */
	db = server.DB(database.NewPrismaDB(opts.PrismaHost))

	whitelist, err := whitelist.NewWhitelist(opts.RedisHost)
	if err != nil {
		panic("Failed to connect to redis database: " + err.Error())
	}
	tokenWhitelist = whitelist

	blacklist, err := blacklist.NewBlacklist(opts.RedisHost, opts.AccountLockDuration)
	if err != nil {
		panic("Failed to connect to redis database: " + err.Error())
	}
	loginBlacklist = blacklist

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
	removeUnverifiedAccountAfterNDays = opts.RemoveUnverifiedAccountAfterNDays
	accountAllowedNOfFailedLoginAttempts = opts.AccountAllowedNOfFailedLoginAttempts
	accountLockDuration = opts.AccountLockDuration

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
	apiEndpoint.GET("/profile", profile)
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
