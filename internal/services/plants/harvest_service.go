package plants

import (
	"errors"
	"farm-backend/internal/middleware"
	plantModels "farm-backend/internal/models/plants"
	"strings"

	"gorm.io/gorm"
)

type HarvestService struct {
	DB *gorm.DB
}

func NewHarvestService(db *gorm.DB) *HarvestService {
	return &HarvestService{DB: db}
}

func (s *HarvestService) validateHarvest(userID uint, harvest *plantModels.Harvest) error {
	if harvest.SeasonID == 0 {
		return errors.New("season_id is required")
	}
	if harvest.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	if strings.TrimSpace(harvest.Unit) == "" {
		return errors.New("unit is required")
	}
	harvest.Unit = strings.TrimSpace(harvest.Unit)

	if err := s.DB.Where("id = ? AND user_id = ?", harvest.SeasonID, userID).
		First(&plantModels.Season{}).Error; err != nil {
		return errors.New("season not found or does not belong to user")
	}

	return middleware.ValidateStruct(harvest)
}

func (s *HarvestService) Create(userID uint, harvest *plantModels.Harvest) error {
	harvest.UserID = userID
	if err := s.validateHarvest(userID, harvest); err != nil {
		return err
	}
	return s.DB.Create(harvest).Error
}

func (s *HarvestService) List(userID uint, seasonID uint) ([]plantModels.Harvest, error) {
	var harvests []plantModels.Harvest
	q := s.DB.Where("user_id = ?", userID)
	if seasonID != 0 {
		q = q.Where("season_id = ?", seasonID)
	}
	return harvests, q.Order("date DESC").Find(&harvests).Error
}

func (s *HarvestService) Get(userID uint, id uint) (*plantModels.Harvest, error) {
	var harvest plantModels.Harvest
	err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&harvest).Error
	if err != nil {
		return nil, err
	}
	return &harvest, nil
}

func (s *HarvestService) Update(userID, id uint, harvest *plantModels.Harvest) error {
	if err := s.validateHarvest(userID, harvest); err != nil {
		return err
	}
	return s.DB.Model(&plantModels.Harvest{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(harvest).Error
}

func (s *HarvestService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).
		Delete(&plantModels.Harvest{}).Error
}