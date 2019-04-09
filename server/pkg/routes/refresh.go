package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/diogox/REST-JWT/server/prisma-client"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func refreshToken(c echo.Context) error {

	// Get logger
	logger := c.Logger()

	// Request body
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	// Get POST body
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request body!",
		})
	}

	// Validate request
	err = c.Validate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Get from whitelist
	_, err = refreshTokenWhitelist.Get(req.RefreshToken).Result()
	if err != nil {

		// Not found (most likely)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid Token!",
		})
	}

	// Get token from string
	refreshToken, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret, nil
	})

	// Get expiration time
	claims := refreshToken.Claims.(jwt.MapClaims)
	expiration, _ := claims["exp"].(int64)
	userId, _ := claims["user_id"].(string)

	// Make sure token can only be refreshed 30 seconds away from its expiration
	if time.Unix(expiration, 0).Sub(time.Now()) > 30*time.Second {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Token still valid for sufficient time. Try again later!",
		})
	}

	// Generate encoded token and send it as response.
	refreshTokenOpts := token.RefreshTokenOptions{
		JWTSecret: jwtSecret,
		DurationInMinutes: tokenDurationInMinutes,
		UserId: userId,
	}
	newRefreshtokenStr, err := token.NewRefreshTokenToken(refreshTokenOpts)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// TODO: Refactor repetitive code into smaller functions

	// Remove previous token from whitelist
	_, err = refreshTokenWhitelist.Del(req.RefreshToken).Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Add new token to whitelist
	spareTime := 1
	expiresIn := time.Minute * time.Duration(tokenDurationInMinutes + spareTime)
	_, err = refreshTokenWhitelist.Set(newRefreshtokenStr, "", expiresIn).Result()
	if err != nil {

		// Not found (most likely)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid Token!",
		})
	}

	// Get user // TODO: Remove this later? Unnecessary overhead
	query := prisma.UserWhereUniqueInput{
		ID: &userId,
	}

	// Get context
	ctx := c.Request().Context()

	user, err := client.User(query).Exec(ctx)
	if err != nil {
		// TODO: Maybe it's more helpful to specify that the username doesn't exist?
		// No user found
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Failed to get user",
		})
	}

	// Generate encoded token and send it as response.
	authTokenOpts := token.AuthTokenOptions{
		JWTSecret: jwtSecret,
		DurationInMinutes: tokenDurationInMinutes,
		Username: user.Username, // TODO: !!!!IMPORTANT!!!! Need to change auth token's options to take user id instead
	}
	newAuthtokenStr, err := token.NewAuthToken(authTokenOpts)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Return new token
	return c.JSON(http.StatusOK, auth.Response{
		AuthToken: newAuthtokenStr, // TODO: Probably should return different struct (don't need this field)
		RefreshToken: newRefreshtokenStr,
		ExpirationIntervalInMinutes: tokenDurationInMinutes,
	})
}
