package auth

type NewRegistration struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,isValidUsername"`
	Password string `json:"password" validate:"required,isValidPassword"`
}
