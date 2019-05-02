package server

import (
	"context"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"time"
)

type DB interface {
	GetUserByID(ctx context.Context, userId string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error)
	CreateUser(ctx context.Context, newUser *auth.NewRegistration) (*models.User, error)
	UpdateUserByID(ctx context.Context, userID string, user *models.User) (*models.User, error)
	DeleteUserByID(ctx context.Context, userID string) (*models.User, error)
}

type Whitelist interface {
	// Refresh Token
	GetRefreshTokenByUserID(userID string) (string, error)
	SetRefreshTokenByUserID(userID string, value string, valueDurationInMinutes int) error
	DelRefreshTokenByUserID(userID string) (int64, error)

	// Reset Password Token
	GetResetPasswordTokenByUserID(userID string) (string, error)
	SetResetPasswordTokenByUserID(userID string, value string, valueDurationInMinutes int) error
	DelResetPasswordTokenByUserID(userID string) (int64, error)
}

type Blacklist interface {
	// Failed Login Limit
	GetFailedLoginCountByUserID(userID string) (string, *time.Duration, error)
	IncrementFailedLoginCountByUserID(userID string) error
	ResetFailedLoginCountByUserID(userID string) error
}