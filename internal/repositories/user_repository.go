package repositories

import (
	"context"
	"shop/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Insert(user *models.User, ctx context.Context) error
	FindByID(id uint, ctx context.Context) (*models.User, error)
	FindByEmail(email string, ctx context.Context) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) Insert(user *models.User, ctx context.Context) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u *userRepository) FindByID(id uint, ctx context.Context) (*models.User, error) {
	var user models.User
	err := u.db.WithContext(ctx).Preload("Cart").First(&user, id).Error
	return &user, err
}

func (u *userRepository) FindByEmail(email string, ctx context.Context) (*models.User, error) {
	var user models.User
	err := u.db.WithContext(ctx).Preload("Cart").First(&user, "email = ?", email).Error
	return &user, err
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
