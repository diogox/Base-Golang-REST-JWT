package routes

import (
	"encoding/json"
	"github.com/diogox/REST-JWT/server/pkg/routes/custom_middleware/authentication"
	"github.com/diogox/REST-JWT/server/pkg/routes/mocks"
	"github.com/labstack/echo"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogout(t *testing.T) {

	const logoutEndpoint = "/api/auth/logout"
	const logoutMethod = http.MethodPost

	userID := "super_secret_refresh_token"
	refreshToken := "secret_refresh_token"

	// For each test in the testing table
	t.Run("Successful Logout", func(t *testing.T) {

		// Setup
		e := SetupTestServer()

		// Create mock db
		memoryDB := mocks.NewWhitelist()

		// Marshal credentials
		var r struct {
			RefreshToken string `json:"refresh_token"`
		}

		r.RefreshToken = refreshToken
		reqJSON, _ := json.Marshal(r)

		req := httptest.NewRequest(logoutMethod, logoutEndpoint, strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Prep
		c.Set(authentication.USER_ID_PARAM, userID)

		// Assertions
		assert.Equal(t, logoutHandler(c, memoryDB), nil)
		assert.Equal(t, rec.Code, http.StatusOK)

		if value, err := memoryDB.GetRefreshTokenByUserID(userID); err != nil {
			assert.Equal(t, value, "")
		}
	})
}
