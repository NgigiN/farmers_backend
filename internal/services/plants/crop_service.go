package plants

import (
	"farm-backend/internal/middleware"
	plantModels "farm-backend/internal/models/plants"

	"gorm.io/gorm"
)

type PlantService struct {
	DB *gorm.DB
}

func NewPlantService(db *gorm.DB) *PlantService {
	return &PlantService{DB: db}
}

func (s *PlantService) Create(UserID uint, plant *plantModels.Plant) error {
	plant.UserID = UserID
	if err := middleware.ValidateStruct(plant); err != nil {
		return err
	}
	return s.DB.Create(plant).Error
}

func (s *PlantService) List(UserID uint) ([]plantModels.Plant, error) {
	var plants []plantModels.Plant
	return plants, s.DB.Where("user_id = ?", UserID).Find(&plants).Error
}

func (s *PlantService) Get(UserID uint, id uint) (*plantModels.Plant, error) {
	var plant plantModels.Plant
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).First(&plant).Error
	if err != nil {
		return nil, err
	}
	return &plant, nil
}

func (s *PlantService) Update(userID, id uint, plant *plantModels.Plant) error {
	if err := middleware.ValidateStruct(plant); err != nil {
		return err
	}
	return s.DB.Model(&plantModels.Plant{}).Where("id = ? AND user_id = ?", id, userID).Updates(plant).Error
}

func (s *PlantService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&plantModels.Plant{}).Error
}
