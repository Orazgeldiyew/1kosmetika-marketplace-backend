package routes

import (
	"1kosmetika-marketplace-backend/handlers"
	"1kosmetika-marketplace-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(r *gin.Engine, productHandler *handlers.ProductHandler) {
	products := r.Group("/api/products")
	{
		products.GET("/", productHandler.GetProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.GET("/paginated", productHandler.GetProductsPaginated)
		products.GET("/search", productHandler.SearchProducts)
		products.GET("/categories", productHandler.GetCategories)
		products.GET("/brands", productHandler.GetBrands)
		
		// Admin only routes
		adminRoutes := products.Group("")
		adminRoutes.Use(middlewares.JWTAuth(), middlewares.AdminOnly())
		{
			adminRoutes.POST("/", productHandler.CreateProduct)
			adminRoutes.PUT("/:id", productHandler.UpdateProduct)
			adminRoutes.DELETE("/:id", productHandler.DeleteProduct)
		}
	}
}