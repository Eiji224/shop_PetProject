package models

import (
	"time"
)

type Product struct {
	ID          uint      `gorm:"primaryKey;AUTO_INCREMENT"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"size:255"`
	Price       float64   `gorm:"not null"`
	ImageUrl    string    `gorm:"size:255"`
	CategoryID  uint      `gorm:"not null"`
	UserID      uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`

	Category Category `gorm:"foreignKey:CategoryID"`
}
