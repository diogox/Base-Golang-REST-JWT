package routes

import (
	"fmt"
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/routes/custom_middleware/authentication"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func login(c echo.Context) error {
	return loginHandler(c, db, tokenWhitelist, loginBlacklist)
}

const incorrectCredentialsError = "Username and password don't match."
const invalidCredentialsError = "Username and/or password not valid!"
const emailNotVerifiedError = "Email not verified!"
const accountLockedError = "Too many failed attempts, your account has been locked."

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
			Message: invalidCredentialsError,
		})
	}

	// Get user
	user, err := db.GetUserByUsername(ctx, loginCreds.Username)
	if err != nil {
		// No user found
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: incorrectCredentialsError,
		})
	}

	// Make sure the account hasn't been locked
	countStr, timeUntilExpire, err := blacklist.GetFailedLoginCountByUserID(user.ID)
	if err != nil {
		// It probably returned an empty string, which means the value is 0.
		countStr = "0"
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		logger.Error("Failed to assert the number of failed logins!")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Unspecified Internal Error!",
		})
	}

	// Check if the password matches
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCreds.Password))
	if err != nil {

		// Should be locked
		if count + 1  >= accountAllowedNOfFailedLoginAttempts {
			msg := fmt.Sprintf("%s Try again in %.0f minutes or reset your password!", accountLockedError, timeUntilExpire.Minutes())

			return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: msg,
			})
		}

		// Can try again...
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: handleFailedLogin(logger, blacklist, user.ID, count).Error(),
		})
	}

	// Reset failed login counter
	err = blacklist.ResetFailedLoginCountByUserID(user.ID)
	if err != nil {
		logger.Error("Failed to reset 'failed-login' attempt counter!")
	}

	// Check if email is verified
	if !user.IsEmailVerified {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: emailNotVerifiedError,
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


// Returns the error to be used
func handleFailedLogin(logger echo.Logger, blacklist server.Blacklist, userID string, failedLoginCount int) error {
	count := failedLoginCount + 1

	// Increment 'Failed Login Attempt' counter
	err := blacklist.IncrementFailedLoginCountByUserID(userID)
	if err != nil {
		logger.Error("Failed to increment 'failed-login' count.")
	}

	attemptsLeft := accountAllowedNOfFailedLoginAttempts - count
	msg := fmt.Sprintf("%s You have %d tries left.", incorrectCredentialsError, attemptsLeft)
	return errors.New(msg)
}