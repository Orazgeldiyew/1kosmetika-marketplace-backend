package handlers

import (
	"1kosmetika-marketplace-backend/models"
	"1kosmetika-marketplace-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// @Summary Получить уведомления
// @Description Получить уведомления пользователя с пагинацией
// @Tags notifications
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит на странице" default(20)
// @Success 200 {object} gin.H
// @Router /api/notifications [get]
// @Security BearerAuth
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	notifications, err := h.notificationService.GetUserNotifications(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"page":          page,
		"limit":         limit,
	})
}

// @Summary Получить количество непрочитанных
// @Description Получить количество непрочитанных уведомлений
// @Tags notifications
// @Produce json
// @Success 200 {object} gin.H
// @Router /api/notifications/unread-count [get]
// @Security BearerAuth
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unread count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// @Summary Отметить как прочитанное
// @Description Отметить уведомление как прочитанное
// @Tags notifications
// @Produce json
// @Param id path int true "ID уведомления"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /api/notifications/{id}/read [put]
// @Security BearerAuth
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	if err := h.notificationService.MarkAsRead(uint(notificationID), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// @Summary Отметить все как прочитанные
// @Description Отметить все уведомления как прочитанные
// @Tags notifications
// @Produce json
// @Success 200 {object} gin.H
// @Router /api/notifications/mark-all-read [put]
// @Security BearerAuth
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := h.notificationService.MarkAllAsRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark all as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

// @Summary Создать уведомление
// @Description Создать новое уведомление (только для админов)
// @Tags notifications
// @Accept json
// @Produce json
// @Param input body models.Notification true "Данные уведомления"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /api/notifications [post]
// @Security BearerAuth
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var input models.Notification
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data: " + err.Error()})
		return
	}

	if input.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is required"})
		return
	}
	if input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	input.IsRead = false

	if err := h.notificationService.CreateNotification(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Notification created",
		"notification":  input,
	})
}