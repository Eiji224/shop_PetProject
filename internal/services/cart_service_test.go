package services

import (
	"context"
	"errors"
	"net/http"
	"shop/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockCartRepository struct {
	mock.Mock
}

func (repo *MockCartRepository) Get(userID uint, ctx context.Context) (*models.Cart, error) {
	args := repo.Called(userID, ctx)

	cart, _ := args.Get(0).(*models.Cart)
	return cart, args.Error(1)
}

func (repo *MockCartRepository) Create(cart *models.Cart, ctx context.Context) error {
	return repo.Called(cart, ctx).Error(0)
}

type MockCartItemRepository struct {
	mock.Mock
}

func (repo *MockCartItemRepository) Create(cartItem *models.CartItem, ctx context.Context) error {
	return repo.Called(cartItem, ctx).Error(0)
}

func (repo *MockCartItemRepository) GetItem(cartItemID uint, ctx context.Context) (*models.CartItem, error) {
	args := repo.Called(cartItemID, ctx)

	cartItem, _ := args.Get(0).(*models.CartItem)
	return cartItem, args.Error(1)
}

func (repo *MockCartItemRepository) GetAllByCartID(cartID uint, ctx context.Context) ([]models.CartItem, error) {
	args := repo.Called(cartID, ctx)

	items, _ := args.Get(0).([]models.CartItem)
	return items, args.Error(1)
}

func (repo *MockCartItemRepository) UpdateQty(cartItemID uint, qty int, ctx context.Context) error {
	return repo.Called(cartItemID, qty, ctx).Error(0)
}

func (repo *MockCartItemRepository) DeleteItem(cartItemID uint, ctx context.Context) error {
	return repo.Called(cartItemID, ctx).Error(0)
}

func (repo *MockCartItemRepository) DeleteAll(cartID uint, ctx context.Context) error {
	return repo.Called(cartID, ctx).Error(0)
}

type checkItemTestCase struct {
	name           string
	repoItem       *models.CartItem // что вернёт репозиторий
	repoErr        error
	inputID        int // что передаём в сервис
	expectedStatus int
	expectedErr    error
}

func newCartService() (*MockCartRepository, *MockCartItemRepository, *CartService) {
	cartRepo := new(MockCartRepository)
	cartItemRepo := new(MockCartItemRepository)
	cartService := NewCartService(cartRepo, cartItemRepo)
	return cartRepo, cartItemRepo, cartService
}

func TestValidateUser(t *testing.T) {
	_, _, cartService := newCartService()
	userCustomer := &models.User{
		Type: "customer",
	}
	customerStatus, customerError := cartService.ValidateUser(userCustomer)

	userSeller := &models.User{
		Type: "seller",
	}
	sellerStatus, sellerError := cartService.ValidateUser(userSeller)

	assert.Equal(t, http.StatusOK, customerStatus)
	assert.Nil(t, customerError)

	assert.Equal(t, http.StatusForbidden, sellerStatus)
	assert.Equal(t, "user doesn't have permissions", sellerError.Error())
}

func TestCheckItemBelongsToCart(t *testing.T) {
	_, cartItemRepo, cartService := newCartService()
	cartID := uint(0)

	ctx := context.Background()

	user := &models.User{
		ID:   0,
		Type: "customer",
		Cart: &models.Cart{ID: cartID},
	}

	testCases := []checkItemTestCase{
		{
			name: "item belongs to user's cart",
			repoItem: &models.CartItem{
				ID:        0,
				CartID:    cartID,
				ProductID: 0,
				Quantity:  1,
				CreatedAt: time.Now(),
			},
			repoErr:        nil,
			inputID:        0,
			expectedStatus: http.StatusOK,
			expectedErr:    nil,
		},
		{
			name:           "item not found",
			repoItem:       nil,
			repoErr:        gorm.ErrRecordNotFound,
			inputID:        1,
			expectedStatus: http.StatusNotFound,
			expectedErr:    errors.New("not exists"),
		},
		{
			name: "item belongs to another cart",
			repoItem: &models.CartItem{
				ID:        2,
				CartID:    99999,
				ProductID: 0,
				Quantity:  1,
				CreatedAt: time.Now(),
			},
			repoErr:        nil,
			inputID:        2,
			expectedStatus: http.StatusForbidden,
			expectedErr:    errors.New("you are not allow to do this"),
		},
	}

	for _, tc := range testCases {
		cartItemRepo.On("GetItem", uint(tc.inputID), ctx).Return(tc.repoItem, tc.repoErr)

		status, err := cartService.CheckItemBelongsToCart(tc.inputID, user, ctx)

		assert.Equal(t, tc.expectedStatus, status, tc.name)

		if tc.expectedErr == nil {
			assert.NoError(t, err, tc.name)
		} else {
			assert.EqualError(t, err, tc.expectedErr.Error(), tc.name)
		}

		cartItemRepo.AssertExpectations(t)
	}
}
