package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	AuthToken              = "AUTH_TOKEN"
	EmailVerificationToken = "EMAIL_VERIFICATION_TOKEN"
)

type AuthTokenOptions struct {
	JWTSecret         []byte
	Username          string
	DurationInMinutes int
}

func NewAuthToken(opts AuthTokenOptions) (string, error) {
	// Create auth token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = opts.Username
	claims["type"] = AuthToken
	//claims["admin"] = false // TODO: Not needed for now...
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(opts.DurationInMinutes)).Unix()

	// Generate encoded token and send it as response.
	signedToken, err := token.SignedString(opts.JWTSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}