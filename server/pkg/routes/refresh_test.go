package routes

import (
	"context"
	"encoding/json"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/routes/mocks"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/labstack/echo"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const refreshEndpoint = "/api/auth/refresh"
const refreshMethod = http.MethodPost

var (
	// Testing table
	refreshTT = []struct {
		name               string
		isValidToken       bool
		expectedStatusCode int
	}{
		{
			name:               "Successful Refresh",
			isValidToken:       true,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Failed Refresh",
			isValidToken:       false,
			expectedStatusCode: http.StatusBadRequest,
		},
	}
)

func TestRefresh(t *testing.T) {

	// For each test in the testing table
	for _, tc := range refreshTT {
		t.Run(tc.name, func(t *testing.T) {

			// Create mock db
			db := mocks.NewMockDb()
			memoryDB := mocks.NewWhitelist()

			newUser, _ := db.CreateUser(context.Background(), &auth.NewRegistration{
				Email: "",
				Username: "",
				Password: "",
			})

			// Setup
			e := SetupTestServer()

			tokenStr, _ := token.NewRefreshToken(token.RefreshTokenOptions{
				UserId: newUser.ID,
				DurationInMinutes: 5,
				JWTSecret: []byte(testJWTSecret),
			})

			_ = memoryDB.SetRefreshTokenByUserID(newUser.ID, tokenStr, 5)

			// Request struct
			type RefreshReq struct {
				RefreshToken string `json:"refresh_token"`
			}

			refreshReq := RefreshReq {
				RefreshToken: tokenStr,
			}

			// If it's supposed to have an invalid token
			if !tc.isValidToken {
				refreshReq.RefreshToken = ""
			}

			// Marshal credentials
			reqJSON, _ := json.Marshal(refreshReq)

			req := httptest.NewRequest(refreshMethod, refreshEndpoint, strings.NewReader(string(reqJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Assertions
			errorGot := refreshTokenHandler(c, db, memoryDB)

			if tc.isValidToken {
				assert.Equal(t, errorGot, nil)

				memoryToken, _ := memoryDB.GetRefreshTokenByUserID(newUser.ID)

				if memoryToken == tokenStr {
					panic("Token should be invalidated!")
				}
			}

			assert.Equal(t, rec.Code, tc.expectedStatusCode)
		})
	}
}
