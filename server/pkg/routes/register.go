package routes

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/REST-JWT/server"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func register(c echo.Context) error {
	return registerHandler(c, db, emailService)
}

// For testing purposes
func registerHandler(c echo.Context, db server.SqlDB, emailService server.EmailService) error {
	// Get context
	ctx := c.Request().Context()

	// Get logger
	logger := c.Logger()

	// Request body
	var req auth.NewRegistration

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

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8
	// (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 8)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Create user
	newUser, err := db.CreateUser(ctx, &auth.NewRegistration{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Email or Username already in use!",
		})
	}

	// Generate encoded token to verify email
	opts := token.EmailVerificationTokenOptions{
		JWTSecret:         jwtSecret,
		DurationInMinutes: tokenDurationInMinutes,
		UserId:            newUser.ID,
	}

	verificationToken, err := token.NewEmailVerificationToken(opts)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// TODO: Delete this in the future!
	logger.Info("Verify Token: " + verificationToken)

	// Send Verification Email
	err = emailService.SendEmail(*newUser, models.NewEmail{
		Subject: "Registration",
		Message: fmt.Sprintf("Congrats %s you are now a user. Use this token to verify your account: %s.", newUser.Username, verificationToken),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to send verification email!\n" + err.Error(),
		})
	}

	// Return Successful Response
	return c.JSON(http.StatusCreated, models.User{
		Email:    newUser.Email,
		Username: newUser.Username,
	})
}

func verifyEmail(c echo.Context) error {

	// Get context
	ctx := c.Request().Context()

	// Get token
	tokenString := c.Param("token")
	tokn, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret, nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid token!",
		})
	}

	// Check if the token is valid
	if !token.AssertAndValidate(tokn, token.EmailVerificationToken) {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid token!",
		})
	}

	// Get user id associated with token
	claims := tokn.Claims.(jwt.MapClaims)
	userId, _ := claims["user_id"].(string)

	// Update user to have email verified
	user, err := db.GetUserByID(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Update user to have email verified
	updatedUser := models.User{
		ID: user.ID,
		Email: user.Email,
		Username: user.Username,
		Password: user.Password,
		IsEmailVerified: true,
	}

	if _, err := db.UpdateUserByID(ctx, userId, &updatedUser); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.String(http.StatusOK, "Verification Successful")
}
