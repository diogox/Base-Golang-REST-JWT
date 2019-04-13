package database

import (
	"context"
	"github.com/diogox/REST-JWT/server/pkg/models"
	"github.com/diogox/REST-JWT/server/pkg/models/auth"
	"github.com/diogox/REST-JWT/server/prisma-client"
)

func NewPrismaDB(host string) *PrismaDB {
	// Create Config
	opts := prisma.Options{
		Endpoint: "http://" + host + ":4467",
	}

	// Make Client
	client := prisma.New(&opts)

	return &PrismaDB{
		client: client,
	}
}

type PrismaDB struct {
	client *prisma.Client
}

func (p *PrismaDB) GetUserByID(ctx context.Context, userId string) (*models.User, error) {

	// Get User
	user, err := p.client.User(prisma.UserWhereUniqueInput{
		ID: &userId,
	}).Exec(ctx)

	// Check If Not Found
	if err != nil {
		return nil, err
	}

	// Return User
	return &models.User{
		ID:              user.ID,
		Email:           user.Email,
		Username:        user.Username,
		Password:        user.Password,
		IsEmailVerified: user.IsEmailVerified,
	}, nil
}

func (p *PrismaDB) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {

	// Get User
	user, err := p.client.User(prisma.UserWhereUniqueInput{
		Username: &username,
	}).Exec(ctx)

	// Check If Not Found
	if err != nil {
		return nil, err
	}

	// Return User
	return &models.User{
		ID:              user.ID,
		Email:           user.Email,
		Username:        user.Username,
		Password:        user.Password,
		IsEmailVerified: user.IsEmailVerified,
	}, nil
}

func (p *PrismaDB) GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error) {

	// Get User
	user, err := p.client.User(prisma.UserWhereUniqueInput{
		Email: &userEmail,
	}).Exec(ctx)

	// Check If Not Found
	if err != nil {
		return nil, err
	}

	// Return User
	return &models.User{
		ID:              user.ID,
		Email:           user.Email,
		Username:        user.Username,
		Password:        user.Password,
		IsEmailVerified: user.IsEmailVerified,
	}, nil
}

func (p *PrismaDB) CreateUser(ctx context.Context, req *auth.NewRegistration) (*models.User, error) {

	// Request
	query := prisma.UserCreateInput{
		Email:           req.Email,
		Username:        req.Username,
		Password:        req.Password,
		IsEmailVerified: false,
	}

	// Create User
	res, err := p.client.CreateUser(query).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Return User Info
	return &models.User{
		ID:       res.ID,
		Email:    res.Email,
		Username: res.Username,
	}, nil
}

func (p *PrismaDB) UpdateUserByID(ctx context.Context, userID string, user *models.User) (*models.User, error) {

	// Update User
	updatedUser, err := p.client.UpdateUser(prisma.UserUpdateParams{
		Where: prisma.UserWhereUniqueInput{
			ID: &userID,
		},
		Data: prisma.UserUpdateInput{
			Email:           &user.Email,
			Username:        &user.Username,
			Password:        &user.Password,
			IsEmailVerified: &user.IsEmailVerified,
		},
	}).Exec(ctx)

	// Check If Update Failed
	if err != nil {
		return nil, err
	}

	// Return Updated User
	return &models.User{
		ID:              updatedUser.ID,
		Email:           updatedUser.Email,
		Username:        updatedUser.Username,
		Password:        updatedUser.Password,
		IsEmailVerified: updatedUser.IsEmailVerified,
	}, nil
}
