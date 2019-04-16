package routes

import (
	"context"

	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/routes/mocks"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func SetupTestServer() *echo.Echo {
	// Create server
	e := echo.New()

	// Setup the routes
	opts := RouteOptions{
		PrismaHost: "localhost",
		RedisHost:  "localhost",
		JWTSecret:  []byte("SuperSecretSecret"),
	}
	SetupRoutes(e, opts)

	return e
}

func RegisterTestUser(db *mocks.MockDB, login auth.LoginCredentials, shouldVerifyEmail bool) {

	// Create User
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(login.Password), 8)
	user, _ := db.CreateUser(context.Background(), &auth.NewRegistration{
		Email:    "email@email.com",
		Username: login.Username,
		Password: string(hashedPassword),
	})

	if shouldVerifyEmail {
		// Validate email
		_, _ = db.UpdateUserByID(context.Background(), user.ID, &models.User{
			ID:              user.ID,
			Email:           user.Email,
			Username:        user.Username,
			Password:        user.Password,
			IsEmailVerified: true,
		})
	}
}
