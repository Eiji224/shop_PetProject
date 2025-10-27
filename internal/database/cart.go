package database

import (
	"context"

	"gorm.io/gorm"
)

type CartModel struct {
	DB *gorm.DB
}

type Cart struct {
	ID     uint `gorm:"primaryKey;AUTO_INCREMENT	"`
	UserID uint `gorm:"not null"`

	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (cm *CartModel) Create(ctx context.Context, userID uint) error {
	cart := &Cart{UserID: userID}
	err := cm.DB.WithContext(ctx).Create(cart).Error

	return err
}

func (cm *CartModel) Get(ctx context.Context, userID uint) (*Cart, error) {
	var cart Cart
	err := cm.DB.WithContext(ctx).First(&cart, "user_id = ?", userID).Error

	return &cart, err
}

func (cm *CartModel) Preload(ctx context.Context, user *User) error {
	return cm.DB.WithContext(ctx).Preload("Cart").Where("id = ?", user.ID).First(user).Error
}
