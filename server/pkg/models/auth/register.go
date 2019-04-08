package auth

type NewRegistration struct {
	Email    string `json:"email" validate:"email"`
	Username string `json:"username"`
	Password string `json:"password" validate:"isValidPassword"`
}
