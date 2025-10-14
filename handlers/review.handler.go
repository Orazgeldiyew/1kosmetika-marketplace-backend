package handlers

import (
	"1kosmetika-marketplace-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService services.ReviewService
}

func NewReviewHandler(reviewService services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

type CreateReviewRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"max=500"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"max=500"`
}

// @Summary Получить отзывы товара
// @Description Получить все отзывы для конкретного товара
// @Tags reviews
// @Produce json
// @Param productId path int true "ID товара"
// @Success 200 {object} gin.H
// @Router /api/reviews/product/{productId} [get]
func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	reviews, avgRating, totalCount, err := h.reviewService.GetProductReviews(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews":     reviews,
		"avg_rating":  avgRating,
		"total_count": totalCount,
	})
}

// @Summary Получить мои отзывы
// @Description Получить все отзывы текущего пользователя
// @Tags reviews
// @Produce json
// @Success 200 {array} models.Review
// @Router /api/reviews/user [get]
// @Security BearerAuth
func (h *ReviewHandler) GetUserReviews(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	reviews, err := h.reviewService.GetUserReviews(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// @Summary Создать отзыв
// @Description Создать новый отзыв для товара
// @Tags reviews
// @Accept json
// @Produce json
// @Param input body CreateReviewRequest true "Данные отзыва"
// @Success 201 {object} models.Review
// @Failure 400 {object} gin.H
// @Router /api/reviews [post]
// @Security BearerAuth
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req CreateReviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewService.CreateReview(userID, req.ProductID, req.Rating, req.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// @Summary Обновить отзыв
// @Description Обновить существующий отзыв
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "ID отзыва"
// @Param input body UpdateReviewRequest true "Обновленные данные отзыва"
// @Success 200 {object} models.Review
// @Failure 400 {object} gin.H
// @Router /api/reviews/{id} [put]
// @Security BearerAuth
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewService.UpdateReview(uint(reviewID), userID, req.Rating, req.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// @Summary Удалить отзыв
// @Description Удалить отзыв
// @Tags reviews
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /api/reviews/{id} [delete]
// @Security BearerAuth
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := h.reviewService.DeleteReview(uint(reviewID), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}