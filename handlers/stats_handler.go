package handlers

import (
	"1kosmetika-marketplace-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *services.StatsService
}

func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) GetAdminStats(c *gin.Context) {
	stats, err := h.statsService.GetAdminStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetMonthlyStats(c *gin.Context) {
	stats, err := h.statsService.GetMonthlyStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch monthly stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetPopularProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	products, err := h.statsService.GetPopularProducts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch popular products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *StatsHandler) GetSalesByCategory(c *gin.Context) {
	stats, err := h.statsService.GetSalesByCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales by category"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetChartData(c *gin.Context) {
	chartData, err := h.statsService.GetChartData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chart data"})
		return
	}
	c.JSON(http.StatusOK, chartData)
}

func (h *StatsHandler) GetAdvancedStats(c *gin.Context) {
	stats, err := h.statsService.GetAdvancedStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch advanced stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetRealTimeStats(c *gin.Context) {
	stats, err := h.statsService.GetRealTimeStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch realtime stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetUserStats(c *gin.Context) {
	stats, err := h.statsService.GetUserStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetUserPurchaseStats(c *gin.Context) {
	stats, err := h.statsService.GetUserPurchaseStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user purchase stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *StatsHandler) GetTrafficStats(c *gin.Context) {
	stats, err := h.statsService.GetTrafficStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traffic stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}


func (h *StatsHandler) GetConversionStats(c *gin.Context) {
	stats, err := h.statsService.GetConversionStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversion stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}


func (h *StatsHandler) GetRefundStats(c *gin.Context) {
	stats, err := h.statsService.GetRefundStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch refund stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}


func (h *StatsHandler) GetProfitStats(c *gin.Context) {
	stats, err := h.statsService.GetProfitStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profit stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
