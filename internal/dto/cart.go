package dto

import "time"

type CartItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type CartItemUpdateRequest struct {
	Quantity int `json:"quantity"`
}

type CartItemResponse struct {
	ID        uint      `json:"id"`
	CartID    uint      `json:"cart_id"`
	ProductID uint      `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}
