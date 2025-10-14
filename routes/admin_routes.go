package routes

import (
	"1kosmetika-marketplace-backend/handlers"
	"1kosmetika-marketplace-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(r *gin.Engine, statsHandler *handlers.StatsHandler) {
	admin := r.Group("/api/admin")
	admin.Use(middlewares.JWTAuth(), middlewares.AdminOnly())

	// Статистика
	admin.GET("/stats", statsHandler.GetAdminStats)
	admin.GET("/stats/monthly", statsHandler.GetMonthlyStats)
	admin.GET("/stats/popular-products", statsHandler.GetPopularProducts)
	admin.GET("/stats/sales-by-category", statsHandler.GetSalesByCategory)
	admin.GET("/stats/charts", statsHandler.GetChartData)
	// new endpoints
	admin.GET("/stats/advanced", statsHandler.GetAdvancedStats)
    admin.GET("/stats/realtime", statsHandler.GetRealTimeStats)
    admin.GET("/stats/users", statsHandler.GetUserStats)
    admin.GET("/stats/dashboard", statsHandler.GetDashboardStats)
	admin.GET("/stats/user-purchases", statsHandler.GetUserPurchaseStats)
	// 2nd new endpoints
	admin.GET("/stats/traffic", statsHandler.GetTrafficStats)
	admin.GET("/stats/conversion", statsHandler.GetConversionStats)
	admin.GET("/stats/refunds", statsHandler.GetRefundStats)
	admin.GET("/stats/profit", statsHandler.GetProfitStats)

}