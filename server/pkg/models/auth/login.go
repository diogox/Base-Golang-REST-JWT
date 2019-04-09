package auth

type LoginCredentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	AuthToken string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
	ExpirationIntervalInMinutes int `json:"expiration_interval_in_minutes"`
}
