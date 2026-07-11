package animals

import (
	"errors"
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
	herd.CurrentHeadCount = herd.InitialHeadCount
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
	return s.DB.Model(&animalModels.Herd{}).
		Where("id = ? AND user_id = ?", id, userID).
		Select("Name", "AnimalTypeID", "Location").
		Updates(herd).Error
}

func (s *HerdService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&animalModels.Herd{}).Error
}

func (s *HerdService) RecordActivity(userID uint, herdID uint, activity *animalModels.HerdActivity) error {
	if activity.Count <= 0 {
		return errors.New("count must be greater than zero")
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var herd animalModels.Herd
		if err := tx.Where("id = ? AND user_id = ?", herdID, userID).First(&herd).Error; err != nil {
			return err
		}

		switch activity.ActivityType {
		case "fatality":
			if herd.CurrentHeadCount-activity.Count < 0 {
				return errors.New("cannot record more fatalities than current headcount")
			}
			herd.CurrentHeadCount -= activity.Count
		case "birth":
			herd.CurrentHeadCount += activity.Count
		default:
			return errors.New("invalid activity type")
		}

		activity.HerdID = herd.ID
		if err := tx.Create(activity).Error; err != nil {
			return err
		}

		return tx.Save(&herd).Error
	})
}
