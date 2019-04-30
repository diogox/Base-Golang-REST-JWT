package routes

import (
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/routes/custom_middleware/authentication"
	"github.com/labstack/echo"
	"net/http"
)

func logout(c echo.Context) error {
	return logoutHandler(c, tokenWhitelist)
}

// Invalidates the previous refresh_token
func logoutHandler(c echo.Context, whitelist server.InMemoryDB) error {

	// Get logger
	logger := c.Logger()

	// Get auth's user token
	userId := c.Get(authentication.USER_ID_PARAM).(string)

	// Remove the refresh token, for the requesting user, from the whitelist, if it exists.
	_, err := whitelist.DelRefreshTokenByUserID(userId)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.String(http.StatusOK, "Successfully Logged Out!")
}
