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

// @Summary Получить избранное
// @Description Получить список избранных товаров пользователя
// @Tags favorites
// @Produce json
// @Success 200 {array} models.Favorite
// @Router /api/favorites [get]
// @Security BearerAuth
func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")
	favorites, err := h.favoriteService.GetUserFavorites(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorites"})
		return
	}
	c.JSON(http.StatusOK, favorites)
}

// @Summary Добавить в избранное
// @Description Добавить товар в избранное
// @Tags favorites
// @Produce json
// @Param productId path int true "ID товара"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /api/favorites/{productId} [post]
// @Security BearerAuth
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

// @Summary Удалить из избранного
// @Description Удалить товар из избранного
// @Tags favorites
// @Produce json
// @Param productId path int true "ID товара"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /api/favorites/{productId} [delete]
// @Security BearerAuth
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

// @Summary Проверить избранное
// @Description Проверить, есть ли товар в избранном
// @Tags favorites
// @Produce json
// @Param productId path int true "ID товара"
// @Success 200 {object} gin.H
// @Router /api/favorites/check/{productId} [get]
// @Security BearerAuth
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