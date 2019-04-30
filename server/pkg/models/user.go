package models

type User struct {
	ID              string `json:"id,omitempty"`
	Email           string `json:"email,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	IsEmailVerified bool   `json:"is_email_verified,omitempty"`
	IsPaidUser      bool   `json:"is_paid_user,omitempty"`
}
