package dto

import (
	"shop/internal/models"
	"time"
)

type CreateUpdateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	ImageUrl    string  `json:"image_url"`
	CategoryID  uint    `json:"category_id" binding:"required,gt=0"`
}

type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageUrl    string    `json:"image_url"`
	CategoryID  uint      `json:"category_id"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func ProductToResp(product *models.Product) *ProductResponse {
	return &ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ImageUrl:    product.ImageUrl,
		CategoryID:  product.CategoryID,
		UserID:      product.UserID,
		CreatedAt:   product.CreatedAt,
	}
}
