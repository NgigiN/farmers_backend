package services

import (
	"farm-backend/internal/middleware"
	"farm-backend/internal/models"

	"gorm.io/gorm"
)

type InputService struct {
	DB *gorm.DB
}

func NewInputService(db *gorm.DB) *InputService {
	return &InputService{DB: db}
}

func (s *InputService) Create(UserID uint, input *models.Input) error {
	if err := s.DB.Where("id = ? AND user_id = ?", input.SeasonID, UserID).First(&models.Season{}).Error; err != nil {
		return err
	}
	if err := middleware.ValidateStruct(input); err != nil {
		return err
	}
	return s.DB.Create(input).Error
}

func (s *InputService) List(UserID uint) ([]models.Input, error) {
	var inputs []models.Input
	err := s.DB.Table("inputs").
		Joins("JOIN seasons ON inputs.season_id = seasons.id").
		Where("seasons.user_id = ?", UserID).
		Find(&inputs).Error
	return inputs, err
}

func (s *InputService) Get(UserID uint, id uint) (*models.Input, error) {
	var input models.Input
	err := s.DB.Table("inputs").
		Joins("JOIN seasons ON inputs.season_id = seasons.id").
		Where("inputs.id = ? AND seasons.user_id = ?", id, UserID).
		First(&input).Error
	if err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *InputService) Update(userID, id uint, input *models.Input) error {
	if err := s.DB.Where("id = ? AND user_id = ?", input.SeasonID, userID).First(&models.Season{}).Error; err != nil {
		return err
	}
	if err := middleware.ValidateStruct(input); err != nil {
		return err
	}
	return s.DB.Table("inputs").
		Joins("JOIN seasons ON inputs.season_id = seasons.id").
		Where("inputs.id = ? AND seasons.user_id = ?", id, userID).
		Updates(input).Error
}

func (s *InputService) Delete(userID, id uint) error {
	return s.DB.Table("inputs").
		Joins("JOIN seasons ON inputs.season_id = seasons.id").
		Where("inputs.id = ? AND seasons.user_id = ?", id, userID).
		Delete(&models.Input{}).Error
}
