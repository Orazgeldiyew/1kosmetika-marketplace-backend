package routes

import (
	"1kosmetika-marketplace-backend/handlers"
	"1kosmetika-marketplace-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupCartRoutes(r *gin.Engine, cartHandler *handlers.CartHandler) {
	cart := r.Group("/api/cart")
	cart.Use(middlewares.JWTAuth())
	{
		cart.GET("/", cartHandler.GetCart)
		cart.POST("/items", cartHandler.AddToCart)
		cart.DELETE("/items/:id", cartHandler.RemoveFromCart)
		cart.DELETE("/clear", cartHandler.ClearCart)
	}
}