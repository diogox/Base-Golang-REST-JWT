package routes

import (
	"encoding/json"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/routes/mocks"
	"github.com/labstack/echo"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const loginEndpoint = "/api/auth/login"
const loginMethod = http.MethodPost

var (
	// Valid Credential to use
	newLogin = auth.LoginCredentials{
		Username: "username",
		Password: "password",
	}

	loginInvalidPassword = "invalidPassword"

	// Testing table
	loginTT = []struct {
		name                 string
		isRegisteredUsername bool
		hasVerifiedEmail     bool
		isUsingValidPassword bool
		isUsingValidRequest  bool
		expectedStatusCode   int
	}{
		{
			name:                 "Sucessful Login",
			isRegisteredUsername: true,
			hasVerifiedEmail:     true,
			isUsingValidPassword: true,
			isUsingValidRequest:  true,
			expectedStatusCode:   http.StatusOK,
		},
		{
			name:                 "Invalid Password Fails",
			isRegisteredUsername: true,
			hasVerifiedEmail:     true,
			isUsingValidPassword: false,
			isUsingValidRequest:  true,
			expectedStatusCode:   http.StatusUnauthorized,
		},
		{
			name:                 "Unverified Email Fails",
			isRegisteredUsername: true,
			hasVerifiedEmail:     false,
			isUsingValidPassword: true,
			isUsingValidRequest:  true,
			expectedStatusCode:   http.StatusUnauthorized,
		},
		{
			name:                 "Unregistred User Fails",
			isRegisteredUsername: false,
			hasVerifiedEmail:     true,
			isUsingValidPassword: true,
			isUsingValidRequest:  true,
			expectedStatusCode:   http.StatusUnauthorized,
		},
		{
			name:                 "Invalid Request Fails",
			isRegisteredUsername: false,
			hasVerifiedEmail:     true,
			isUsingValidPassword: true,
			isUsingValidRequest:  false,
			expectedStatusCode:   http.StatusBadRequest,
		},
	}
)

func TestLogin(t *testing.T) {

	// For each test in the testing table
	for _, tc := range loginTT {
		t.Run(tc.name, func(t *testing.T) {

			// Setup
			e := SetupTestServer()

			// Create mock db
			db := mocks.NewMockDb()

			// Register user so that the login is valid
			if tc.isRegisteredUsername {
				RegisterTestUser(db, newLogin, tc.hasVerifiedEmail)
			}

			// Define credentials
			creds := newLogin
			if !tc.isUsingValidPassword {

				// Use invalid password instead
				creds.Password = loginInvalidPassword
			}

			if !tc.isUsingValidRequest {

				// Use invalid field value
				creds.Password = "" // This field is required, therefore, this should fail
			}

			// Marshal credentials
			reqJSON, _ := json.Marshal(creds)

			req := httptest.NewRequest(loginMethod, loginEndpoint, strings.NewReader(string(reqJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Assertions
			assert.Equal(t, loginHandler(c, db, mocks.NewMockInMemoryDB()), nil)
			assert.Equal(t, rec.Code, tc.expectedStatusCode)
		})
	}
}
