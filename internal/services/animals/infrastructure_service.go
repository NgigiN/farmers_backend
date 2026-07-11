// Infrastructure service will represent one time costs like fences,
// stores, buildings
package animals

import (
	animalModels "farm-backend/internal/models/animals"

	"gorm.io/gorm"
)

type InfrastructureService struct {
	DB *gorm.DB
}

func NewInfrastructureService(db *gorm.DB) *InfrastructureService {
	return &InfrastructureService{DB: db}
}

func (s *InfrastructureService) Create(UserID uint, infrastructure *animalModels.Infrastructure) error {
	infrastructure.UserID = UserID
	return s.DB.Create(infrastructure).Error
}

func (s *InfrastructureService) List(UserID uint) ([]animalModels.Infrastructure, error) {
	var infrastructures []animalModels.Infrastructure
	return infrastructures, s.DB.Where("user_id = ?", UserID).Find(&infrastructures).Error
}

func (s *InfrastructureService) Get(UserID uint, id uint) (*animalModels.Infrastructure, error) {
	var infrastructure animalModels.Infrastructure
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).First(&infrastructure).Error
	if err != nil {
		return nil, err
	}
	return &infrastructure, nil
}

func (s *InfrastructureService) Update(userID, id uint, infrastructure *animalModels.Infrastructure) error {
	return s.DB.Model(&animalModels.Infrastructure{}).
		Where("id = ? AND user_id = ?", id, userID).
		Select("Type", "Name", "Location", "Cost", "Date", "Notes").
		Updates(infrastructure).Error
}

func (s *InfrastructureService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&animalModels.Infrastructure{}).Error
}