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
}

type InMemoryDB interface {
	Get(key string) (string, error)
	Set(key string, value string, valueDurationInMinutes int) error
	Del(key string) (int64, error)
}
