package summaries

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

var defaultAnimalInputCategories = []summariesModels.CostCategory{
	{Name: "Feed", Type: "animal", Category: "input", IsDefault: true},
	{Name: "Vaccination", Type: "animal", Category: "input", IsDefault: true},
	{Name: "Medicine", Type: "animal", Category: "input", IsDefault: true},
	{Name: "Labor", Type: "animal", Category: "input", IsDefault: true},
	{Name: "Transport", Type: "animal", Category: "input", IsDefault: true},
	{Name: "Miscellaneous", Type: "animal", Category: "input", IsDefault: true},
}

var defaultAnimalActivityCategories = []summariesModels.CostCategory{
	{Name: "Milking", Type: "animal", Category: "activity", IsDefault: true},
	{Name: "Breeding", Type: "animal", Category: "activity", IsDefault: true},
	{Name: "Health Check", Type: "animal", Category: "activity", IsDefault: true},
	{Name: "Grazing", Type: "animal", Category: "activity", IsDefault: true},
	{Name: "Miscellaneous", Type: "animal", Category: "activity", IsDefault: true},
}

var defaultPlantInputCategories = []summariesModels.CostCategory{
	{Name: "Seeds", Type: "plant", Category: "input", IsDefault: true},
	{Name: "Nursery", Type: "plant", Category: "input", IsDefault: true},
	{Name: "Water", Type: "plant", Category: "input", IsDefault: true},
	{Name: "Labor", Type: "plant", Category: "input", IsDefault: true},
	{Name: "Transport", Type: "plant", Category: "input", IsDefault: true},
	{Name: "Fertilizer", Type: "plant", Category: "input", IsDefault: true},
	{Name: "Miscellaneous", Type: "plant", Category: "input", IsDefault: true},
}

var defaultPlantActivityCategories = []summariesModels.CostCategory{
	{Name: "Planting", Type: "plant", Category: "activity", IsDefault: true},
	{Name: "Harvesting", Type: "plant", Category: "activity", IsDefault: true},
	{Name: "Weeding", Type: "plant", Category: "activity", IsDefault: true},
	{Name: "Irrigation", Type: "plant", Category: "activity", IsDefault: true},
	{Name: "Pruning", Type: "plant", Category: "activity", IsDefault: true},
	{Name: "Miscellaneous", Type: "plant", Category: "activity", IsDefault: true},
}

func (s *CostCategoryService) Create(UserID uint, costCategory *summariesModels.CostCategory) error {
	costCategory.UserID = UserID
	if err := middleware.ValidateStruct(costCategory); err != nil {
		return err
	}
	return s.DB.Create(costCategory).Error
}

func (s *CostCategoryService) List(UserID uint) ([]summariesModels.CostCategory, error) {
	var count int64
	s.DB.Model(&summariesModels.CostCategory{}).Where("user_id = ?", UserID).Count(&count)

	if count == 0 {
		if err := s.InitializeDefaultCategories(UserID); err != nil {
			return nil, err
		}
	}

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
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&summariesModels.CostCategory{}).Error
}

func (s *CostCategoryService) GetDefaultCategories() ([]summariesModels.CostCategory, error) {
	var costCategories []summariesModels.CostCategory
	return costCategories, s.DB.Where("is_default = true").Find(&costCategories).Error
}

func (s *CostCategoryService) InitializeDefaultCategories(userID uint) error {
	allDefaults := append(
		append(
			append(defaultAnimalInputCategories, defaultAnimalActivityCategories...),
			defaultPlantInputCategories...),
		defaultPlantActivityCategories...)

	for _, category := range allDefaults {
		category.UserID = userID
		categoryCopy := category
		if err := s.Create(userID, &categoryCopy); err != nil {
			return err
		}
	}
	return nil
}
