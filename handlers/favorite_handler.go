package handlers

import (
	"1kosmetika-marketplace-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favoriteService services.FavoriteService
}

func NewFavoriteHandler(favoriteService services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favoriteService: favoriteService}
}

func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")
	favorites, err := h.favoriteService.GetUserFavorites(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorites"})
		return
	}
	c.JSON(http.StatusOK, favorites)
}


func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	
	if err := h.favoriteService.AddFavorite(userID, uint(productID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Added to favorites"})
}

func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	
	if err := h.favoriteService.RemoveFavorite(userID, uint(productID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Removed from favorites"})
}

func (h *FavoriteHandler) CheckFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	
	isFavorite, err := h.favoriteService.IsFavorite(userID, uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check favorite"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_favorite": isFavorite})
}