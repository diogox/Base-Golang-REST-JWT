package routes

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/REST-JWT/server/pkg/email"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/prisma-client"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func sendVerificationEmail(c echo.Context) error {

	// Get context
	ctx := c.Request().Context()

	// Get logger
	logger := c.Logger()

	// Request body
	var req struct {
		Email string `json:"email" validate:"email,required"`
	}

	// Get POST body
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request body!",
		})
	}

	// Validate request
	err = c.Validate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Get user
	query := prisma.UserWhereUniqueInput{
		Email: &req.Email,
	}

	reqUser, err := client.User(query).Exec(ctx)
	if err != nil {
		// TODO: Maybe it's more helpful to specify that the username doesn't exist?
		// No user found
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "No user found with that email!",
		})
	}

	// Check if email is verified
	if reqUser.IsEmailVerified {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Email already verified!",
		})
	}

	// Create verification token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = reqUser.ID
	claims["type"] = "verification"
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenDurationInMinutes)).Unix()

	// Generate encoded token and send it as response.
	verificationToken, err := token.SignedString(jwtSecret)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	logger.Info("Verify Token: " + verificationToken)

	// Create `User` response
	user := models.User{
		Email:    reqUser.Email,
		Username: reqUser.Username,
	}

	// Send verification email
	err = emailClient.SendEmail(user, email.NewEmailOptions{
		Subject: "Registration",
		Message: fmt.Sprintf("Congrats %s you are now a user. Use this token to verify your account.", user.Username, verificationToken),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to send verification email!\n" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}
