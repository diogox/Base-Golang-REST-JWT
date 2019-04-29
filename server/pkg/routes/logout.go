package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/labstack/echo"
	"net/http"
)

func logout(c echo.Context) error {
	return logoutHandler(c, refreshTokenWhitelist)
}

// Invalidates the previous refresh_token
func logoutHandler(c echo.Context, whitelist server.InMemoryDB) error {

	// Get logger
	logger := c.Logger()

	// Get auth's user token
	tokn := c.Get("user").(*jwt.Token)

	claims := tokn.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(string)

	// Remove the refresh token, for the requesting user, from the whitelist, if it exists.
	_, err := whitelist.Del(userId)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.String(http.StatusOK, "Successfully Logged Out!")
}
