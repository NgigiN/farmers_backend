package services

import (
	"farm-backend/internal/middleware"
	inputModels "farm-backend/internal/models/plants"

	"gorm.io/gorm"
)

type InputService struct {
	DB *gorm.DB
}

func NewInputService(db *gorm.DB) *InputService {
	return &InputService{DB: db}
}

func (s *InputService) Create(UserID uint, input *inputModels.Input) error {
	if err := s.DB.Where("id = ? AND user_id = ?", input.SeasonID, UserID).First(&inputModels.Season{}).Error; err != nil {
		return err
	}
	if err := middleware.ValidateStruct(input); err != nil {
		return err
	}
	return s.DB.Create(input).Error
}

func (s *InputService) List(UserID uint) ([]inputModels.Input, error) {
	var inputs []inputModels.Input
	err := s.DB.Table("inputs").
		Joins("JOIN seasons ON inputs.season_id = seasons.id").
		Where("seasons.user_id = ?", UserID).
		Find(&inputs).Error
	return inputs, err
}

func (s *InputService) Get(UserID uint, id uint) (*inputModels.Input, error) {
	var input inputModels.Input
	err := s.DB.Table("inputs").
		Joins("JOIN seasons ON inputs.season_id = seasons.id").
		Where("inputs.id = ? AND seasons.user_id = ?", id, UserID).
		First(&input).Error
	if err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *InputService) Update(userID, id uint, input *inputModels.Input) error {
	if err := s.DB.Where("id = ? AND user_id = ?", input.SeasonID, userID).First(&inputModels.Season{}).Error; err != nil {
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
		Delete(&inputModels.Input{}).Error
}
