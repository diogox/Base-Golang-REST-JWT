package routes

import (
	"github.com/labstack/echo"
	"net/http"
)

func availableEmail(c echo.Context) error {
	ctx := c.Request().Context()

	email := c.Param("email")
	_, err := db.GetUserByEmail(ctx, email)
	if err != nil {
		return c.String(http.StatusOK, "Email is available.")
	}

	return c.String(http.StatusNotFound, "Email is already taken.")
}

