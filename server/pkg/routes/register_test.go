package routes

import (
	"encoding/json"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/routes/tests"
	"github.com/labstack/echo"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	newRegistration = auth.NewRegistration{
		Email:    "email@gmail.com",
		Username: "username",
		Password: "Password",
	}
)

func TestRegister(t *testing.T) {

	// Setup
	e := echo.New()

	reqJSON, _ := json.Marshal(newRegistration)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(string(reqJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	opts := RouteOptions{
		PrismaHost: "localhost",
		RedisHost:  "localhost",
		JWTSecret:  []byte("SuperSecretSecret"),
	}
	SetupRoutes(e, opts)

	// Assertions
	assert.Equal(t, registerHandler(c, tests.NewMockDb(), tests.NewMockEmailService()), nil)
	assert.Equal(t, rec.Code, http.StatusCreated)
}
