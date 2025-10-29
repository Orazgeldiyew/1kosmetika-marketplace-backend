package main

import (
	"log"

	"1kosmetika-marketplace-backend/config"
	"1kosmetika-marketplace-backend/database"
	"1kosmetika-marketplace-backend/handlers"
	"1kosmetika-marketplace-backend/middlewares"
	"1kosmetika-marketplace-backend/repositories"
	"1kosmetika-marketplace-backend/routes"
	"1kosmetika-marketplace-backend/scheduler"
	"1kosmetika-marketplace-backend/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Kosmetika Marketplace API
// @version 1.0
// @description API сервер для Kosmetika Marketplace
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	if err := database.ConnectDB(cfg); err != nil {
		log.Fatal("Database connection failed:", err)
	}
	if err := database.Migrate(); err != nil {
		log.Fatal("Database migration failed:", err)
	}

	// repos
	userRepo := repositories.NewUserRepository(database.DB)
	productRepo := repositories.NewProductRepository(database.DB)
	orderRepo := repositories.NewOrderRepository(database.DB)
	cartRepo := repositories.NewCartRepository(database.DB)
	favoriteRepo := repositories.NewFavoriteRepository(database.DB)
	reviewRepo := repositories.NewReviewRepository(database.DB)
	notificationRepo := repositories.NewNotificationRepository(database.DB)
	statsRepo := repositories.NewStatsRepository()

	// services
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(orderRepo, productRepo, cartRepo, notificationRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	favoriteService := services.NewFavoriteService(favoriteRepo, productRepo)
	reviewService := services.NewReviewService(reviewRepo, productRepo)
	notificationService := services.NewNotificationService(notificationRepo, userRepo)
	statsService := services.NewStatsService(statsRepo)

	// handlers
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	cartHandler := handlers.NewCartHandler(cartService)
	favoriteHandler := handlers.NewFavoriteHandler(favoriteService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	statsHandler := handlers.NewStatsHandler(statsService)

	r := gin.Default()
	r.MaxMultipartMemory = 10 << 20 // 10MB

	r.Use(middlewares.CORS())
	r.Static("/static", "./uploads")

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Kosmetika Marketplace Backend is running 🚀",
			"version": "1.0.0",
			"port":    cfg.ServerPort,
		})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// routes
	routes.SetupUserRoutes(r, userHandler)
	routes.SetupProductRoutes(r, productHandler)
	routes.SetupOrderRoutes(r, orderHandler)
	routes.SetupCartRoutes(r, cartHandler)
	routes.SetupFavoriteRoutes(r, favoriteHandler)
	routes.SetupReviewRoutes(r, reviewHandler)
	routes.SetupNotificationRoutes(r, notificationHandler)
	routes.SetupAdminRoutes(r, statsHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	scheduler.StartCronJobs()

	log.Printf("🚀 Server running on http://localhost:%s", cfg.ServerPort)
	log.Printf("📚 Swagger docs available on http://localhost:%5s/swagger/index.html", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
