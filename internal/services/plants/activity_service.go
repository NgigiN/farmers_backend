package services

import (
	"farm-backend/internal/middleware"
	activityModels "farm-backend/internal/models/plants"

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
	if err := middleware.ValidateStruct(activity); err != nil {
		return err
	}
	return s.DB.Model(&activityModels.Activity{}).Where("id = ? AND user_id = ?", id, userID).Updates(activity).Error
}

func (s *ActivityService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(activityModels.Activity{}).Error
}
