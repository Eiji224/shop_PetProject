package repositories

import (
	"context"
	"shop/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *models.Product, ctx context.Context) error
	GetProduct(id uint, ctx context.Context) (*models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	GetAllBySellerID(sellerID uint, ctx context.Context) ([]models.Product, error)
	UpdateProduct(product *models.Product, ctx context.Context) error
	DeleteProduct(id uint, ctx context.Context) error
}

type productRepository struct {
	db *gorm.DB
}

func (p *productRepository) CreateProduct(product *models.Product, ctx context.Context) error {
	return p.db.WithContext(ctx).Create(product).Error
}

func (p *productRepository) GetProduct(id uint, ctx context.Context) (*models.Product, error) {
	var product models.Product
	err := p.db.WithContext(ctx).First(&product, id).Error
	return &product, err
}

func (p *productRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := p.db.WithContext(ctx).Find(&products).Error
	return products, err
}

func (p *productRepository) GetAllBySellerID(sellerID uint, ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := p.db.WithContext(ctx).Find(&products, "user_id = ?", sellerID).Error
	return products, err
}

func (p *productRepository) UpdateProduct(product *models.Product, ctx context.Context) error {
	return p.db.WithContext(ctx).Save(product).Error
}

func (p *productRepository) DeleteProduct(id uint, ctx context.Context) error {
	return p.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}
