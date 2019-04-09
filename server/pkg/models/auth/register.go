package auth

type NewRegistration struct {
	Email    string `json:"email" validate:"email,required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"isValidPassword,required"`
}
