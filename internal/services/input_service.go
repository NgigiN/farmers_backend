// Package services provides functions necessary for business logic
package services

import (
	"sync"

	"farm-backend/internal/models"

	"gorm.io/gorm"
)

type CostService struct {
	DB *gorm.DB
}

func NewCostService(db *gorm.DB) *CostService {
	return &CostService{DB: db}
}

func (s *CostService) GetBreakdownByInputType(userID uint, seasonID uint) ([]InputCostBreakdown, error) {
	var inputs []models.Input
	if err := s.DB.Where("season_id = ? AND season_id IN (SELECT id FROM seasons where user_id = ?)", seasonID, userID).Find(&inputs).Error; err != nil {
		return nil, err
	}
	types := []string{"Seeds", "Nursery", "Water", "Labor", "Transport", "Miscellaneous", "Fertelizer", "Salaries"}

	breakdown := make([]InputCostBreakdown, len(types))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, t := range types {
		wg.Add(1)
		go func(idx int, inputType string) {
			defer wg.Done()
			var totalCost float64
			for _, inp := range inputs {
				if inp.Type == inputType {
					totalCost += inp.Cost
				}
			}
			mu.Lock()
			breakdown[idx] = InputCostBreakdown{Type: inputType, TotalCost: totalCost}
			mu.Unlock()
		}(i, t)
	}
	return breakdown, nil
}

type InputCostBreakdown struct {
	Type      string
	TotalCost float64
}
