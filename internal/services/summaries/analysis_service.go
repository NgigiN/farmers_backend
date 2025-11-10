package services

import (
	"gorm.io/gorm"
)

type AnalysisService struct {
	DB *gorm.DB
}

func NewAnalysisService(db *gorm.DB) *AnalysisService {
	return &AnalysisService{DB: db}
}

// Total revenue
func (s *AnalysisService) GetTotalRevenue(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("revenues").Where("user_id = ?", UserID).Select("SUM(total) as total").Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetPlantRevenue(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("revenues").Where("user_id = ? AND type = 'Plant'").Select("SUM(total) as total").Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetAnimalRevenue(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("revenues").Where("user_id = ? AND type = 'Animal'").Select("SUM(total) as total").Scan(&total).Error
	return total, err
}

// Total costs
func (s *AnalysisService) GetTotalPlantCosts(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("inputs").Where("user_id = ?", UserID).Select("SUM(cost) as total").Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetTotalAnimalCosts(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("animal_inputs").Where("user_id = ?", UserID).Select("SUM(cost) as total").Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetTotalCosts(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("inputs").Where("user_id = ?", UserID).Select("SUM(cost) as total").Scan(&total).Error
	return total, err
}

// profit analysis
func (s *AnalysisService) GetProfit(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("revenues").Where("user_id = ?", UserID).Select("SUM(total) as total").Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetPlantProfit(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("revenues").Where("user_id = ? AND type = 'Plant'").Select("SUM(total) as total").Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetAnimalProfit(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("revenues").Where("user_id = ? AND type = 'Animal'").Select("SUM(total) as total").Scan(&total).Error
	return total, err
}

// Breakdowns
func (s *AnalysisService) GetCostBreakdownByCategory(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("inputs").Where("user_id = ?").Select("SUM(cost) as total").Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetCostBreakdownByInputType(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("inputs").Where("user_id = ?").Select("SUM(cost) as total").Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetCostBreakdownBySource(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("inputs").Where("user_id = ?").Select("SUM(cost) as total").Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetMonthlySummary(UserID uint) ([]MonthlySummary, error) {
	var results []MonthlySummary
	err := s.DB.Table("revenues").Where("user_id = ?").Select("strftime('%Y-%m', date) as month, SUM(total) as total").Group("strftime('%Y-%m', date)").Scan(&results).Error
	return results, err
}

type MonthlySummary struct {
	Month string
	Total float64
}
