package plants

import (
	"errors"
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

func (s *ActivityService) validateActivity(userID uint, activity *activityModels.Activity) error {
	if activity.SourceType != "plant" && activity.SourceType != "animal" {
		return errors.New("source_type must be either 'plant' or 'animal'")
	}

	if activity.SourceType == "plant" {
		if err := s.DB.Where("id = ? AND user_id = ?", activity.SourceID, userID).First(&activityModels.Season{}).Error; err != nil {
			return errors.New("season not found or does not belong to user")
		}
	} else {
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
	if err := s.DB.Where("user_id = ? AND name = ? AND type = ? AND category = ?",
		userID, activity.Type, activity.SourceType, "activity").First(&category).Error; err != nil {
		return errors.New("activity type does not exist in cost categories — please create it first")
	}

	return nil
}

func (s *ActivityService) Create(userID uint, activity *activityModels.Activity) error {
	activity.UserID = userID
	if err := s.validateActivity(userID, activity); err != nil {
		return err
	}
	return s.DB.Create(activity).Error
}

// List returns all activities for a user, optionally filtered by source_type at the DB level.
func (s *ActivityService) List(userID uint, sourceType string) ([]activityModels.Activity, error) {
	var activities []activityModels.Activity
	q := s.DB.Where("user_id = ?", userID)
	if sourceType != "" {
		q = q.Where("source_type = ?", sourceType)
	}
	return activities, q.Find(&activities).Error
}

func (s *ActivityService) Get(userID uint, id uint) (*activityModels.Activity, error) {
	var activity activityModels.Activity
	err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&activity).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (s *ActivityService) Update(userID, id uint, activity *activityModels.Activity) error {
	if err := s.validateActivity(userID, activity); err != nil {
		return err
	}
	return s.DB.Model(&activityModels.Activity{}).
		Where("id = ? AND user_id = ?", id, userID).
		Select("SourceType", "SourceID", "AnimalID", "Type", "Details", "Cost", "Date", "Notes").
		Updates(activity).Error
}

func (s *ActivityService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&activityModels.Activity{}).Error
}
