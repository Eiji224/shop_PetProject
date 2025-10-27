package handlers

import (
	"context"
	"errors"
	"net/http"
	"shop/internal/database"
	"shop/internal/dto"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// addCartItem add product to cart
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
func (r *Router) addCartItem(c *gin.Context) {
	user, status := getUser(r, c)
	if user == nil {
		c.AbortWithStatusJSON(status, gin.H{"error": "You are not allow to do this"})
		return
	}

	var itemReq dto.CartItemRequest
	if err := c.ShouldBindJSON(&itemReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productCtx, productCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer productCancel()
	if _, err := r.models.Products.GetProduct(productCtx, itemReq.ProductID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartCtx, cartCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cartCancel()
	cart, err := r.models.Carts.Get(cartCtx, user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	cartItem := &database.CartItem{
		CartID:    cart.ID,
		ProductID: itemReq.ProductID,
		Quantity:  itemReq.Quantity,
		CreatedAt: time.Now(),
	}

	itemCtx, itemCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer itemCancel()
	err = r.models.CartItems.Create(itemCtx, cartItem)
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

// getCartItems get cart items
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
func (r *Router) getCartItems(c *gin.Context) {
	user, status := getUser(r, c)
	if user == nil {
		c.AbortWithStatusJSON(status, gin.H{"error": "You are not allow to do this"})
		return
	}

	itemCtx, itemCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer itemCancel()
	cartItems, err := r.models.CartItems.GetAllByCartID(itemCtx, user.Cart.ID)
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

// updateCartItemQuantity update cart item quantity
// @Summary Updates cart item quantity
// @Description Updates cart item quantity
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path uint true "Cart item Id"
// @Param Quantity body dto.CartItemUpdateRequest true "Item quantity"
// @Success 200 {object} map[string]string "ok"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/cart/item/{id} [patch]
func (r *Router) updateCartItemQuantity(c *gin.Context) {
	user, status := getUser(r, c)
	if user == nil {
		c.AbortWithStatusJSON(status, gin.H{"error": "You are not allow to do this"})
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
	if err, status := checkItemBelongsCart(r, id, user); status != http.StatusOK {
		c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
		return
	}

	itemCtx, itemCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer itemCancel()
	if err := r.models.CartItems.UpdateQuantity(itemCtx, uint(id), itemReq.Quantity); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// deleteCartItem delete cart item
// @Summary Deletes item from cart
// @Description Deletes item from cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path uint true "Cart item Id"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/cart/item/{id} [delete]
func (r *Router) deleteCartItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, status := getUser(r, c)
	if user == nil {
		c.AbortWithStatusJSON(status, gin.H{"error": "You are not allow to do this"})
		return
	}

	if err, status := checkItemBelongsCart(r, id, user); status != http.StatusOK {
		c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
		return
	}

	deleteCtx, deleteCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer deleteCancel()
	if err := r.models.CartItems.DeleteCartItem(deleteCtx, uint(id)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// deleteAllCartItems delete all cart items
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
func (r *Router) deleteAllCartItems(c *gin.Context) {
	user, status := getUser(r, c)
	if user == nil {
		c.AbortWithStatusJSON(status, gin.H{"error": "You are not allow to do this"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := r.models.CartItems.DeleteAllCartItems(ctx, user.Cart.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func getUser(r *Router, c *gin.Context) (*database.User, int) {
	user := r.services.AuthService.GetUserFromContext(c)
	if user == nil {
		return nil, http.StatusUnauthorized
	}
	if user.Type != string(dto.TypeCustomer) {
		return nil, http.StatusForbidden
	}

	preloadCtx, preloadCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer preloadCancel()
	if err := r.models.Carts.Preload(preloadCtx, user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, http.StatusInternalServerError
	}

	return user, http.StatusOK
}

func checkItemBelongsCart(r *Router, id int, user *database.User) (error, int) {
	getCtx, getCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer getCancel()
	existingItem, err := r.models.CartItems.GetCartItem(getCtx, uint(id))
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if existingItem == nil {
		return errors.New("Not exist"), http.StatusNotFound
	}
	if user.Cart.ID != existingItem.CartID {
		return errors.New("You are not allow to do this"), http.StatusForbidden
	}

	return nil, http.StatusOK
}
