package database

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type ProductModel struct {
	DB *gorm.DB
}

type Product struct {
	ID          uint      `gorm:"primaryKey;AUTO_INCREMENT"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"size:255"`
	Price       float64   `gorm:"not null"`
	ImageUrl    string    `gorm:"size:255"`
	CategoryID  uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`

	Category Category `gorm:"foreignKey:CategoryID"`
}

func (pm *ProductModel) GetAll(ctx context.Context) ([]Product, error) {
	var products []Product
	err := pm.DB.WithContext(ctx).Find(&products).Error

	return products, err
}
