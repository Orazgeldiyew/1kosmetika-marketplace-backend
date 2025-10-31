package services

import "1kosmetika-marketplace-backend/repositories"

type StatsService struct {
	repo *repositories.StatsRepository
}

func NewStatsService(repo *repositories.StatsRepository) *StatsService {
	return &StatsService{repo: repo}
}

func (s *StatsService) GetAdminStats() (repositories.AdminStats, error) {
	return s.repo.GetAdminStats()
}

func (s *StatsService) GetMonthlyStats() ([]repositories.MonthlyStats, error) {
	return s.repo.GetMonthlyStats()
}

func (s *StatsService) GetPopularProducts(limit int) ([]repositories.ProductStats, error) {
	return s.repo.GetPopularProducts(limit)
}

func (s *StatsService) GetSalesByCategory() ([]repositories.CategoryStats, error) {
	return s.repo.GetSalesByCategory()
}

func (s *StatsService) GetChartData() (repositories.ChartData, error) {
	return s.repo.GetChartData()
}

func (s *StatsService) GetAdvancedStats() (repositories.AdvancedStats, error) {
	return s.repo.GetAdvancedStats()
}

func (s *StatsService) GetRealTimeStats() (repositories.RealTimeStats, error) {
	return s.repo.GetRealTimeStats()
}

func (s *StatsService) GetUserStats() (repositories.UserStats, error) {
	return s.repo.GetUserStats()
}

func (s *StatsService) GetDashboardStats() (map[string]interface{}, error) {
	return s.repo.GetDashboardStats()
}

func (s *StatsService) GetUserPurchaseStats() ([]repositories.OrderStats, error) {
	return s.repo.GetUserPurchaseStats()
}



func (s *StatsService) GetTrafficStats() (repositories.TrafficStats, error) {
	return s.repo.GetTrafficStats()
}


func (s *StatsService) GetConversionStats() (repositories.ConversionStats, error) {
	return s.repo.GetConversionStats()
}


func (s *StatsService) GetRefundStats() (repositories.RefundStats, error) {
	return s.repo.GetRefundStats()
}


func (s *StatsService) GetProfitStats() (repositories.ProfitStats, error) {
	return s.repo.GetProfitStats()
}
