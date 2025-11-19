package animals

import (
	"farm-backend/internal/middleware"
	animalModels "farm-backend/internal/models/animals"

	"gorm.io/gorm"
)

type AnimalService struct {
	DB *gorm.DB
}

func NewAnimalService(db *gorm.DB) *AnimalService {
	return &AnimalService{DB: db}
}

func (s *AnimalService) Create(UserID uint, animal *animalModels.Animal) error {
	animal.UserID = UserID
	if err := s.DB.Where("id = ? AND user_id = ?", animal.AnimalTypeID, UserID).First(&animalModels.AnimalType{}).Error; err != nil {
		return err
	}
	if animal.HerdID != 0 {
		if err := s.DB.Where("id = ? AND user_id = ?", animal.HerdID, UserID).First(&animalModels.Herd{}).Error; err != nil {
			return err
		}
	}
	if err := middleware.ValidateStruct(animal); err != nil {
		return err
	}
	return s.DB.Create(animal).Error
}

func (s *AnimalService) List(UserID uint) ([]animalModels.Animal, error) {
	var animals []animalModels.Animal
	return animals, s.DB.Where("user_id = ?", UserID).Find(&animals).Error
}

func (s *AnimalService) Get(UserID uint, id uint) (*animalModels.Animal, error) {
	var animal animalModels.Animal
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).First(&animal).Error
	if err != nil {
		return nil, err
	}
	return &animal, nil
}

func (s *AnimalService) Update(userID, id uint, animal *animalModels.Animal) error {
	if err := s.DB.Where("id = ? AND user_id = ?", animal.AnimalTypeID, userID).First(&animalModels.AnimalType{}).Error; err != nil {
		return err
	}
	if animal.HerdID != 0 {
		if err := s.DB.Where("id = ? AND user_id = ?", animal.HerdID, userID).First(&animalModels.Herd{}).Error; err != nil {
			return err
		}
	}
	if err := middleware.ValidateStruct(animal); err != nil {
		return err
	}
	return s.DB.Model(&animalModels.Animal{}).Where("id = ? AND user_id = ?", id, userID).Updates(animal).Error
}

func (s *AnimalService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&animalModels.Animal{}).Error
}
