package routes

import (
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/routes/tests"
	"github.com/labstack/echo"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	newUserJSON = `{"email": "email@gmail.com","username": "Username","password": "Password"}`
)

func TestRegister(t *testing.T) {

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(newUserJSON))
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
	assert.Equal(t, registerHandler(c, server.SqlDB(tests.NewMockDb())), nil)

	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), newUserJSON)
}