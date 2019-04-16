package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/routes/mocks"
	"github.com/labstack/echo"
	"github.com/magiconair/properties/assert"
)

const registerEndpoint = "/api/auth/register"
const registerMethod = http.MethodPost

var (
	// Valid Credential to use
	newRegistration = auth.NewRegistration{
		Email:    "email@gmail.com",
		Username: "username",
		Password: "Password",
	}

	registerInvalidPassword = "2Short" // Password is too short

	// Testing table
	registerTT = []struct {
		name                 string
		isUsingValidPassword bool
		isUsingValidRequest  bool
		expectedStatusCode   int
	}{
		{
			name:                 "Sucessful Registration",
			isUsingValidPassword: true,
			isUsingValidRequest:  true,
			expectedStatusCode:   http.StatusCreated,
		},
		{
			name:                 "Password Too Short To Register",
			isUsingValidPassword: false,
			isUsingValidRequest:  true,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name:                 "Password Too Short To Register",
			isUsingValidPassword: true,
			isUsingValidRequest:  false,
			expectedStatusCode:   http.StatusBadRequest,
		},
	}
)

func TestRegister(t *testing.T) {

	for _, tc := range registerTT {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			e := SetupTestServer()

			// Create mock db
			db := mocks.NewMockDb()

			// Define credentials
			reqBody := newRegistration
			if !tc.isUsingValidPassword {

				// Use invalid password instead
				reqBody.Password = registerInvalidPassword
			}

			if !tc.isUsingValidRequest {

				// Use invalid field value
				reqBody.Password = "" // This field is required, therefore, this should fail
			}

			reqJSON, _ := json.Marshal(reqBody)

			req := httptest.NewRequest(registerMethod, registerEndpoint, strings.NewReader(string(reqJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Assertions
			assert.Equal(t, registerHandler(c, db, mocks.NewMockEmailService()), nil)
			assert.Equal(t, rec.Code,tc.expectedStatusCode)
		})
	}
}
