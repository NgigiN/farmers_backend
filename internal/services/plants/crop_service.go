package services

import (
	"farm-backend/internal/middleware"
	plantModels "farm-backend/internal/models/plants"

	"gorm.io/gorm"
)

type CropService struct {
	DB *gorm.DB
}

func NewCropService(db *gorm.DB) *CropService {
	return &CropService{DB: db}
}

func (s *CropService) Create(UserID uint, crop *plantModels.Crop) error {
	crop.UserID = UserID
	if err := middleware.ValidateStruct(crop); err != nil {
		return err
	}
	return s.DB.Create(crop).Error
}

func (s *CropService) List(UserID uint) ([]plantModels.Crop, error) {
	var crops []plantModels.Crop
	return crops, s.DB.Where("user_id = ?", UserID).Find(&crops).Error
}

func (s *CropService) Get(UserID uint, id uint) (*plantModels.Crop, error) {
	var crop plantModels.Crop
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).First(&crop).Error
	if err != nil {
		return nil, err
	}
	return &crop, nil
}

func (s *CropService) Update(userID, id uint, crop *plantModels.Crop) error {
	if err := middleware.ValidateStruct(crop); err != nil {
		return err
	}
	return s.DB.Model(&plantModels.Crop{}).Where("id = ? AND user_id = ?", id, userID).Updates(crop).Error
}

func (s *CropService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(plantModels.Crop{}).Error
}
