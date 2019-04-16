package mocks

import (
	"context"
	"errors"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"strconv"
)

func NewMockDb() *MockDB {
	return &MockDB{
		nextId: 0,
		users:  make([]models.User, 0),
	}
}

type MockDB struct {
	nextId int
	users  []models.User
}

func (m *MockDB) GetUserByID(ctx context.Context, userId string) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == userId {
			return &user, nil
		}
	}

	return nil, errors.New("User not found")
}

func (m *MockDB) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, errors.New("User not found")
}

func (m *MockDB) GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == userEmail {
			return &user, nil
		}
	}

	return nil, errors.New("User not found")
}

func (m *MockDB) CreateUser(ctx context.Context, user *auth.NewRegistration) (*models.User, error) {
	newUser := models.User{
		ID:              strconv.Itoa(m.nextId),
		Email:           user.Email,
		Username:        user.Username,
		Password:        user.Password,
		IsEmailVerified: false,
	}
	m.users = append(m.users, newUser)

	m.nextId = m.nextId + 1
	return &newUser, nil
}

func (m *MockDB) UpdateUserByID(ctx context.Context, userID string, user *models.User) (*models.User, error) {
	for i, u := range m.users {
		if u.ID == userID {
			updatedUser := models.User{
				ID: u.ID,
				Email: user.Email,
				Username: user.Username,
				Password: user.Password,
				IsEmailVerified: user.IsEmailVerified,
			}

			m.users[i] = updatedUser

			return &updatedUser, nil
		}
	}

	return nil, errors.New("User not found")
}