package services

import (
	"farm-backend/internal/middleware"
	"farm-backend/internal/models"

	"gorm.io/gorm"
)

type SeasonService struct {
	DB *gorm.DB
}

func NewSeasonService(db *gorm.DB) *SeasonService {
	return &SeasonService{DB: db}
}

func (s *SeasonService) Create(UserID uint, season *models.Season) error {
	season.UserID = UserID
	if err := middleware.ValidateStruct(season); err != nil {
		return err
	}
	return s.DB.Create(season).Error
}

func (s *SeasonService) List(UserID uint) ([]models.Season, error) {
	var seasons []models.Season
	return seasons, s.DB.Where("user_id = ?", UserID).Find(&seasons).Error
}

func (s *SeasonService) Get(UserID uint, id uint) (*[]models.Season, error) {
	var seasons []models.Season
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).Find(&seasons).Error
	return &seasons, err
}

func (s *SeasonService) Update(userID, id uint, season *models.Season) error {
	if err := middleware.ValidateStruct(season); err != nil {
		return err
	}
	return s.DB.Model(&models.season{}).Where("id = ? AND user_id = ?", id, userID).Updates(activity).Error
}

func (s *SeasonService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(models.Season{}).Error
}
