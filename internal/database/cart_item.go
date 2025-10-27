package database

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type CartItemModel struct {
	DB *gorm.DB
}

type CartItem struct {
	ID        uint      `gorm:"primaryKey;AUTO_INCREMENT"`
	CartID    uint      `gorm:"not null"`
	ProductID uint      `gorm:"not null"`
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`

	Cart    Cart    `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

func (cm *CartItemModel) Create(ctx context.Context, cartItem *CartItem) error {
	err := cm.DB.WithContext(ctx).Create(&cartItem).Error
	return err
}

func (cm *CartItemModel) GetCartItem(ctx context.Context, cartItemID uint) (*CartItem, error) {
	var cartItem CartItem
	err := cm.DB.WithContext(ctx).Find(&cartItem, "id = ?", cartItemID).Error
	return &cartItem, err
}

func (cm *CartItemModel) GetAllByCartID(ctx context.Context, cartID uint) ([]CartItem, error) {
	var cartItems []CartItem
	err := cm.DB.WithContext(ctx).Find(&cartItems, "cart_id = ?", cartID).Error

	return cartItems, err
}

func (cm *CartItemModel) UpdateQuantity(ctx context.Context, cartItemID uint, quantity int) error {
	cartItem, err := cm.GetCartItem(ctx, cartItemID)
	if err != nil {
		return err
	}

	cartItem.Quantity = quantity
	err = cm.DB.WithContext(ctx).Save(&cartItem).Error
	return err
}

func (cm *CartItemModel) DeleteCartItem(ctx context.Context, cartItemID uint) error {
	cartItem, err := cm.GetCartItem(ctx, cartItemID)
	if err != nil {
		return err
	}

	err = cm.DB.WithContext(ctx).Delete(&cartItem).Error
	return err
}

func (cm *CartItemModel) DeleteAllCartItems(ctx context.Context, cartID uint) error {
	var cartItems []CartItem
	err := cm.DB.WithContext(ctx).Find(&cartItems, "cart_id = ?", cartID).Error
	if err != nil {
		return err
	}

	for _, cartItem := range cartItems {
		err = cm.DeleteCartItem(ctx, cartItem.CartID)
		if err != nil {
			return err
		}
	}

	return nil
}
