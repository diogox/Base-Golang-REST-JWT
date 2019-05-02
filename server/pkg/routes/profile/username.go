package profile

import (
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/routes/custom_middleware/authentication"
	"github.com/labstack/echo"
	"net/http"
)

func SetUsername(db server.DB) func (c echo.Context) error {
	return func(c echo.Context) error {
		return handleSetUsername(c, db)
	}
}

func handleSetUsername(c echo.Context, db server.DB) error {
	ctx := c.Request().Context()
	//logger := c.Logger()

	// Request body
	var req struct {
		Username string `json:"username" validate:"isValidUsername,required"`
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

	// Check if Email is in usage
	_, err = db.GetUserByUsername(ctx, req.Username)
	if err == nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Username is already in use!",
		})
	}

	// Get user ID
	userID, _ := c.Get(authentication.USER_ID_PARAM).(string)

	// Get user
	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Message: "User not found!",
		})
	}

	// Change email
	user.Username = req.Username

	// Update user
	updatedUser, err := db.UpdateUserByID(ctx, userID, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to update user's email!",
		})
	}

	// Return user info
	return c.JSON(http.StatusOK, models.User{
		Email:      updatedUser.Email,
		Username:   updatedUser.Username,
		IsPaidUser: updatedUser.IsPaidUser,
	})
}
