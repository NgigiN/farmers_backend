package services

import (
	"farm-backend/internal/middleware"
	"farm-backend/internal/models"

	"gorm.io/gorm"
)

type CropService struct {
	DB *gorm.DB
}

func NewCropService(db *gorm.DB) *CropService {
	return &CropService{DB: db}
}

func (s *CropService) Create(UserID uint, crop *models.Crop) error {
	crop.UserID = UserID
	if err := middleware.ValidateStruct(crop); err != nil {
		return err
	}
	return s.DB.Create(crop).Error
}

func (s *CropService) List(UserID uint) ([]models.Crop, error) {
	var crops []models.Crop
	return crops, s.DB.Where("user_id = ?", UserID).Find(&crops).Error
}

func (s *CropService) Get(UserID uint, id uint) (*[]models.Crop, error) {
	var crops []models.Crop
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).Find(&crops).Error
	return &crops, err
}

func (s *CropService) Update(userID, id uint, crop *models.Crop) error {
	if err := middleware.ValidateStruct(crop); err != nil {
		return err
	}
	return s.DB.Model(&models.Crop{}).Where("id = ? AND user_id = ?", id, userID).Updates(crop).Error
}

func (s *CropService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(models.Crop{}).Error
}
