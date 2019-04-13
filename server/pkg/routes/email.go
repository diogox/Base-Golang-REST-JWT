package routes

import (
	"fmt"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/labstack/echo"
	"net/http"
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
	reqUser, err := db.GetUserByEmail(ctx, req.Email)
	if err != nil {
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

	// Generate encoded token and send it as response.
	opts := token.EmailVerificationTokenOptions{
		JWTSecret: jwtSecret,
		UserId: reqUser.ID,
		DurationInMinutes: tokenDurationInMinutes,
	}

	verificationToken, err := token.NewEmailVerificationToken(opts)
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
	err = emailService.SendEmail(user, models.NewEmail{
		Subject: "Registration",
		Message: fmt.Sprintf("Congrats %s you are now a user. Use this token to verify your account: %s", user.Username, verificationToken),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to send verification email!\n" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}
