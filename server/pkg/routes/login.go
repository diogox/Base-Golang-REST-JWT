package routes

import (
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/routes/custom_middleware/authentication"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func login(c echo.Context) error {
	return loginHandler(c, db, tokenWhitelist, loginBlacklist)
}

// For testing purposes
func loginHandler(c echo.Context, db server.DB, whitelist server.Whitelist, blacklist server.Blacklist) error {
	// The previous auth token doesn't get invalidated, we just have to rely on the short duration of each token.
	// New logins invalidate the previous refresh_token. Logouts do the same.

	// Get context
	ctx := c.Request().Context()

	// Get logger
	logger := c.Logger()

	// Request body
	var loginCreds auth.LoginCredentials

	// Get POST body
	err := c.Bind(&loginCreds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request body!",
		})
	}

	// Validate request
	err = c.Validate(loginCreds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Get user
	user, err := db.GetUserByUsername(ctx, loginCreds.Username)
	if err != nil {
		// No user found
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Make sure the account hasn't been locked
	count, err := blacklist.GetFailedLoginCountByUserID(user.ID)
	if err != nil {
		// It probably returned an empty string, which means the value is 0.
		count = "0"
	}

	if cnt, err := strconv.Atoi(count); err != nil || cnt >= 5 {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Too many failed attempts, your account has been locked. Try again in 10 minutes!",
		})
	}

	// Check if the password matches
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCreds.Password))
	if err != nil {
		// Increment 'Failed Login Attempt' counter
		_ = blacklist.IncrementFailedLoginCountByUserID(user.ID)

		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Reset failed login counter
	err = blacklist.ResetFailedLoginCountByUserID(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Check if email is verified
	if !user.IsEmailVerified {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Email not verified!",
		})
	}

	// Generate encoded token and send it as response.
	opts := token.AuthTokenOptions{
		JWTSecret:         jwtSecret,
		UserID:            user.ID,
		UserRole:          authentication.ResolveUserRole(user.IsPaidUser),
		DurationInMinutes: authTokenDurationInMinutes,
	}

	tokenStr, err := token.NewAuthToken(opts)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Generate refresh token
	refreshOpts := token.RefreshTokenOptions{
		JWTSecret:         jwtSecret,
		UserId:            user.ID,
		DurationInMinutes: refreshTokenDurationInMinutes,
	}

	refreshTokenStr, err := token.NewRefreshToken(refreshOpts)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Add the refresh token to the whitelist (and make it expire after a determined amount of time)
	err = whitelist.SetRefreshTokenByUserID(user.ID, refreshTokenStr, refreshTokenDurationInMinutes)
	if err != nil {
		logger.Error(err.Error())
	}

	// Create response
	res := auth.LoginResponse{
		AuthToken:                   tokenStr,
		RefreshToken:                refreshTokenStr,
		ExpirationIntervalInMinutes: authTokenDurationInMinutes,
	}

	return c.JSON(http.StatusOK, res)
}
