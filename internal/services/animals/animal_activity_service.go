package animals

import (
	"farm-backend/internal/middleware"
	animalModels "farm-backend/internal/models/animals"

	"gorm.io/gorm"
)

type AnimalActivityService struct {
	DB *gorm.DB
}

func NewAnimalActivityService(db *gorm.DB) *AnimalActivityService {
	return &AnimalActivityService{DB: db}
}

func (s *AnimalActivityService) Create(UserID uint, animalActivity *animalModels.AnimalActivity) error {
	if err := s.DB.Where("id = ? AND user_id = ?", animalActivity.HerdID, UserID).First(&animalModels.Herd{}).Error; err != nil {
		return err
	}
	if err := middleware.ValidateStruct(animalActivity); err != nil {
		return err
	}
	return s.DB.Create(animalActivity).Error
}

func (s *AnimalActivityService) List(UserID uint) ([]animalModels.AnimalActivity, error) {
	var animalActivities []animalModels.AnimalActivity
	return animalActivities, s.DB.Table("animal_activities").
		Joins("JOIN herds ON animal_activities.herd_id = herds.id").
		Where("herds.user_id = ?", UserID).
		Find(&animalActivities).Error
}

func (s *AnimalActivityService) Get(UserID uint, id uint) (*animalModels.AnimalActivity, error) {
	var animalActivity animalModels.AnimalActivity
	err := s.DB.Table("animal_activities").
		Joins("JOIN herds ON animal_activities.herd_id = herds.id").
		Where("animal_activities.id = ? AND herds.user_id = ?", id, UserID).
		First(&animalActivity).Error
	if err != nil {
		return nil, err
	}
	return &animalActivity, nil
}

func (s *AnimalActivityService) Update(userID, id uint, animalActivity *animalModels.AnimalActivity) error {
	if err := s.DB.Where("id = ? AND user_id = ?", animalActivity.HerdID, userID).First(&animalModels.Herd{}).Error; err != nil {
		return err
	}
	if err := middleware.ValidateStruct(animalActivity); err != nil {
		return err
	}
	return s.DB.Table("animal_activities").
		Joins("JOIN herds ON animal_activities.herd_id = herds.id").
		Where("animal_activities.id = ? AND herds.user_id = ?", id, userID).
		Updates(animalActivity).Error
}

func (s *AnimalActivityService) Delete(userID, id uint) error {
	return s.DB.Table("animal_activities").
		Joins("JOIN herds ON animal_activities.herd_id = herds.id").
		Where("animal_activities.id = ? AND herds.user_id = ?", id, userID).
		Delete(&animalModels.AnimalActivity{}).Error
}
