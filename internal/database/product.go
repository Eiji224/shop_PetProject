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
	UserID      uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`

	Category Category `gorm:"foreignKey:CategoryID"`
}

func (pm *ProductModel) GetProduct(ctx context.Context, id uint) (*Product, error) {
	var product Product
	err := pm.DB.WithContext(ctx).First(&product, id).Error

	return &product, err
}

func (pm *ProductModel) GetAll(ctx context.Context) ([]Product, error) {
	var products []Product
	err := pm.DB.WithContext(ctx).Find(&products).Error

	return products, err
}

func (pm *ProductModel) CreateProduct(ctx context.Context, product *Product) error {
	err := pm.DB.WithContext(ctx).Create(product).Error
	return err
}

func (pm *ProductModel) UpdateProduct(ctx context.Context, product *Product) error {
	err := pm.DB.WithContext(ctx).Save(product).Error
	return err
}

func (pm *ProductModel) DeleteProduct(ctx context.Context, id uint) error {
	err := pm.DB.WithContext(ctx).Delete(&Product{}, id).Error
	return err
}
