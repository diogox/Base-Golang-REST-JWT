package profile

import (
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/routes/custom_middleware/authentication"
	"github.com/labstack/echo"
	"net/http"
)

func GetProfile(db server.DB) func (c echo.Context) error {
	return func(c echo.Context) error {
		return handleGetProfile(c, db)
	}
}

func handleGetProfile(c echo.Context, db server.DB) error {
	ctx := c.Request().Context()
	//logger := c.Logger()

	// Get user ID
	userID, _ := c.Get(authentication.USER_ID_PARAM).(string)

	// Get user
	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Message: "User not found!",
		})
	}

	// Return user info
	return c.JSON(http.StatusOK, models.User{
		Email:      user.Email,
		Username:   user.Username,
		IsPaidUser: user.IsPaidUser,
	})
}
