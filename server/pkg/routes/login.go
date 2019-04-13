package routes

import (
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/refresh_whitelist"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/diogox/REST-JWT/server/prisma-client"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// TODO: Need to check refresh token whitelist for previous entries for a given user, if it exists,
//  remove it and return new refresh token.

func login(c echo.Context) error {
	// The previous token doesn't get invalidated, we just have to rely on the short duration of each token.
	// To invalidate, we'd need to hold a token 'blacklist' in a database (probably Redis), but we're not doing that here.

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
	query := prisma.UserWhereUniqueInput{
		Username: &loginCreds.Username,
	}

	user, err := client.User(query).Exec(ctx)
	if err != nil {
		// TODO: Maybe it's more helpful to specify that the username doesn't exist?
		// No user found
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Check if the password matches
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCreds.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
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
		JWTSecret: jwtSecret,
		Username: user.Username,
		DurationInMinutes: tokenDurationInMinutes,
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
		JWTSecret: jwtSecret,
		UserId: user.ID,
		DurationInMinutes: tokenDurationInMinutes,
	}

	refreshtokenStr, err := token.NewRefreshTokenToken(refreshOpts)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Add the refresh token to the whitelist (and make it expire after a determined amount of time)
	err = refresh_whitelist.AddToWhitelist(refreshTokenWhitelist, refreshtokenStr, tokenDurationInMinutes)
	if err != nil {
		panic(err)
	}

	// Create response
	res := auth.Response{
		AuthToken: tokenStr,
		RefreshToken: refreshtokenStr,
		ExpirationIntervalInMinutes: tokenDurationInMinutes,
	}

	return c.JSON(http.StatusOK, res)
}