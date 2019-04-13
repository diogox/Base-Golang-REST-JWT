package routes

import (
	"context"
	"encoding/json"
	"github.com/diogox/REST-JWT/server/pkg/models"
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
	newLogin = auth.LoginCredentials{
		Username: "username",
		Password: "password",
	}
)

func TestLogin(t *testing.T) {

	// Setup
	e := echo.New()

	reqJSON, _ := json.Marshal(newRegistration)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(string(reqJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	opts := RouteOptions{
		PrismaHost: "localhost",
		RedisHost:  "localhost",
		JWTSecret:  []byte("SuperSecretSecret"),
	}
	SetupRoutes(e, opts)

	db := tests.NewMockDb()

	// TODO: THe problem here is that the password gets compared against the hashed version,
	//  and in the mock db, we're saving it as plain text.
	user, _ := db.CreateUser(context.Background(), &auth.NewRegistration{
		Email: "email@email.com",
		Username: "username",
		Password: "password",
	})
	_, _ = db.UpdateUserByID(context.Background(), user.ID, &models.User{
		IsEmailVerified: true,
	})

	// Assertions
	assert.Equal(t, loginHandler(c, tests.NewMockDb(), tests.NewMockInMemoryDB()), nil)
	assert.Equal(t, rec.Code, http.StatusOK)
}
