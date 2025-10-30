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

type CartHandler struct {
	cartRepository     repositories.CartRepository
	cartItemRepository repositories.CartItemRepository
	productRepository  repositories.ProductRepository
	userService        *services.UserService
	cartService        *services.CartService
}

// AddCartItem add product to cart
// @Summary Add product to cart
// @Description Add existing product to cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param credentials body dto.CartItemRequest true "Data for add item to cart"
// @Success 201 {object} dto.CartItemResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/cart/item [post]
func (ch *CartHandler) AddCartItem(c *gin.Context) {
	user := ch.getValidatedUser(c)
	if user == nil {
		return
	}

	var itemReq dto.CartItemRequest
	if err := c.ShouldBindJSON(&itemReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productCtx, productCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer productCancel()
	if _, err := ch.productRepository.GetProduct(itemReq.ProductID, productCtx); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartCtx, cartCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cartCancel()
	cart, err := ch.cartRepository.Get(user.ID, cartCtx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	cartItem := &models.CartItem{
		CartID:    cart.ID,
		ProductID: itemReq.ProductID,
		Quantity:  itemReq.Quantity,
		CreatedAt: time.Now(),
	}

	itemCtx, itemCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer itemCancel()
	err = ch.cartItemRepository.Create(cartItem, itemCtx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	itemResponse := dto.CartItemResponse{
		ID:        cartItem.ID,
		CartID:    cart.ID,
		ProductID: cartItem.ProductID,
		Quantity:  itemReq.Quantity,
		CreatedAt: cartItem.CreatedAt,
	}
	c.JSON(http.StatusCreated, itemResponse)
}

// GetCartItems get cart items
// @Summary Gets all items in cart
// @Description Gets all items in cart
// @Tags Cart
// @Accept json
// @Produce json
// @Success 200 {object} []dto.CartItemResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/cart/item [get]
func (ch *CartHandler) GetCartItems(c *gin.Context) {
	user := ch.getValidatedUser(c)
	if user == nil {
		return
	}

	itemCtx, itemCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer itemCancel()
	cartItems, err := ch.cartItemRepository.GetAllByCartID(user.Cart.ID, itemCtx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var itemsResponse []dto.CartItemResponse
	for _, cartItem := range cartItems {
		itemsResponse = append(itemsResponse, dto.CartItemResponse{
			ID:        cartItem.ID,
			CartID:    cartItem.CartID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			CreatedAt: cartItem.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, itemsResponse)
}

// UpdateCartItemQuantity update cart item quantity
// @Summary Updates cart item quantity
// @Description Updates cart item quantity
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path uint true "Cart item ID"
// @Param Quantity body dto.CartItemUpdateRequest true "Item quantity"
// @Success 200 {object} map[string]string "ok"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/cart/item/{id} [patch]
func (ch *CartHandler) UpdateCartItemQuantity(c *gin.Context) {
	user := ch.getValidatedUser(c)
	if user == nil {
		return
	}

	var itemReq dto.CartItemUpdateRequest
	if err := c.ShouldBindJSON(&itemReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if itemReq.Quantity <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Quantity must be greater than 0"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	checkCtx, checkCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer checkCancel()
	if status, err := ch.cartService.CheckItemBelongsToCart(id, user, checkCtx); status != http.StatusOK {
		c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
		return
	}

	itemCtx, itemCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer itemCancel()
	if err := ch.cartItemRepository.UpdateQty(uint(id), itemReq.Quantity, itemCtx); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DeleteCartItem delete cart item
// @Summary Deletes item from cart
// @Description Deletes item from cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path uint true "Cart item ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/cart/item/{id} [delete]
func (ch *CartHandler) DeleteCartItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := ch.getValidatedUser(c)
	if user == nil {
		return
	}

	checkCtx, checkCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer checkCancel()
	if status, err := ch.cartService.CheckItemBelongsToCart(id, user, checkCtx); status != http.StatusOK {
		c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
		return
	}

	deleteCtx, deleteCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer deleteCancel()
	if err := ch.cartItemRepository.DeleteItem(uint(id), deleteCtx); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteAllCartItems delete all cart items
// @Summary Deletes all items from cart
// @Description Deletes all items from cart
// @Tags Cart
// @Accept json
// @Produce json
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/cart/item [delete]
func (ch *CartHandler) DeleteAllCartItems(c *gin.Context) {
	user := ch.getValidatedUser(c)
	if user == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := ch.cartItemRepository.DeleteAll(user.Cart.ID, ctx); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (ch *CartHandler) getValidatedUser(c *gin.Context) *models.User {
	user, status := ch.userService.GetUserFromContext(c)
	if status == http.StatusUnauthorized {
		c.AbortWithStatusJSON(status, gin.H{"error": "unauthorized"})
		return nil
	}
	status, err := ch.cartService.ValidateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
		return nil
	}

	return user
}

func NewCartHandler(
	cartRepository repositories.CartRepository,
	cartItemRepository repositories.CartItemRepository,
	productRepository repositories.ProductRepository,
	userService *services.UserService,
	cartService *services.CartService,
) *CartHandler {

	return &CartHandler{
		cartRepository:     cartRepository,
		cartItemRepository: cartItemRepository,
		productRepository:  productRepository,
		userService:        userService,
		cartService:        cartService,
	}
}
