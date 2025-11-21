package services

import (
	"context"
	"shop/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockCartRepository struct {
	mock.Mock
}

func (repo *MockCartRepository) Get(userID uint, ctx context.Context) (*models.Cart, error) {
	args := repo.Called(ctx, userID)
	if cart, ok := args.Get(0).(*models.Cart); ok {
		return cart, nil
	}
	return nil, args.Error(1)
}

func (repo *MockCartRepository) Create(cart *models.Cart, ctx context.Context) error {
	return repo.Called(ctx, cart).Error(0)
}
