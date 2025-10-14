package routes

import (
	"1kosmetika-marketplace-backend/handlers"
	"1kosmetika-marketplace-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupReviewRoutes(r *gin.Engine, reviewHandler *handlers.ReviewHandler) {
	reviews := r.Group("/api/reviews")
	{
		reviews.GET("/product/:productId", reviewHandler.GetProductReviews)
		reviews.GET("/user", middlewares.JWTAuth(), reviewHandler.GetUserReviews)
		reviews.POST("/", middlewares.JWTAuth(), reviewHandler.CreateReview)
		reviews.PUT("/:id", middlewares.JWTAuth(), reviewHandler.UpdateReview)
		reviews.DELETE("/:id", middlewares.JWTAuth(), reviewHandler.DeleteReview)
	}
}