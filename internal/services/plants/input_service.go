package plants

import (
	"errors"
	animalModels "farm-backend/internal/models/animals"
	inputModels "farm-backend/internal/models/plants"
	summaryModels "farm-backend/internal/models/summaries"

	"gorm.io/gorm"
)

type InputService struct {
	DB *gorm.DB
}

func NewInputService(db *gorm.DB) *InputService {
	return &InputService{DB: db}
}

// validateInput centralises the checks shared by Create and Update.
func (s *InputService) validateInput(userID uint, input *inputModels.Input) error {
	if input.SourceType != "plant" && input.SourceType != "animal" {
		return errors.New("source_type must be either 'plant' or 'animal'")
	}

	if input.SourceType == "plant" {
		if err := s.DB.Where("id = ? AND user_id = ?", input.SourceID, userID).First(&inputModels.Season{}).Error; err != nil {
			return errors.New("season not found or does not belong to user")
		}
	} else {
		if err := s.DB.Where("id = ? AND user_id = ?", input.SourceID, userID).First(&animalModels.Herd{}).Error; err != nil {
			return errors.New("herd not found or does not belong to user")
		}
		if input.AnimalID != 0 {
			if err := s.DB.Where("id = ? AND user_id = ?", input.AnimalID, userID).First(&animalModels.Animal{}).Error; err != nil {
				return errors.New("animal not found or does not belong to user")
			}
		}
	}

	var category summaryModels.CostCategory
	if err := s.DB.Where("user_id = ? AND name = ? AND type = ? AND category = ?",
		userID, input.Type, input.SourceType, "input").First(&category).Error; err != nil {
		return errors.New("input type does not exist in cost categories — please create it first")
	}

	return nil
}

func (s *InputService) Create(userID uint, input *inputModels.Input) error {
	input.UserID = userID
	if err := s.validateInput(userID, input); err != nil {
		return err
	}
	return s.DB.Create(input).Error
}

// List returns all inputs for a user, optionally filtered by source_type at the DB level.
func (s *InputService) List(userID uint, sourceType string) ([]inputModels.Input, error) {
	var inputs []inputModels.Input
	q := s.DB.Where("user_id = ?", userID)
	if sourceType != "" {
		q = q.Where("source_type = ?", sourceType)
	}
	return inputs, q.Find(&inputs).Error
}

func (s *InputService) Get(userID uint, id uint) (*inputModels.Input, error) {
	var input inputModels.Input
	err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&input).Error
	if err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *InputService) Update(userID, id uint, input *inputModels.Input) error {
	if err := s.validateInput(userID, input); err != nil {
		return err
	}
	return s.DB.Model(&inputModels.Input{}).
		Where("id = ? AND user_id = ?", id, userID).
		Select("SourceType", "SourceID", "AnimalID", "Type", "Quantity", "Cost", "Date", "Notes").
		Updates(input).Error
}

func (s *InputService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&inputModels.Input{}).Error
}
