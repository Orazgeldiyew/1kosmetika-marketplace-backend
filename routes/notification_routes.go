package routes

import (
	"1kosmetika-marketplace-backend/handlers"
	"1kosmetika-marketplace-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupNotificationRoutes(r *gin.Engine, notificationHandler *handlers.NotificationHandler) {
	notif := r.Group("/api/notifications")
	notif.Use(middlewares.JWTAuth())
	{
		notif.GET("/", notificationHandler.GetUserNotifications)
		notif.GET("/unread-count", notificationHandler.GetUnreadCount)
		notif.PUT("/:id/read", notificationHandler.MarkAsRead)
		notif.PUT("/mark-all-read", notificationHandler.MarkAllAsRead)
		notif.POST("/", middlewares.AdminOnly(), notificationHandler.CreateNotification)
	}
}