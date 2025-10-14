package handlers

import (
	"1kosmetika-marketplace-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// @Summary Получить корзину
// @Description Получить корзину текущего пользователя
// @Tags cart
// @Produce json
// @Success 200 {object} models.Cart
// @Router /api/cart [get]
// @Security BearerAuth
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}
	c.JSON(http.StatusOK, cart)
}

// @Summary Добавить товар в корзину
// @Description Добавить товар в корзину пользователя
// @Tags cart
// @Accept json
// @Produce json
// @Param input body AddToCartRequest true "Данные товара"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /api/cart/items [post]
// @Security BearerAuth
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req AddToCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.AddToCart(userID, req.ProductID, req.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}

// @Summary Удалить товар из корзины
// @Description Удалить товар из корзины пользователя
// @Tags cart
// @Produce json
// @Param id path int true "ID элемента корзины"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /api/cart/items/{id} [delete]
// @Security BearerAuth
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	if err := h.cartService.RemoveFromCart(userID, uint(itemID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}

// @Summary Очистить корзину
// @Description Очистить всю корзину пользователя
// @Tags cart
// @Produce json
// @Success 200 {object} gin.H
// @Router /api/cart/clear [delete]
// @Security BearerAuth
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.cartService.ClearCart(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}