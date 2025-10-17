package database

import "gorm.io/gorm"

type CartModel struct {
	DB *gorm.DB
}

type Cart struct {
	ID     uint `gorm:"primaryKey;AUTO_INCREMENT	"`
	UserID uint `gorm:"not null"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
