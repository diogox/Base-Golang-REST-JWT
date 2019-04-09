package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const RefreshToken = "REFRESH_TOKEN"

type RefreshTokenOptions struct {
	JWTSecret         []byte
	UserId          string
	DurationInMinutes int
}

func NewRefreshTokenToken(opts RefreshTokenOptions) (string, error) {

	// Create verification token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = opts.UserId
	claims["type"] = RefreshToken
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(opts.DurationInMinutes)).Unix()

	// Generate encoded token
	refreshToken, err := token.SignedString(opts.JWTSecret)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
