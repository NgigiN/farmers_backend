package services

import (
	"farm-backend/internal/middleware"
	summariesModels "farm-backend/internal/models/summaries"

	"gorm.io/gorm"
)

type RevenueService struct {
	DB *gorm.DB
}

func NewRevenueService(db *gorm.DB) *RevenueService {
	return &RevenueService{DB: db}
}

func (s *RevenueService) Create(UserID uint, revenue *summariesModels.Revenue) error {
	revenue.UserID = UserID
	if err := middleware.ValidateStruct(revenue); err != nil {
		return err
	}
	return s.DB.Create(revenue).Error
}

func (s *RevenueService) List(UserID uint) ([]summariesModels.Revenue, error) {
	var revenues []summariesModels.Revenue
	return revenues, s.DB.Where("user_id = ?", UserID).Find(&revenues).Error
}

func (s *RevenueService) Get(UserID uint, id uint) (*summariesModels.Revenue, error) {
	var revenue summariesModels.Revenue
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).First(&revenue).Error
	if err != nil {
		return nil, err
	}
	return &revenue, nil
}
func (s *RevenueService) Update(userID, id uint, revenue *summariesModels.Revenue) error {
	if err := middleware.ValidateStruct(revenue); err != nil {
		return err
	}
	return s.DB.Model(&summariesModels.Revenue{}).Where("id = ? AND user_id = ?", id, userID).Updates(revenue).Error
}

func (s *RevenueService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(summariesModels.Revenue{}).Error
}

func (s *RevenueService) GetTotalRevenue(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("revenues").Where("user_id = ?", UserID).Select("SUM(total) as total").Scan(&total).Error
	return total, err
}

func (s *RevenueService) GetTotalRevenueBySource(UserID uint, source string) (float64, error) {
	var total float64
	err := s.DB.Table("revenues").Where("user_id = ? AND source = ?", UserID, source).Select("SUM(total) as total").Scan(&total).Error
	return total, err
}
