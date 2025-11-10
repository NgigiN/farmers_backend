package services

import (
	summariesModels "farm-backend/internal/models/summaries"

	"farm-backend/internal/middleware"

	"gorm.io/gorm"
)

type CostCategoryService struct {
	DB *gorm.DB
}

func NewCostCategoryService(db *gorm.DB) *CostCategoryService {
	return &CostCategoryService{DB: db}
}

var defaultCategories = map[string]summariesModels.CostCategory{
	"Feed": {
		Name:      "Feed",
		Type:      "Animal",
		IsDefault: true,
	},
	"Medication": {
		Name:      "Medication",
		Type:      "Animal",
		IsDefault: true,
	},
	"Vaccination": {
		Name:      "Vaccination",
		Type:      "Animal",
		IsDefault: true,
	},
	"Labor": {
		Name:      "Labor",
		Type:      "Animal",
		IsDefault: true,
	},
}

func (s *CostCategoryService) Create(UserID uint, costCategory *summariesModels.CostCategory) error {
	costCategory.UserID = UserID
	if err := middleware.ValidateStruct(costCategory); err != nil {
		return err
	}
	return s.DB.Create(costCategory).Error
}

func (s *CostCategoryService) List(UserID uint) ([]summariesModels.CostCategory, error) {
	var costCategories []summariesModels.CostCategory
	return costCategories, s.DB.Where("user_id = ?", UserID).Find(&costCategories).Error
}

func (s *CostCategoryService) Get(UserID uint, id uint) (*summariesModels.CostCategory, error) {
	var costCategory summariesModels.CostCategory
	err := s.DB.Where("id = ? AND user_id = ?", id, UserID).First(&costCategory).Error
	if err != nil {
		return nil, err
	}
	return &costCategory, nil
}

func (s *CostCategoryService) Update(userID, id uint, costCategory *summariesModels.CostCategory) error {
	if err := middleware.ValidateStruct(costCategory); err != nil {
		return err
	}
	return s.DB.Model(&summariesModels.CostCategory{}).Where("id = ? AND user_id = ?", id, userID).Updates(costCategory).Error
}

func (s *CostCategoryService) Delete(userID, id uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(summariesModels.CostCategory{}).Error
}

func (s *CostCategoryService) GetDefaultCategories() ([]summariesModels.CostCategory, error) {
	var costCategories []summariesModels.CostCategory
	return costCategories, s.DB.Where("is_default = true").Find(&costCategories).Error
}

func (s *CostCategoryService) InitializeDefaultCategories(userID uint) error {
	defaultCategories, err := s.GetDefaultCategories()
	if err != nil {
		return err
	}
	for _, category := range defaultCategories {
		category.UserID = userID
		err := s.Create(userID, &category)
		if err != nil {
			return err
		}
	}
	return nil
}
