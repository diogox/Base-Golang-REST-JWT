package auth

type LoginCredentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Token string `json:"token"`
}
