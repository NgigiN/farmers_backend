package services

import (
	"errors"
	"farm-backend/internal/middleware"
	animalModels "farm-backend/internal/models/animals"
	activityModels "farm-backend/internal/models/plants"
	summaryModels "farm-backend/internal/models/summaries"

	"gorm.io/gorm"
)

type ActivityService struct {
	DB *gorm.DB
}

func NewActivityService(db *gorm.DB) *ActivityService {
	return &ActivityService{DB: db}
}

func (s *ActivityService) Create(UserID uint, activity *activityModels.Activity) error {
	activity.UserID = UserID

	if activity.SourceType != "plant" && activity.SourceType != "animal" {
		return errors.New("source_type must be either 'plant' or 'animal'")
	}

	if activity.SourceType == "plant" {
		if err := s.DB.Where("id = ? AND user_id = ?", activity.SourceID, UserID).First(&activityModels.Season{}).Error; err != nil {
			return errors.New("season not found or does not belong to user")
		}
	} else if activity.SourceType == "animal" {
		if err := s.DB.Where("id = ? AND user_id = ?", activity.SourceID, UserID).First(&animalModels.Herd{}).Error; err != nil {
			return errors.New("herd not found or does not belong to user")
		}
		if activity.AnimalID != 0 {
			if err := s.DB.Where("id = ? AND user_id = ?", activity.AnimalID, UserID).First(&animalModels.Animal{}).Error; err != nil {
				return errors.New("animal not found or does not belong to user")
			}
		}
	}

	var category summaryModels.CostCategory
	if err := s.DB.Where("user_id = ? AND name = ? AND type = ? AND category = ?", UserID, activity.Type, activity.SourceType, "activity").First(&category).Error; err != nil {
		return errors.New("activity type does not exist in cost categories. Please create it first")
	}

	if err := middleware.ValidateStruct(activity); err != nil {
		return err
	}
	return s.DB.Create(activity).Error
}

func (s *ActivityService) List(UserID uint) ([]activityModels.Activity, error) {
	var activities []activityModels.Activity
	return activities, s.DB.Where("user_id = ?", UserID).Find(&activities).Error
}

func (s *ActivityService) Get(UserID uint, id uint) (*activityModels.Activity, error) {
	var activity activityModels.Activity
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).First(&activity).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (s *ActivityService) Update(userID, id uint, activity *activityModels.Activity) error {
	if activity.SourceType != "plant" && activity.SourceType != "animal" {
		return errors.New("source_type must be either 'plant' or 'animal'")
	}

	if activity.SourceType == "plant" {
		if err := s.DB.Where("id = ? AND user_id = ?", activity.SourceID, userID).First(&activityModels.Season{}).Error; err != nil {
			return errors.New("season not found or does not belong to user")
		}
	} else if activity.SourceType == "animal" {
		if err := s.DB.Where("id = ? AND user_id = ?", activity.SourceID, userID).First(&animalModels.Herd{}).Error; err != nil {
			return errors.New("herd not found or does not belong to user")
		}
		if activity.AnimalID != 0 {
			if err := s.DB.Where("id = ? AND user_id = ?", activity.AnimalID, userID).First(&animalModels.Animal{}).Error; err != nil {
				return errors.New("animal not found or does not belong to user")
			}
		}
	}

	var category summaryModels.CostCategory
	if err := s.DB.Where("user_id = ? AND name = ? AND type = ? AND category = ?", userID, activity.Type, activity.SourceType, "activity").First(&category).Error; err != nil {
		return errors.New("activity type does not exist in cost categories. Please create it first")
	}

	if err := middleware.ValidateStruct(activity); err != nil {
		return err
	}
	return s.DB.Model(&activityModels.Activity{}).Where("id = ? AND user_id = ?", id, userID).Updates(activity).Error
}

func (s *ActivityService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(activityModels.Activity{}).Error
}
