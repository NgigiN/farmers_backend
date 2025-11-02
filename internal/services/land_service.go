package services

import (
	"farm-backend/internal/middleware"
	"farm-backend/internal/models"

	"gorm.io/gorm"
)

type LandService struct {
	DB *gorm.DB
}

func NewLandService(db *gorm.DB) *LandService {
	return &LandService{DB: db}
}

func (s *LandService) Create(UserID uint, land *models.Land) error {
	land.UserID = UserID
	if err := middleware.ValidateStruct(land); err != nil {
		return err
	}
	return s.DB.Create(land).Error
}

func (s *LandService) List(UserID uint) ([]models.Land, error) {
	var land []models.Land
	return land, s.DB.Where("user_id = ?", UserID).Find(&land).Error
}

func (s *LandService) Get(UserID uint, id uint) (*[]models.Land, error) {
	var land []models.Land
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).Find(&land).Error
	return &land, err
}

func (s *LandService) Update(userID, id uint, crop *models.Land) error {
	if err := middleware.ValidateStruct(crop); err != nil {
		return err
	}
	return s.DB.Model(&models.Land{}).Where("id = ? AND user_id = ?", id, userID).Updates(crop).Error
}

func (s *LandService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(models.Land{}).Error
}
