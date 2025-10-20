package handlers

import (
	"context"
	"net/http"
	"shop/internal/database"
	"shop/internal/dto"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// getProduct return product by id
// @Summary Returns product by id
// @Description Returns product by its id
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product Id"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/products/{id} [get]
func (r *Router) getProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	product, err := r.models.Products.GetProduct(ctx, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	response := dto.ProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ImageUrl:    product.ImageUrl,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// getAllProducts return all products
// @Summary Returns all products
// @Description Returns all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} []dto.ProductResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/products [get]
func (r *Router) getAllProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	products, err := r.models.Products.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
	}
	var response []dto.ProductResponse
	for _, product := range products {
		response = append(response, dto.ProductResponse{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageUrl:    product.ImageUrl,
			CategoryID:  product.CategoryID,
			CreatedAt:   product.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

// createProduct create product
// @Summary Creates product
// @Description Creates product and returns one
// @Tags Products
// @Accept json
// @Produce json
// @Param credentials body dto.CreateUpdateProductRequest true "Data for create product"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/products [post]
func (r *Router) createProduct(c *gin.Context) {
	var createReq dto.CreateUpdateProductRequest

	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//cUser, exists := c.Get("user")
	//user, ok := cUser.(database.User)
	//if !exists || !ok {
	//	c.AbortWithStatus(http.StatusUnauthorized)
	//	return
	//}
	//
	//product.sellerId = user.ID

	product := database.Product{
		Name:        createReq.Name,
		Description: createReq.Description,
		Price:       createReq.Price,
		ImageUrl:    createReq.ImageUrl,
		CategoryID:  createReq.CategoryID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := r.models.Products.CreateProduct(ctx, &product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	response := dto.ProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ImageUrl:    product.ImageUrl,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// updateProduct update existing product
// @Summary Updates existing product
// @Description Updates existing product and returns one
// @Tags Products
// @Accept json
// @Produce json
// @Param id path uint true "Product Id"
// @Param credentials body dto.CreateUpdateProductRequest true "Data for update existing product"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/products/{id} [put]
func (r *Router) updateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//cUser, exists := c.Get("user")
	//user, ok := cUser.(database.User)
	//if !exists || !ok {
	//	c.AbortWithStatus(http.StatusUnauthorized)
	//	return
	//}
	ctx, getCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer getCancel()
	existingProduct, err := r.models.Products.GetProduct(ctx, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}
	if existingProduct == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	//if user.ID != existingProduct.ID {
	//	c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this product"})
	//	return
	//}
	updateReq := dto.CreateUpdateProductRequest{}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedProduct := database.Product{
		Name:        updateReq.Name,
		Description: updateReq.Description,
		Price:       updateReq.Price,
		ImageUrl:    updateReq.ImageUrl,
		CategoryID:  updateReq.CategoryID,
	}

	updatedProduct.ID = uint(id)
	updatedProduct.CreatedAt = existingProduct.CreatedAt

	ctx, updateCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer updateCancel()
	if err := r.models.Products.UpdateProduct(ctx, &updatedProduct); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	response := dto.ProductResponse{
		Id:          updatedProduct.ID,
		Name:        updatedProduct.Name,
		Description: updatedProduct.Description,
		Price:       updatedProduct.Price,
		ImageUrl:    updatedProduct.ImageUrl,
		CategoryID:  updatedProduct.CategoryID,
		CreatedAt:   updatedProduct.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// deleteProduct delete product
// @Summary Deletes existing product
// @Description Deletes existing product
// @Tags Products
// @Accept json
// @Produce json
// @Success 204 {object} nil
// @Param id path uint true "Product Id"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/products/{id} [delete]
func (r *Router) deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//cUser, exists := c.Get("user")
	//user, ok := cUser.(database.User)
	//if !exists || !ok {
	//	c.AbortWithStatus(http.StatusUnauthorized)
	//	return
	//}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	existingProduct, err := r.models.Products.GetProduct(ctx, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}
	if existingProduct == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	//if user.ID != existingProduct.ID {
	//	c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this product"})
	//	return
	//}
	if err := r.models.Products.DeleteProduct(ctx, uint(id)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
