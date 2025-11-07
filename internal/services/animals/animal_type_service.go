package animals

import (
	"farm-backend/internal/middleware"
	animalModels "farm-backend/internal/models/animals"

	"gorm.io/gorm"
)

type AnimalTypeService struct {
	DB *gorm.DB
}

func NewAnimalTypeService(db *gorm.DB) *AnimalTypeService {
	return &AnimalTypeService{DB: db}
}

func (s *AnimalTypeService) Create(UserID uint, animalType *animalModels.AnimalType) error {
	animalType.UserID = UserID
	if err := middleware.ValidateStruct(animalType); err != nil {
		return err
	}
	return s.DB.Create(animalType).Error
}

func (s *AnimalTypeService) List(UserID uint) ([]animalModels.AnimalType, error) {
	var animalTypes []animalModels.AnimalType
	return animalTypes, s.DB.Where("user_id = ?", UserID).Find(&animalTypes).Error
}

func (s *AnimalTypeService) Get(UserID uint, id uint) (*animalModels.AnimalType, error) {
	var animalType animalModels.AnimalType
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).First(&animalType).Error
	if err != nil {
		return nil, err
	}
	return &animalType, nil
}

func (s *AnimalTypeService) Update(userID, id uint, animalType *animalModels.AnimalType) error {
	if err := middleware.ValidateStruct(animalType); err != nil {
		return err
	}
	return s.DB.Model(&animalModels.AnimalType{}).Where("id = ? AND user_id = ?", id, userID).Updates(animalType).Error
}

func (s *AnimalTypeService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(animalModels.AnimalType{}).Error
}
