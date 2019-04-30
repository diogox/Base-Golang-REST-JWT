package server

import (
	"context"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
)

type SqlDB interface {
	GetUserByID(ctx context.Context, userId string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error)
	CreateUser(ctx context.Context, newUser *auth.NewRegistration) (*models.User, error)
	UpdateUserByID(ctx context.Context, userID string, user *models.User) (*models.User, error)
	DeleteUserByID(ctx context.Context, userID string) (*models.User, error)
}

type InMemoryDB interface {
	// Refresh Token
	GetRefreshTokenByUserID(key string) (string, error)
	SetRefreshTokenByUserID(key string, value string, valueDurationInMinutes int) error
	DelRefreshTokenByUserID(key string) (int64, error)

	// Reset Password Token
	GetResetPasswordTokenByUserID(key string) (string, error)
	SetResetPasswordTokenByUserID(key string, value string, valueDurationInMinutes int) error
	DelResetPasswordTokenByUserID(key string) (int64, error)
}
