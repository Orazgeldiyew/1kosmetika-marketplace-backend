package repositories

import (
	"1kosmetika-marketplace-backend/database"
	"1kosmetika-marketplace-backend/models"
	"time"
)

type OrderStats struct {
	UserName       string
	TotalOrders    int64
	TotalSpent     float64
	MostUsedMethod string
	LastOrderDate  time.Time
}

type StatsRepository struct{}

func NewStatsRepository() *StatsRepository {
	return &StatsRepository{}
}

type AdminStats struct {
	TotalUsers    int64   `json:"total_users"`
	TotalOrders   int64   `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
	TotalProducts int64   `json:"total_products"`
	TotalReviews  int64   `json:"total_reviews"`
}

type MonthlyStats struct {
	Month        string  `json:"month"`
	Orders       int64   `json:"orders"`
	Revenue      float64 `json:"revenue"`
	NewUsers     int64   `json:"new_users"`
}

type ChartData struct {
	MonthlyStats     []MonthlyStats   `json:"monthly_stats"`
	PopularProducts  []ProductStats   `json:"popular_products"`
	SalesByCategory  []CategoryStats  `json:"sales_by_category"`
}

type ProductStats struct {
	ProductName string `json:"product_name"`
	SalesCount  int64  `json:"sales_count"`
}

type CategoryStats struct {
	CategoryName string  `json:"category_name"`
	TotalSales   float64 `json:"total_sales"`
	OrderCount   int64   `json:"order_count"`
}

type AdvancedStats struct {
	AverageOrderValue float64 `json:"average_order_value"`
	ConversionRate    float64 `json:"conversion_rate"`
	CustomerLifetime  float64 `json:"customer_lifetime_value"`
	RepeatCustomers   int64   `json:"repeat_customers"`
}

type RealTimeStats struct {
	TodayOrders    int64   `json:"today_orders"`
	TodayRevenue   float64 `json:"today_revenue"`
	TodayUsers     int64   `json:"today_users"`
	PendingOrders  int64   `json:"pending_orders"`
}

type UserStats struct {
	ActiveUsers     int64 `json:"active_users"`
	NewUsersToday   int64 `json:"new_users_today"`
	NewUsersWeek    int64 `json:"new_users_week"`
	TotalAdmins     int64 `json:"total_admins"`
}
// ✅ Новые типы для трафика, конверсии, возвратов и прибыли
type TrafficStats struct {
	TotalVisits      int64   `json:"total_visits"`
	UniqueUsers      int64   `json:"unique_users"`
	AverageSession   float64 `json:"average_session"`
}

type ConversionStats struct {
	Visitors     int64   `json:"visitors"`
	Buyers       int64   `json:"buyers"`
	Conversion   float64 `json:"conversion_rate"`
}

type RefundStats struct {
	TotalRefunds    int64   `json:"total_refunds"`
	RefundedAmount  float64 `json:"refunded_amount"`
	RefundRate      float64 `json:"refund_rate"`
}

type ProfitStats struct {
	TotalRevenue float64 `json:"total_revenue"`
	TotalCost    float64 `json:"total_cost"`
	NetProfit    float64 `json:"net_profit"`
	ProfitMargin float64 `json:"profit_margin"`
}
// ✅ 1. Трафик
func (r *StatsRepository) GetTrafficStats() (TrafficStats, error) {
	db := database.DB
	var stats TrafficStats

	// Если у тебя нет таблицы traffic_logs, просто пример:
	db.Raw(`SELECT COUNT(*) FROM traffic_logs`).Scan(&stats.TotalVisits)
	db.Raw(`SELECT COUNT(DISTINCT user_id) FROM traffic_logs WHERE user_id IS NOT NULL`).Scan(&stats.UniqueUsers)
	db.Raw(`SELECT COALESCE(AVG(session_duration),0) FROM traffic_logs`).Scan(&stats.AverageSession)

	return stats, nil
}

// ✅ 2. Конверсия
func (r *StatsRepository) GetConversionStats() (ConversionStats, error) {
	db := database.DB
	var stats ConversionStats

	db.Table("users").Count(&stats.Visitors)
	db.Table("orders").Where("status = 'completed'").Distinct("user_id").Count(&stats.Buyers)

	if stats.Visitors > 0 {
		stats.Conversion = float64(stats.Buyers) / float64(stats.Visitors) * 100
	}

	return stats, nil
}

// ✅ 3. Возвраты
func (r *StatsRepository) GetRefundStats() (RefundStats, error) {
	db := database.DB
	var stats RefundStats

	db.Table("orders").Where("status = 'refunded'").Count(&stats.TotalRefunds)
	db.Table("orders").Select("COALESCE(SUM(total), 0)").Where("status = 'refunded'").Scan(&stats.RefundedAmount)

	var totalOrders int64
	db.Table("orders").Count(&totalOrders)

	if totalOrders > 0 {
		stats.RefundRate = float64(stats.TotalRefunds) / float64(totalOrders) * 100
	}

	return stats, nil
}

// ✅ 4. Прибыль
func (r *StatsRepository) GetProfitStats() (ProfitStats, error) {
	db := database.DB
	var stats ProfitStats

	db.Table("orders").Select("COALESCE(SUM(total), 0)").Where("status = 'completed'").Scan(&stats.TotalRevenue)
	db.Table("products").Select("COALESCE(SUM(cost_price * stock), 0)").Scan(&stats.TotalCost)

	stats.NetProfit = stats.TotalRevenue - stats.TotalCost
	if stats.TotalRevenue > 0 {
		stats.ProfitMargin = (stats.NetProfit / stats.TotalRevenue) * 100
	}

	return stats, nil
}

func (r *StatsRepository) GetAdminStats() (AdminStats, error) {
	db := database.DB
	var stats AdminStats

	if err := db.Table("users").Count(&stats.TotalUsers).Error; err != nil {
		return stats, err
	}
	if err := db.Table("orders").Count(&stats.TotalOrders).Error; err != nil {
		return stats, err
	}
	if err := db.Table("orders").Select("COALESCE(SUM(total),0)").Scan(&stats.TotalRevenue).Error; err != nil {
		return stats, err
	}
	if err := db.Table("products").Count(&stats.TotalProducts).Error; err != nil {
		return stats, err
	}
	if err := db.Table("reviews").Count(&stats.TotalReviews).Error; err != nil {
		return stats, err
	}

	return stats, nil
}

func (r *StatsRepository) GetMonthlyStats() ([]MonthlyStats, error) {
	db := database.DB
	var monthlyStats []MonthlyStats

	query := `
		WITH months AS (
			SELECT 
				TO_CHAR(date_trunc('month', CURRENT_DATE - INTERVAL '5 months' + (n || ' months')::interval), 'YYYY-MM') as month,
				date_trunc('month', CURRENT_DATE - INTERVAL '5 months' + (n || ' months')::interval) as month_start
			FROM generate_series(0, 5) n
		),
		order_stats AS (
			SELECT 
				TO_CHAR(date_trunc('month', created_at), 'YYYY-MM') as month,
				COUNT(*) as orders,
				COALESCE(SUM(total), 0) as revenue
			FROM orders 
			WHERE created_at >= (SELECT MIN(month_start) FROM months)
			GROUP BY date_trunc('month', created_at)
		),
		user_stats AS (
			SELECT 
				TO_CHAR(date_trunc('month', created_at), 'YYYY-MM') as month,
				COUNT(*) as new_users
			FROM users
			WHERE created_at >= (SELECT MIN(month_start) FROM months)
			GROUP BY date_trunc('month', created_at)
		)
		SELECT 
			m.month,
			COALESCE(os.orders, 0) as orders,
			COALESCE(os.revenue, 0) as revenue,
			COALESCE(us.new_users, 0) as new_users
		FROM months m
		LEFT JOIN order_stats os ON os.month = m.month
		LEFT JOIN user_stats us ON us.month = m.month
		ORDER BY m.month DESC
	`

	if err := db.Raw(query).Scan(&monthlyStats).Error; err != nil {
		return nil, err
	}

	return monthlyStats, nil
}

func (r *StatsRepository) GetPopularProducts(limit int) ([]ProductStats, error) {
	db := database.DB
	var popularProducts []ProductStats

	query := `
		SELECT 
			p.name as product_name,
			COUNT(oi.id) as sales_count
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		GROUP BY p.id, p.name
		ORDER BY sales_count DESC
		LIMIT ?
	`

	if err := db.Raw(query, limit).Scan(&popularProducts).Error; err != nil {
		return nil, err
	}

	return popularProducts, nil
}

func (r *StatsRepository) GetSalesByCategory() ([]CategoryStats, error) {
	db := database.DB
	var categoryStats []CategoryStats

	query := `
		SELECT 
			c.name as category_name,
			COALESCE(SUM(oi.quantity * oi.price), 0) as total_sales,
			COUNT(DISTINCT o.id) as order_count
		FROM categories c
		LEFT JOIN products p ON p.category_id = c.id
		LEFT JOIN order_items oi ON oi.product_id = p.id
		LEFT JOIN orders o ON o.id = oi.order_id AND o.status = 'completed'
		GROUP BY c.id, c.name
		ORDER BY total_sales DESC
	`

	if err := db.Raw(query).Scan(&categoryStats).Error; err != nil {
		return nil, err
	}

	return categoryStats, nil
}

func (r *StatsRepository) GetChartData() (ChartData, error) {
	var chartData ChartData
	var err error

	chartData.MonthlyStats, err = r.GetMonthlyStats()
	if err != nil {
		return chartData, err
	}

	chartData.PopularProducts, err = r.GetPopularProducts(10)
	if err != nil {
		return chartData, err
	}

	chartData.SalesByCategory, err = r.GetSalesByCategory()
	if err != nil {
		return chartData, err
	}

	return chartData, nil
}

func (r *StatsRepository) GetAdvancedStats() (AdvancedStats, error) {
	db := database.DB
	var stats AdvancedStats

	// Средний чек
	db.Table("orders").Where("status = 'completed'").
		Select("COALESCE(AVG(total), 0)").Scan(&stats.AverageOrderValue)

	// Повторные клиенты
	db.Raw(`
		SELECT COUNT(*) FROM (
			SELECT user_id FROM orders 
			WHERE status = 'completed' 
			GROUP BY user_id HAVING COUNT(*) > 1
		) as repeat_customers
	`).Scan(&stats.RepeatCustomers)

	// Lifetime Value (упрощенный расчет)
	db.Raw(`
		SELECT COALESCE(AVG(total_orders), 0) FROM (
			SELECT user_id, COUNT(*) as order_count, SUM(total) as total_orders 
			FROM orders WHERE status = 'completed' 
			GROUP BY user_id
		) user_orders
	`).Scan(&stats.CustomerLifetime)

	// Conversion Rate (упрощенный)
	var totalVisitors int64
	db.Table("users").Count(&totalVisitors)
	var buyers int64
	db.Table("orders").Where("status = 'completed'").Distinct("user_id").Count(&buyers)
	
	if totalVisitors > 0 {
		stats.ConversionRate = float64(buyers) / float64(totalVisitors) * 100
	}

	return stats, nil
}

func (r *StatsRepository) GetRealTimeStats() (RealTimeStats, error) {
	db := database.DB
	var stats RealTimeStats

	today := time.Now().Format("2006-01-02")

	db.Table("orders").Where("DATE(created_at) = ?", today).
		Count(&stats.TodayOrders)
	db.Table("orders").Where("DATE(created_at) = ?", today).
		Select("COALESCE(SUM(total), 0)").Scan(&stats.TodayRevenue)
	db.Table("users").Where("DATE(created_at) = ?", today).
		Count(&stats.TodayUsers)
	db.Table("orders").Where("status = 'pending'").
		Count(&stats.PendingOrders)

	return stats, nil
}

func (r *StatsRepository) GetUserStats() (UserStats, error) {
	db := database.DB
	var stats UserStats

	today := time.Now().Format("2006-01-02")
	weekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	// Активные пользователи (с заказами за последний месяц)
	db.Raw(`
		SELECT COUNT(DISTINCT user_id) FROM orders 
		WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
	`).Scan(&stats.ActiveUsers)

	db.Table("users").Where("DATE(created_at) = ?", today).
		Count(&stats.NewUsersToday)
	db.Table("users").Where("created_at >= ?", weekAgo).
		Count(&stats.NewUsersWeek)
	db.Table("users").Where("role = 'admin'").
		Count(&stats.TotalAdmins)

	return stats, nil
}

func (r *StatsRepository) GetDashboardStats() (map[string]interface{}, error) {
	dashboard := make(map[string]interface{})

	basicStats, err := r.GetAdminStats()
	if err != nil {
		return nil, err
	}

	realTimeStats, err := r.GetRealTimeStats()
	if err != nil {
		return nil, err
	}

	userStats, err := r.GetUserStats()
	if err != nil {
		return nil, err
	}

	advancedStats, err := r.GetAdvancedStats()
	if err != nil {
		return nil, err
	}

	chartData, err := r.GetChartData()
	if err != nil {
		return nil, err
	}

	dashboard["basic"] = basicStats
	dashboard["realtime"] = realTimeStats
	dashboard["users"] = userStats
	dashboard["advanced"] = advancedStats
	dashboard["charts"] = chartData

	return dashboard, nil
}

func (r *StatsRepository) GetUserPurchaseStats() ([]OrderStats, error) {
	db := database.DB
	var stats []OrderStats

	query := `
		SELECT 
			users.full_name AS user_name,
			COUNT(orders.id) AS total_orders,
			SUM(orders.total_amount) AS total_spent,
			MAX(orders.created_at) AS last_order_date,
			(
				SELECT payment_method
				FROM orders o2
				WHERE o2.user_id = users.id
				GROUP BY payment_method
				ORDER BY COUNT(*) DESC
				LIMIT 1
			) AS most_used_method
		FROM orders
		JOIN users ON orders.user_id = users.id
		GROUP BY users.id, users.full_name
		ORDER BY total_spent DESC;
	`
	result := db.Raw(query).Scan(&stats)
	return stats, result.Error
}
func (r *StatsRepository) SaveDailyStats(date time.Time) error {
	db := database.DB

	var totalUsers, totalOrders, totalProducts int64
	var totalRevenue float64

	db.Model(&models.User{}).Count(&totalUsers)
	db.Model(&models.Order{}).Count(&totalOrders)
	db.Model(&models.Product{}).Count(&totalProducts)
	db.Model(&models.Order{}).Select("SUM(total_amount)").Scan(&totalRevenue)

	daily := models.DailyStats{
		Date:          date.Format("2006-01-02"),
		TotalUsers:    totalUsers,
		TotalOrders:   totalOrders,
		TotalProducts: totalProducts,
		TotalRevenue:  totalRevenue,
	}

	return db.Create(&daily).Error
}
