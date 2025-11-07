package animals

import (
	"farm-backend/internal/middleware"
	animalModels "farm-backend/internal/models/animals"

	"gorm.io/gorm"
)

type AnimalInputService struct {
	DB *gorm.DB
}

func NewAnimalInputService(db *gorm.DB) *AnimalInputService {
	return &AnimalInputService{DB: db}
}

func (s *AnimalInputService) Create(UserID uint, animalInput *animalModels.AnimalInput) error {
	if err := s.DB.Where("id = ? AND user_id = ?", animalInput.HerdID, UserID).First(&animalModels.Herd{}).Error; err != nil {
		return err
	}
	if err := middleware.ValidateStruct(animalInput); err != nil {
		return err
	}
	return s.DB.Create(animalInput).Error
}

func (s *AnimalInputService) List(UserID uint) ([]animalModels.AnimalInput, error) {
	var animalInputs []animalModels.AnimalInput
	return animalInputs, s.DB.Table("animal_inputs").
		Joins("JOIN herds ON animal_inputs.herd_id = herds.id").
		Where("herds.user_id = ?", UserID).
		Find(&animalInputs).Error
}

func (s *AnimalInputService) Get(UserID uint, id uint) (*animalModels.AnimalInput, error) {
	var animalInput animalModels.AnimalInput
	err := s.DB.Table("animal_inputs").
		Joins("JOIN herds ON animal_inputs.herd_id = herds.id").
		Where("animal_inputs.id = ? AND herds.user_id = ?", id, UserID).
		First(&animalInput).Error
	if err != nil {
		return nil, err
	}
	return &animalInput, nil
}

func (s *AnimalInputService) Update(userID, id uint, animalInput *animalModels.AnimalInput) error {
	if err := s.DB.Where("id = ? AND user_id = ?", animalInput.HerdID, userID).First(&animalModels.Herd{}).Error; err != nil {
		return err
	}
	if err := middleware.ValidateStruct(animalInput); err != nil {
		return err
	}
	return s.DB.Table("animal_inputs").
		Joins("JOIN herds ON animal_inputs.herd_id = herds.id").
		Where("animal_inputs.id = ? AND herds.user_id = ?", id, userID).
		Updates(animalInput).Error
}

func (s *AnimalInputService) Delete(userID, id uint) error {
	return s.DB.Table("animal_inputs").
		Joins("JOIN herds ON animal_inputs.herd_id = herds.id").
		Where("animal_inputs.id = ? AND herds.user_id = ?", id, userID).
		Delete(&animalModels.AnimalInput{}).Error
}
