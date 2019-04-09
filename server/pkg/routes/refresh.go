package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func refreshToken(c echo.Context) error {

	// Get token
	tokn := c.Get("user").(*jwt.Token)

	// Get logger
	logger := c.Logger()

	// Get expiration time
	claims := tokn.Claims.(jwt.MapClaims)
	expiration, _ := claims["exp"].(int64)

	// Make sure token can only be refreshed 30 seconds away from its expiration
	if time.Unix(expiration, 0).Sub(time.Now()) > 30*time.Second {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Token still valid for sufficient time. Try again later!",
		})
	}

	// Generate encoded token and send it as response.
	tokenStr, err := tokn.SignedString(jwtSecret)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Return new token
	return c.JSON(http.StatusOK, auth.Response{
		Token: tokenStr,
	})
}
