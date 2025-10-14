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

// @Summary Получить статистику для администратора
// @Description Возвращает общие данные по пользователям, заказам, продажам и т.д.
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/stats [get]
// @Security BearerAuth
func (h *StatsHandler) GetAdminStats(c *gin.Context) {
	stats, err := h.statsService.GetAdminStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// @Summary Получить месячную статистику
// @Description Возвращает статистику по месяцам для графиков
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/stats/monthly [get]
// @Security BearerAuth
func (h *StatsHandler) GetMonthlyStats(c *gin.Context) {
	stats, err := h.statsService.GetMonthlyStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch monthly stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// @Summary Получить популярные товары
// @Description Возвращает список самых продаваемых товаров
// @Tags admin
// @Produce json
// @Param limit query int false "Количество товаров (по умолчанию 10)"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/stats/popular-products [get]
// @Security BearerAuth
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

// @Summary Получить продажи по категориям
// @Description Возвращает статистику продаж по категориям товаров
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/stats/sales-by-category [get]
// @Security BearerAuth
func (h *StatsHandler) GetSalesByCategory(c *gin.Context) {
	stats, err := h.statsService.GetSalesByCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales by category"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// @Summary Получить все данные для графиков
// @Description Возвращает полные данные для построения графиков
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/stats/charts [get]
// @Security BearerAuth
func (h *StatsHandler) GetChartData(c *gin.Context) {
	chartData, err := h.statsService.GetChartData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chart data"})
		return
	}
	c.JSON(http.StatusOK, chartData)
}

// @Summary Получить расширенную статистику
// @Description Возвращает расширенные метрики: средний чек, конверсия, LTV и т.д.
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/stats/advanced [get]
// @Security BearerAuth
func (h *StatsHandler) GetAdvancedStats(c *gin.Context) {
	stats, err := h.statsService.GetAdvancedStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch advanced stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// @Summary Получить статистику в реальном времени
// @Description Возвращает данные за сегодня: заказы, доходы, пользователи
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/stats/realtime [get]
// @Security BearerAuth
func (h *StatsHandler) GetRealTimeStats(c *gin.Context) {
	stats, err := h.statsService.GetRealTimeStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch realtime stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// @Summary Получить статистику пользователей
// @Description Возвращает аналитику по пользователям: активные, новые, администраторы
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/stats/users [get]
// @Security BearerAuth
func (h *StatsHandler) GetUserStats(c *gin.Context) {
	stats, err := h.statsService.GetUserStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// @Summary Получить все данные для дашборда
// @Description Возвращает полные данные для админ-панели
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/stats/dashboard [get]
// @Security BearerAuth
func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// ✅ Покупательская активность пользователей
// @Summary Получить аналитику покупок пользователей
// @Description Возвращает информацию: кто, когда, сколько заказов сделал и чем оплатил
// @Tags admin
// @Produce json
// @Success 200 {object} []repositories.OrderStats
// @Router /api/admin/stats/user-purchases [get]
// @Security BearerAuth
func (h *StatsHandler) GetUserPurchaseStats(c *gin.Context) {
	stats, err := h.statsService.GetUserPurchaseStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user purchase stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// ✅ Новый эндпоинт: Трафик
// @Summary Получить статистику трафика
// @Description Возвращает количество визитов, уникальных пользователей и среднюю продолжительность сессии
// @Tags admin
// @Produce json
// @Router /api/admin/stats/traffic [get]
// @Security BearerAuth
func (h *StatsHandler) GetTrafficStats(c *gin.Context) {
	stats, err := h.statsService.GetTrafficStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traffic stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// ✅ Новый эндпоинт: Конверсия
// @Summary Получить статистику конверсии
// @Description Возвращает данные о конверсии посетителей в покупателей
// @Tags admin
// @Produce json
// @Router /api/admin/stats/conversion [get]
// @Security BearerAuth
func (h *StatsHandler) GetConversionStats(c *gin.Context) {
	stats, err := h.statsService.GetConversionStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversion stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// ✅ Новый эндпоинт: Возвраты
// @Summary Получить статистику возвратов
// @Description Возвращает количество возвратов, их причины и процент от общего числа заказов
// @Tags admin
// @Produce json
// @Router /api/admin/stats/refunds [get]
// @Security BearerAuth
func (h *StatsHandler) GetRefundStats(c *gin.Context) {
	stats, err := h.statsService.GetRefundStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch refund stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// ✅ Новый эндпоинт: Прибыль
// @Summary Получить статистику прибыли
// @Description Возвращает доход, расходы и чистую прибыль
// @Tags admin
// @Produce json
// @Router /api/admin/stats/profit [get]
// @Security BearerAuth
func (h *StatsHandler) GetProfitStats(c *gin.Context) {
	stats, err := h.statsService.GetProfitStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profit stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
