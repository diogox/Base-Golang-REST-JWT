package routes

import (
	"github.com/labstack/echo"
	"net/http"
)

func availableUsername(c echo.Context) error {
	ctx := c.Request().Context()

	username := c.Param("username")
	_, err := db.GetUserByUsername(ctx, username)
	if err != nil {
		return c.String(http.StatusOK, "Username is available.")
	}

	return c.String(http.StatusNotFound, "Username is already taken.")
}
