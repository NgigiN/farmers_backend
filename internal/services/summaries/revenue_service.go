package summaries

import (
	"errors"
	"time"

	animalModels "farm-backend/internal/models/animals"
	plantModels "farm-backend/internal/models/plants"
	summariesModels "farm-backend/internal/models/summaries"

	"gorm.io/gorm"
)

type RevenueService struct {
	DB *gorm.DB
}

func NewRevenueService(db *gorm.DB) *RevenueService {
	return &RevenueService{DB: db}
}

func (s *RevenueService) validateRevenue(userID uint, revenue *summariesModels.Revenue) error {
	if revenue.Source != "plant" && revenue.Source != "animal" {
		return errors.New("source must be either 'plant' or 'animal'")
	}
	if revenue.Quantity <= 0 || revenue.UnitPrice <= 0 {
		return errors.New("quantity and unit_price must be greater than zero")
	}
	if revenue.Source == "plant" {
		if err := s.DB.Where("id = ? AND user_id = ?", revenue.SourceID, userID).
			First(&plantModels.Season{}).Error; err != nil {
			return errors.New("season not found or does not belong to user")
		}
	} else {
		if err := s.DB.Where("id = ? AND user_id = ?", revenue.SourceID, userID).
			First(&animalModels.Herd{}).Error; err != nil {
			return errors.New("herd not found or does not belong to user")
		}
	}
	return nil
}

func (s *RevenueService) Create(UserID uint, revenue *summariesModels.Revenue) error {
	revenue.UserID = UserID
	if revenue.Total == 0 && revenue.Quantity > 0 && revenue.UnitPrice > 0 {
		revenue.Total = revenue.Quantity * revenue.UnitPrice
	}
	if err := s.validateRevenue(UserID, revenue); err != nil {
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
	if revenue.Total == 0 && revenue.Quantity > 0 && revenue.UnitPrice > 0 {
		revenue.Total = revenue.Quantity * revenue.UnitPrice
	}
	if err := s.validateRevenue(userID, revenue); err != nil {
		return err
	}
	return s.DB.Model(&summariesModels.Revenue{}).
		Where("id = ? AND user_id = ?", id, userID).
		Select("Source", "SourceID", "Type", "Quantity", "UnitPrice", "Total", "Date", "Notes").
		Updates(revenue).Error
}

func (s *RevenueService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&summariesModels.Revenue{}).Error
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

func (s *RevenueService) ListBySource(UserID uint, source string) ([]summariesModels.Revenue, error) {
	var revenues []summariesModels.Revenue
	return revenues, s.DB.Where("user_id = ? AND source = ?", UserID, source).Find(&revenues).Error
}

func (s *RevenueService) ListByDateRange(UserID uint, startDate, endDate time.Time) ([]summariesModels.Revenue, error) {
	var revenues []summariesModels.Revenue
	return revenues, s.DB.Where("user_id = ? AND date >= ? AND date <= ?", UserID, startDate, endDate).Find(&revenues).Error
}