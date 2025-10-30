package repositories

import (
	"context"
	"shop/internal/models"

	"gorm.io/gorm"
)

type CartItemRepository interface {
	Create(cartItem *models.CartItem, ctx context.Context) error
	GetItem(cartItemID uint, ctx context.Context) (*models.CartItem, error)
	GetAllByCartID(cartID uint, ctx context.Context) ([]models.CartItem, error)
	UpdateQty(cartItemID uint, qty int, ctx context.Context) error
	DeleteItem(cartItemID uint, ctx context.Context) error
	DeleteAll(cartID uint, ctx context.Context) error
}

type cartItemRepository struct {
	db *gorm.DB
}

func (c *cartItemRepository) Create(cartItem *models.CartItem, ctx context.Context) error {
	return c.db.WithContext(ctx).Create(cartItem).Error
}

func (c *cartItemRepository) GetItem(cartItemID uint, ctx context.Context) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := c.db.WithContext(ctx).First(&cartItem, cartItemID).Error
	return &cartItem, err
}

func (c *cartItemRepository) GetAllByCartID(cartID uint, ctx context.Context) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	err := c.db.WithContext(ctx).Find(&cartItems, "cart_id = ?", cartID).Error
	return cartItems, err
}

func (c *cartItemRepository) UpdateQty(cartItemID uint, qty int, ctx context.Context) error {
	return c.db.WithContext(ctx).
		Model(&models.CartItem{}).
		Where("id = ?", cartItemID).
		Update("quantity", qty).
		Error
}

func (c *cartItemRepository) DeleteItem(cartItemID uint, ctx context.Context) error {
	return c.db.WithContext(ctx).Delete(&models.CartItem{}, cartItemID).Error
}

func (c *cartItemRepository) DeleteAll(cartID uint, ctx context.Context) error {
	return c.db.WithContext(ctx).
		Model(&models.CartItem{}).
		Where("cart_id = ?", cartID).
		Delete(&models.CartItem{}).Error
}

func NewCartItemRepository(db *gorm.DB) CartItemRepository {
	return &cartItemRepository{db: db}
}
