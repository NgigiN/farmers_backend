package plants

import (
	"farm-backend/internal/middleware"
	plantModels "farm-backend/internal/models/plants"

	"gorm.io/gorm"
)

type LandService struct {
	DB *gorm.DB
}

func NewLandService(db *gorm.DB) *LandService {
	return &LandService{DB: db}
}

func (s *LandService) Create(UserID uint, land *plantModels.Land) error {
	land.UserID = UserID
	if err := middleware.ValidateStruct(land); err != nil {
		return err
	}
	return s.DB.Create(land).Error
}

func (s *LandService) List(UserID uint) ([]plantModels.Land, error) {
	var land []plantModels.Land
	return land, s.DB.Where("user_id = ?", UserID).Find(&land).Error
}

func (s *LandService) Get(UserID uint, id uint) (*plantModels.Land, error) {
	var land plantModels.Land
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).First(&land).Error
	if err != nil {
		return nil, err
	}
	return &land, nil
}

func (s *LandService) Update(userID, id uint, land *plantModels.Land) error {
	if err := middleware.ValidateStruct(land); err != nil {
		return err
	}
	return s.DB.Model(&plantModels.Land{}).Where("id = ? AND user_id = ?", id, userID).Updates(land).Error
}

func (s *LandService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(plantModels.Land{}).Error
}
