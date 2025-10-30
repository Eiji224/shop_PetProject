package handlers

import (
	"context"
	"net/http"
	"shop/internal/dto"
	"shop/internal/models"
	"shop/internal/repositories"
	"shop/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productRepository repositories.ProductRepository
	userService       *services.UserService
	productService    *services.ProductService
}

// GetProduct return product by id
// @Summary Returns product by id
// @Description Returns product by its id
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/products/{id} [get]
func (ph *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	product, err := ph.productRepository.GetProduct(uint(id), ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	response := dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ImageUrl:    product.ImageUrl,
		CategoryID:  product.CategoryID,
		UserID:      product.UserID,
		CreatedAt:   product.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetAllProducts return all products
// @Summary Returns all products
// @Description Returns all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} []dto.ProductResponse
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/products [get]
func (ph *ProductHandler) GetAllProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	products, err := ph.productRepository.GetAll(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}
	var response []dto.ProductResponse
	for _, product := range products {
		response = append(response, dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageUrl:    product.ImageUrl,
			CategoryID:  product.CategoryID,
			UserID:      product.UserID,
			CreatedAt:   product.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetProductsBySeller return seller products
// @Summary Returns seller products
// @Description Returns all products that were created by seller
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} []dto.ProductResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users/{id}/products [get]
func (ph *ProductHandler) GetProductsBySeller(c *gin.Context) {
	sellerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	products, err := ph.productRepository.GetAllBySellerID(uint(sellerID), ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}
	if len(products) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Products not found"})
		return
	}

	var response []dto.ProductResponse
	for _, product := range products {
		response = append(response, dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageUrl:    product.ImageUrl,
			CategoryID:  product.CategoryID,
			UserID:      product.UserID,
			CreatedAt:   product.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

// CreateProduct create product
// @Summary Creates product
// @Description Creates product and returns one
// @Tags Products
// @Accept json
// @Produce json
// @Param credentials body dto.CreateUpdateProductRequest true "Data for create product"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/products [post]
func (ph *ProductHandler) CreateProduct(c *gin.Context) {
	var createReq dto.CreateUpdateProductRequest

	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, status := ph.userService.GetUserFromContext(c)
	if status == http.StatusUnauthorized {
		c.AbortWithStatusJSON(status, gin.H{"error": "You are unauthorized"})
		return
	}
	if user.Type != string(dto.TypeSeller) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not allow to create product"})
		return
	}

	product := &models.Product{
		Name:        createReq.Name,
		Description: createReq.Description,
		Price:       createReq.Price,
		ImageUrl:    createReq.ImageUrl,
		CategoryID:  createReq.CategoryID,
		UserID:      user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := ph.productRepository.CreateProduct(product, ctx); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	if product == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	response := dto.ProductToResp(product)

	c.JSON(http.StatusCreated, response)
}

// UpdateProduct update existing product
// @Summary Updates existing product
// @Description Updates existing product and returns one
// @Tags Products
// @Accept json
// @Produce json
// @Param id path uint true "Product ID"
// @Param credentials body dto.CreateUpdateProductRequest true "Data for update existing product"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/products/{id} [put]
func (ph *ProductHandler) UpdateProduct(c *gin.Context) {
	id, user := ph.getUserAndID(c)
	if id == 0 && user == nil {
		return
	}

	productCtx, productCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer productCancel()
	existingProduct, status, err := ph.productService.GetProductIfAuthorized(id, user, productCtx)
	if err != nil {
		c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
		return
	}

	updateReq := dto.CreateUpdateProductRequest{}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedProduct := &models.Product{
		ID:          existingProduct.ID,
		Name:        updateReq.Name,
		Description: updateReq.Description,
		Price:       updateReq.Price,
		ImageUrl:    updateReq.ImageUrl,
		CategoryID:  updateReq.CategoryID,
		UserID:      existingProduct.UserID,
		CreatedAt:   existingProduct.CreatedAt,
	}

	ctx, updateCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer updateCancel()
	if err := ph.productRepository.UpdateProduct(updatedProduct, ctx); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	response := dto.ProductToResp(updatedProduct)

	c.JSON(http.StatusOK, response)
}

// DeleteProduct delete product
// @Summary Deletes existing product
// @Description Deletes existing product
// @Tags Products
// @Accept json
// @Produce json
// @Success 204 {object} nil
// @Param id path uint true "Product ID"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/products/{id} [delete]
func (ph *ProductHandler) DeleteProduct(c *gin.Context) {
	id, user := ph.getUserAndID(c)
	if id == 0 && user == nil {
		return
	}

	productCtx, productCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer productCancel()
	_, status, err := ph.productService.GetProductIfAuthorized(id, user, productCtx)
	if err != nil {
		c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := ph.productRepository.DeleteProduct(id, ctx); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (ph *ProductHandler) getUserAndID(c *gin.Context) (uint, *models.User) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 0, nil
	}
	user, status := ph.userService.GetUserFromContext(c)
	if status == http.StatusUnauthorized {
		c.AbortWithStatusJSON(status, gin.H{"error": "You are unauthorized"})
		return 0, nil
	}

	return uint(id), user
}

func NewProductHandler(
	productRepository repositories.ProductRepository,
	userService *services.UserService,
	productService *services.ProductService,
) *ProductHandler {
	return &ProductHandler{productRepository: productRepository, userService: userService, productService: productService}
}
