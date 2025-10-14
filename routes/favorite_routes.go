package routes

import (
	"1kosmetika-marketplace-backend/handlers"
	"1kosmetika-marketplace-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupFavoriteRoutes(r *gin.Engine, favoriteHandler *handlers.FavoriteHandler) {
	favorites := r.Group("/api/favorites")
	favorites.Use(middlewares.JWTAuth())
	{
		favorites.GET("/", favoriteHandler.GetFavorites)
		favorites.POST("/:productId", favoriteHandler.AddFavorite)
		favorites.DELETE("/:productId", favoriteHandler.RemoveFavorite)
		favorites.GET("/check/:productId", favoriteHandler.CheckFavorite)
	}
}