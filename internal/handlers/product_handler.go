package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetAllProducts return all products
// @Summary Returns all products
// @Description Returns all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} []database.Product
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/products [get]
func (r *Router) GetAllProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	products, err := r.models.Products.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
	}

	c.JSON(http.StatusOK, products)
}
