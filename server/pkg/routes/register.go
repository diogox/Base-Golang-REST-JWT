package routes

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/REST-JWT/server/pkg/email"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/pkg/token"
	"github.com/diogox/REST-JWT/server/prisma-client"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func register(c echo.Context) error {

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
	query := prisma.UserCreateInput{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
		IsEmailVerified: false,
	}

	newUser, err := client.CreateUser(query).Exec(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Email or Username already in use!",
		})
	}

	// Generate encoded token and send it as response.
	opts := token.EmailVerificationTokenOptions{
		JWTSecret: jwtSecret,
		DurationInMinutes: tokenDurationInMinutes,
		UserId: newUser.ID,
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
		Email:    newUser.Email,
		Username: newUser.Username,
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
		return c.JSON(http.StatusBadRequest, models.ErrorResponse {
			Message: "Invalid token!",
		})
	}

	// Get user id associated with token
	claims := tokn.Claims.(jwt.MapClaims)
	userId, _ := claims["user_id"].(string)

	isVerified := true
	// Update user to have email verified
	_, err = client.UpdateUser(prisma.UserUpdateParams{
		Where: prisma.UserWhereUniqueInput{
			ID: &userId,
		},
		Data: prisma.UserUpdateInput{
			IsEmailVerified: &isVerified,
		},
	}).Exec(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.String(http.StatusOK, "Verification Successful")
}
