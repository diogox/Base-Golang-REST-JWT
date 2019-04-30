package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/routes/custom_middleware/authentication"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/labstack/echo"
	"net/http"
)

func refreshToken(c echo.Context) error {
	return refreshTokenHandler(c, db, tokenWhitelist)
}

func refreshTokenHandler(c echo.Context, db server.SqlDB, whitelist server.InMemoryDB) error {

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

	// Get token from string
	refreshToken, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret, nil
	})

	// Get expiration time
	claims := refreshToken.Claims.(jwt.MapClaims)
	userId, _ := claims["user_id"].(string)

	// Make sure the token hasn't expired
	if !token.AssertAndValidate(refreshToken, token.RefreshToken) {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid Token!",
		})
	}

	// Get from whitelist
	previousRefreshToken, err := whitelist.GetRefreshTokenByUserID(userId)
	if err != nil {

		// Not found (most likely)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid Token!",
		})
	}

	// The given refresh_token must match the previous token, otherwise it's invalid
	if previousRefreshToken != req.RefreshToken {

		// Not found (most likely)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid Token!",
		})
	}

	// Generate encoded token and send it as response.
	refreshTokenOpts := token.RefreshTokenOptions{
		JWTSecret:         jwtSecret,
		DurationInMinutes: refreshTokenDurationInMinutes,
		UserId:            userId,
	}
	newRefreshTokenStr, err := token.NewRefreshToken(refreshTokenOpts)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// TODO: Refactor repetitive code into smaller functions

	// Add new token to whitelist (Replaces previous, if it exists)
	err = whitelist.SetRefreshTokenByUserID(userId, newRefreshTokenStr, refreshTokenDurationInMinutes)
	if err != nil {

		// Already exists (most likely)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Token already exists!",
		})
	}

	// Get user
	ctx := c.Request().Context()
	user, err := db.GetUserByID(ctx, userId)
	if err != nil {

		// Already exists (most likely)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to get user!",
		})
	}

	// Generate encoded token and send it as response.
	authTokenOpts := token.AuthTokenOptions{
		JWTSecret:         jwtSecret,
		DurationInMinutes: authTokenDurationInMinutes,
		UserRole:          authentication.ResolveUserRole(user.IsPaidUser),
		UserID:            userId,
	}
	newAuthtokenStr, err := token.NewAuthToken(authTokenOpts)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Return new token
	return c.JSON(http.StatusOK, auth.LoginResponse{
		AuthToken:                   newAuthtokenStr, // TODO: Probably should return different struct (don't need this field)
		RefreshToken:                newRefreshTokenStr,
		ExpirationIntervalInMinutes: authTokenDurationInMinutes, // TODO: Is this necessary?
	})
}
