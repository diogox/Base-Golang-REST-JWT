package token

import "github.com/dgrijalva/jwt-go"

func AssertAndValidate(token *jwt.Token, expectedType string) bool {

	// Check that it is an authentication type
	claims := token.Claims.(jwt.MapClaims)
	tokenType := claims["type"].(string)

	// Assert that it's the right type
	if tokenType != expectedType {
		return false
	}

	// Check if it has not expired
	if !token.Valid {
		return false
	}

	// Perform specific check for the type of token
	switch tokenType {
	case AuthToken:
		return validateAuthToken(token)
	case EmailVerificationToken:
		return validateEmailVerificationToken(token)
	case PasswordResetToken:
		return validatePasswordResetToken(token)
	}

	return false
}

func validateAuthToken(token *jwt.Token) bool {
	// Add specific checks here

	return true
}

func validateEmailVerificationToken(token *jwt.Token) bool {
	// Add specific checks here

	return true
}

func validatePasswordResetToken(token *jwt.Token) bool {
	// Add specific checks here

	return true
}
