package animals

import (
	"farm-backend/internal/middleware"
	animalModels "farm-backend/internal/models/animals"

	"gorm.io/gorm"
)

type HerdService struct {
	DB *gorm.DB
}

func NewHerdService(db *gorm.DB) *HerdService {
	return &HerdService{DB: db}
}

func (s *HerdService) Create(UserID uint, herd *animalModels.Herd) error {
	herd.UserID = UserID
	if err := s.DB.Where("id = ? AND user_id = ?", herd.AnimalTypeID, UserID).First(&animalModels.AnimalType{}).Error; err != nil {
		return err
	}
	if err := middleware.ValidateStruct(herd); err != nil {
		return err
	}
	return s.DB.Create(herd).Error
}

func (s *HerdService) List(UserID uint) ([]animalModels.Herd, error) {
	var herds []animalModels.Herd
	return herds, s.DB.Where("user_id = ?", UserID).Find(&herds).Error
}

func (s *HerdService) Get(UserID uint, id uint) (*animalModels.Herd, error) {
	var herd animalModels.Herd
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).First(&herd).Error
	if err != nil {
		return nil, err
	}
	return &herd, nil
}

func (s *HerdService) Update(userID, id uint, herd *animalModels.Herd) error {
	if err := middleware.ValidateStruct(herd); err != nil {
		return err
	}
	return s.DB.Model(&animalModels.Herd{}).Where("id = ? AND user_id = ?", id, userID).Updates(herd).Error
}

func (s *HerdService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(animalModels.Herd{}).Error
}
