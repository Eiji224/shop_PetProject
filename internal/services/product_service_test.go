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

type MockProductRepository struct {
	mock.Mock
}

func (repo *MockProductRepository) CreateProduct(product *models.Product, ctx context.Context) error {
	return repo.Called(product, ctx).Error(0)
}

func (repo *MockProductRepository) GetProduct(id uint, ctx context.Context) (*models.Product, error) {
	args := repo.Called(id, ctx)

	product, _ := args.Get(0).(*models.Product)
	return product, args.Error(1)
}

func (repo *MockProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	args := repo.Called(ctx)

	products, _ := args.Get(0).([]models.Product)
	return products, args.Error(1)
}

func (repo *MockProductRepository) GetAllBySellerID(sellerID uint, ctx context.Context) ([]models.Product, error) {
	args := repo.Called(sellerID, ctx)

	products, _ := args.Get(0).([]models.Product)
	return products, args.Error(1)
}

func (repo *MockProductRepository) UpdateProduct(product *models.Product, ctx context.Context) error {
	return repo.Called(product, ctx).Error(0)
}

func (repo *MockProductRepository) DeleteProduct(id uint, ctx context.Context) error {
	return repo.Called(id, ctx).Error(0)
}

func newTestProductService() (*ProductService, *MockProductRepository) {
	productRepository := &MockProductRepository{}
	productService := NewProductService(productRepository)

	return productService, productRepository
}

type ProductTestCase struct {
	name           string
	user           *models.User
	repoItem       *models.Product
	repoErr        error
	inputID        uint
	expectedStatus int
	expectedErr    error
}

func TestGetProductIfAuthorized(t *testing.T) {
	productService, productRepo := newTestProductService()
	ctx := context.Background()

	userSeller := &models.User{
		ID:   0,
		Type: "seller",
	}
	userCustomer := &models.User{
		ID:   0,
		Type: "customer",
	}
	userWrongId := &models.User{
		ID:   1,
		Type: "seller",
	}

	testCases := []ProductTestCase{
		{
			name: "successful get product",
			user: userSeller,
			repoItem: &models.Product{
				ID:         0,
				Price:      0,
				CategoryID: 0,
				UserID:     0,
				CreatedAt:  time.Now(),
			},
			repoErr:        nil,
			inputID:        0,
			expectedStatus: http.StatusOK,
			expectedErr:    nil,
		},
		{
			name: "product not found",
			user: userSeller,
			repoItem: &models.Product{
				ID:         0,
				Price:      0,
				CategoryID: 0,
				UserID:     0,
				CreatedAt:  time.Now(),
			},
			repoErr:        gorm.ErrRecordNotFound,
			inputID:        1,
			expectedStatus: http.StatusNotFound,
			expectedErr:    errors.New("product not found"),
		},
		{
			name: "Not allowed (Customer)",
			user: userCustomer,
			repoItem: &models.Product{
				ID:         0,
				Price:      0,
				CategoryID: 0,
				UserID:     0,
				CreatedAt:  time.Now(),
			},
			repoErr:        nil,
			inputID:        2,
			expectedStatus: http.StatusForbidden,
			expectedErr:    errors.New("you are not allowed to modify this product"),
		},
		{
			name: "Not allowed (wrong ID)",
			user: userWrongId,
			repoItem: &models.Product{
				ID:         0,
				Price:      0,
				CategoryID: 0,
				UserID:     0,
				CreatedAt:  time.Now(),
			},
			repoErr:        nil,
			inputID:        3,
			expectedStatus: http.StatusForbidden,
			expectedErr:    errors.New("you are not allowed to modify this product"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			productRepo.On("GetProduct", tc.inputID, ctx).Return(tc.repoItem, tc.repoErr)

			_, status, err := productService.GetProductIfAuthorized(tc.inputID, tc.user, ctx)

			assert.Equal(t, tc.expectedStatus, status)
			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

			productRepo.AssertExpectations(t)
		})
	}
}
