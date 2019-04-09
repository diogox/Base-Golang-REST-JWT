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
	"time"
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

	// Create verification token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = newUser.ID
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

	// Get logger
	logger := c.Logger()

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

	// Check if it's a verification token
	claims := tokn.Claims.(jwt.MapClaims)
	if claims["type"] != "verification" {
		logger.Info("Verification Token")

		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid token!",
		})
	}

	if !token.AssertAndValidate(tokn, token.EmailVerificationToken) {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Expired token!",
		})
	}

	userId, _ := claims["user_id"].(string)
	isVerified := true

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

func login(c echo.Context) error {
	// The previous token doesn't get invalidated, we just have to rely on the short duration of each token.
	// To invalidate, we'd need to hold a token 'blacklist' in a database (probably Redis), but we're not doing that here.

	loginError := models.ErrorResponse{
		Message: "Username or password incorrect!",
	}

	// Get context
	ctx := c.Request().Context()

	// Get logger
	logger := c.Logger()

	// Request body
	var loginCreds auth.LoginCredentials

	// Get POST body
	err := c.Bind(&loginCreds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request body!",
		})
	}

	// Validate request
	err = c.Validate(loginCreds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Get user
	query := prisma.UserWhereUniqueInput{
		Username: &loginCreds.Username,
	}

	user, err := client.User(query).Exec(ctx)
	if err != nil {
		// TODO: Maybe it's more helpful to specify that the username doesn't exist?
		// No user found
		return c.JSON(http.StatusUnauthorized, loginError)
	}

	// Check if the password matches
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCreds.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, loginError)
	}

	// Check if email is verified
	if !user.IsEmailVerified {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Email not verified!",
		})
	}

	// Generate encoded token and send it as response.
	opts := token.AuthTokenOptions{
		JWTSecret: jwtSecret,
		Username: loginCreds.Username,
		DurationInMinutes: tokenDurationInMinutes,
	}

	tokenStr, err := token.NewAuthToken(opts)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Create response
	res := auth.Response{
		Token: tokenStr,
	}

	return c.JSON(http.StatusOK, res)
}

func refreshToken(c echo.Context) error {

	// Get token
	tokn := c.Get("user").(*jwt.Token)

	// Get logger
	logger := c.Logger()

	// Get expiration time
	claims := tokn.Claims.(jwt.MapClaims)
	expiration, _ := claims["exp"].(int64)

	// Make sure token can only be refreshed 30 seconds away from its expiration
	if time.Unix(expiration, 0).Sub(time.Now()) > 30*time.Second {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Token still valid for sufficient time. Try again later!",
		})
	}

	// Generate encoded token and send it as response.
	tokenStr, err := tokn.SignedString(jwtSecret)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Return new token
	return c.JSON(http.StatusOK, auth.Response{
		Token: tokenStr,
	})
}

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