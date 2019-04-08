package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diogox/Calendoer/generated/prisma-client"
	"github.com/diogox/Calendoer/server/pkg/models"
	"github.com/diogox/Calendoer/server/pkg/models/auth"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func register(c echo.Context) error {

	// Get context
	ctx := c.Request().Context()

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
		Email: req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
	}

	newUser, err := client.CreateUser(query).Exec(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Email or Username already in use!",
		})
	}

	// TODO: Should I return a token, initially?

	return c.JSON(http.StatusOK, newUser)
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
		// No user found
		return c.JSON(http.StatusUnauthorized, loginError)
	}

	// Check if the password matches
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCreds.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, loginError)
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = loginCreds.Username
	claims["admin"] = false
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenDurationInMinutes)).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Create response
	res := auth.Response{
		Token: t,
	}

	return c.JSON(http.StatusOK, res)
}

func refreshToken(c echo.Context) error {

	// Get token
	token := c.Get("user").(*jwt.Token)

	// Get logger
	logger := c.Logger()

	// Check if valid
	if !token.Valid {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Invalid Token!",
		})
	}

	// Get expiration time
	claims := token.Claims.(jwt.MapClaims)
	expiration, _ := claims["exp"].(int64)

	// Make sure token can only be refreshed 30 seconds away from its expiration
	if time.Unix(expiration, 0).Sub(time.Now()) > 30*time.Second {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Token still valid for sufficient time. Try again later!",
		})
	}

	// Generate encoded token and send it as response.
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Return new token
	return c.JSON(http.StatusOK, auth.Response{
		Token: t,
	})
}