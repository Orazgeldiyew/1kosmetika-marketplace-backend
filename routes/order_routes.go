package routes

import (
	"1kosmetika-marketplace-backend/handlers"
	"1kosmetika-marketplace-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(r *gin.Engine, orderHandler *handlers.OrderHandler) {
	orders := r.Group("/api/orders")
	orders.Use(middlewares.JWTAuth())
	{
		orders.POST("/", orderHandler.CreateOrder)
		orders.GET("/", orderHandler.GetUserOrders)
		orders.GET("/:id", orderHandler.GetOrderByID)
		
	
		orders.GET("/admin/all", middlewares.AdminOnly(), orderHandler.GetAllOrders)
	}
}