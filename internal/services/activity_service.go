package services

import (
	"farm-backend/internal/middleware"
	"farm-backend/internal/models"

	"gorm.io/gorm"
)

type ActivityService struct {
	DB *gorm.DB
}

func NewActivityService(db *gorm.DB) *ActivityService {
	return &ActivityService{DB: db}
}

func (s *ActivityService) Create(UserID uint, activity *models.Activity) error {
	activity.UserID = UserID
	if err := middleware.ValidateStruct(activity); err != nil {
		return err
	}
	return s.DB.Create(activity).Error
}

func (s *ActivityService) List(UserID uint) ([]models.Activity, error) {
	var activities []models.Activity
	return activities, s.DB.Where("user_id = ?", UserID).Find(&activities).Error
}

func (s *ActivityService) Get(UserID uint, id uint) (*[]models.Activity, error) {
	var activities []models.Activity
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).Find(&activities).Error
	return &activities, err
}

func (s *ActivityService) Update(userID, id uint, activity *models.Activity) error {
	if err := middleware.ValidateStruct(activity); err != nil {
		return err
	}
	return s.DB.Model(&models.Activity{}).Where("id = ? AND user_id = ?", id, userID).Updates(activity).Error
}

func (s *ActivityService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(models.Activity{}).Error
}
