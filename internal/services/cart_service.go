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

type CartService struct {
	cartRepository     repositories.CartRepository
	cartItemRepository repositories.CartItemRepository
}

func (cs *CartService) ValidateUser(user *models.User) (int, error) {
	if user.Type != string(dto.TypeCustomer) {
		return http.StatusForbidden, errors.New("user doesn't have permissions")
	}

	return http.StatusOK, nil
}

func (cs *CartService) CheckItemBelongsToCart(id int, user *models.User, ctx context.Context) (int, error) {
	existingItem, err := cs.cartItemRepository.GetItem(uint(id), ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New("not exists")
		}
		return http.StatusInternalServerError, err
	}
	if user.Cart.ID != existingItem.CartID {
		return http.StatusForbidden, errors.New("you are not allow to do this")
	}

	return http.StatusOK, nil
}

func NewCartService(
	cartRepository repositories.CartRepository,
	cartItemRepository repositories.CartItemRepository) *CartService {
	return &CartService{
		cartRepository:     cartRepository,
		cartItemRepository: cartItemRepository,
	}
}
