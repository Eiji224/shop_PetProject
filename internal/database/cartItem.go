package database

import (
	"time"

	"gorm.io/gorm"
)

type CartItemModel struct {
	DB *gorm.DB
}

type CartItem struct {
	ID        uint      `gorm:"primaryKey;AUTO_INCREMENT"`
	CartID    uint      `gorm:"primaryKey;not null"`
	ProductID uint      `gorm:"primaryKey;not null"`
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`

	Cart    Cart    `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}
