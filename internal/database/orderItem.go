package database

import "gorm.io/gorm"

type OrderItemModel struct {
	DB *gorm.DB
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey;AUTO_INCREMENT"`
	OrderID   uint `gorm:"not null"`
	ProductID uint `gorm:"not null"`
	Quantity  int  `gorm:"not null"`
	Price     int  `gorm:"not null"`

	Order   Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID"`
}
