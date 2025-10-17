package database

import (
	"time"

	"gorm.io/gorm"
)

type OrderModel struct {
	DB *gorm.DB
}

type Order struct {
	ID         uint      `gorm:"primaryKey;AUTO_INCREMENT"`
	UserID     uint      `gorm:"not null"`
	TotalPrice float64   `gorm:"not null"`
	Status     string    `gorm:"size:100;not null"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
