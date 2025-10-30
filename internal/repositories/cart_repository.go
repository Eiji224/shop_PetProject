package repositories

import (
	"context"
	"shop/internal/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	Create(cart *models.Cart, ctx context.Context) error
	Get(userID uint, ctx context.Context) (*models.Cart, error)
}

type cartRepository struct {
	db *gorm.DB
}

func (c *cartRepository) Create(cart *models.Cart, ctx context.Context) error {
	return c.db.WithContext(ctx).Create(cart).Error
}

func (c *cartRepository) Get(userID uint, ctx context.Context) (*models.Cart, error) {
	var cart models.Cart
	err := c.db.WithContext(ctx).First(&cart, "user_id = ?", userID).Error
	return &cart, err
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}
