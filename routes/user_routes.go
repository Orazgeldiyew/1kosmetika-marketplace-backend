package routes

import (
	"1kosmetika-marketplace-backend/handlers"
	"1kosmetika-marketplace-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, userHandler *handlers.UserHandler) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.GET("/profile", middlewares.JWTAuth(), userHandler.GetProfile)
	}

	admin := r.Group("/api/admin")
	admin.Use(middlewares.JWTAuth(), middlewares.AdminOnly())
	{
		admin.GET("/users", userHandler.GetAllUsers)
		admin.PUT("/users/:id/role", userHandler.UpdateUserRole)
		admin.DELETE("/users/:id", userHandler.DeleteUser)
	}
}