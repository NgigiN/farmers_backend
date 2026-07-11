package plants

import (
	"errors"
	seasonModels "farm-backend/internal/models/plants"

	"gorm.io/gorm"
)

type SeasonService struct {
	DB *gorm.DB
}

func NewSeasonService(db *gorm.DB) *SeasonService {
	return &SeasonService{DB: db}
}

func (s *SeasonService) validateSeason(userID uint, season *seasonModels.Season) error {
	if !season.EndDate.IsZero() && season.EndDate.Before(season.StartDate) {
		return errors.New("end date must be on or after start date")
	}
	if err := s.DB.Where("id = ? AND user_id = ?", season.PlantID, userID).
		First(&seasonModels.Plant{}).Error; err != nil {
		return errors.New("plant not found or does not belong to user")
	}
	if err := s.DB.Where("id = ? AND user_id = ?", season.LandID, userID).
		First(&seasonModels.Land{}).Error; err != nil {
		return errors.New("land not found or does not belong to user")
	}
	return nil
}

func (s *SeasonService) Create(UserID uint, season *seasonModels.Season) error {
	season.UserID = UserID
	if err := s.validateSeason(UserID, season); err != nil {
		return err
	}
	return s.DB.Create(season).Error
}

func (s *SeasonService) List(UserID uint) ([]seasonModels.Season, error) {
	var seasons []seasonModels.Season
	return seasons, s.DB.Where("user_id = ?", UserID).Find(&seasons).Error
}

func (s *SeasonService) Get(UserID uint, id uint) (*seasonModels.Season, error) {
	var season seasonModels.Season
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).First(&season).Error
	if err != nil {
		return nil, err
	}
	return &season, nil
}

func (s *SeasonService) Update(userID, id uint, season *seasonModels.Season) error {
	if err := s.validateSeason(userID, season); err != nil {
		return err
	}
	return s.DB.Model(&seasonModels.Season{}).
		Where("id = ? AND user_id = ?", id, userID).
		Select("Name", "PlantID", "LandID", "StartDate", "EndDate").
		Updates(season).Error
}

func (s *SeasonService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&seasonModels.Season{}).Error
}