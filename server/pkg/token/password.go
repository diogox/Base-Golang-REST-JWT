package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const PasswordResetToken = "PASSWORD_RESET_TOKEN"

type PasswordResetTokenOptions struct {
	JWTSecret         []byte
	UserId          string
	DurationInMinutes int
}

func NewPasswordResetToken(opts PasswordResetTokenOptions) (string, error) {

	// Create verification token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = opts.UserId
	claims["type"] = PasswordResetToken
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(opts.DurationInMinutes)).Unix()

	// Generate encoded token and send it as response.
	signedToken, err := token.SignedString(opts.JWTSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
