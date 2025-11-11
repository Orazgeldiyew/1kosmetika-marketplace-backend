package handlers

import (
	"net/http"
	"strconv"

	"1kosmetika-marketplace-backend/services"

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

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}
	c.JSON(http.StatusOK, cart)
}

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

	c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart"})
}

func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	itemID := uint(id64)

	if err := h.cartService.RemoveFromCart(userID, itemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.cartService.ClearCart(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}
	c.Status(http.StatusNoContent)
}
