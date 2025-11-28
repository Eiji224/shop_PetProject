package services

import (
	"context"
	"errors"
	"net/http"
	"shop/internal/dto"
	"shop/internal/models"
	"shop/internal/repositories"

	"gorm.io/gorm"
)

type ProductService struct {
	productRepository repositories.ProductRepository
}

func (ps *ProductService) GetProductIfAuthorized(id uint, user *models.User, ctx context.Context) (*models.Product, int, error) {
	product, err := ps.productRepository.GetProduct(id, ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New("product not found")
		}
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve product")
	}

	if user.ID != product.UserID || user.Type != string(dto.TypeSeller) {
		return nil, http.StatusForbidden, errors.New("you are not allowed to modify this product")
	}

	return product, http.StatusOK, nil
}

func NewProductService(productRepository repositories.ProductRepository) *ProductService {
	return &ProductService{productRepository: productRepository}
}
