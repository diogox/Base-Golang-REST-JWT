package authentication

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

// By default, the key is extracted from the header "Authorization".
// To get it from a field named `token` in the JSON we could add `TokenLookup: "query:token"` to the JWT Configs
func RequireAuth(jwtSecret []byte, roles ...string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// Middleware to check the token is of type `AuthToken`
		f := func(c echo.Context) error {
			t := c.Get("user").(*jwt.Token)

			// Check if valid (the jwt middleware already does this, but we might want to do additional checks...)
			if !token.AssertAndValidate(t, token.AuthToken) {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Message: "Invalid Token!",
				})
			}

			// Verify the role
			claims := t.Claims.(jwt.MapClaims)
			userRole := claims["user_role"].(string)

			isAllowed := checkRoleAllowed(userRole, roles...)
			if !isAllowed {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Message: "Insufficient Privileges!",
				})
			}

			// Set the user's ID to the context
			userID := claims["user_id"].(string)
			c.Set(USER_ID_PARAM, userID)

			return next(c)
		}

		// Return both middleware
		jwtMiddleware := middleware.JWT(jwtSecret)
		return jwtMiddleware(f)
	}
}

func checkRoleAllowed(role string, allowed ...string) bool {
	// If there are no 'allowed' than ALL users are allowed
	if len(allowed) == 0 {
		return true
	}

	for _, a := range allowed {
		if a == role {
			return true
		}
	}

	return false
}