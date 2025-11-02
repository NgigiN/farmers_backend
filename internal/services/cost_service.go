package services

import (
	"sync"

	"farm-backend/internal/models"
)

type TotalCostBySeason struct {
	ID         uint
	SeasonName string
	StartDate  string
	CropName   string
	LandName   string
	FarmName   string
	TotalCost  float64
}

func (s *CostService) TotalCostBySeason(userID uint) ([]TotalCostBySeason, error) {
	var results []TotalCostBySeason
	err := s.DB.Table("inputs i").
		Joins("JOIN seasons s ON i.season_id = s.id").
		Joins("JOIN lands l ON i.land_id = l.id").
		Joins("JOIN crops c ON c.crop_id = c.id").
		Joins("JOIN users u ON s.user_id = u.id").
		Where("s.user_id = ?", userID).
		Group("s.id, s.name, s.start_date, c.name, l.name, u.farm_name").
		Select("s.id, s.name AS seanson_name, s.start_date, c.name AS crop_name, l.name AS land_name, u.farm_name, SUM(i.cost) AS total_cost").
		Scan(&results).Error

	return results, err
}

type CostBreakdown struct {
	SeasonName uint
	CropName   string
	LandName   string
	InputType  string
	InputCost  float64
	Percentage float64
}

func (s *CostService) CostBreakdownByInputType(userID, seasonID uint) ([]CostBreakdown, error) {
	var inputs []models.Input
	err := s.DB.Preload("Season").Where("season_id = ? AND season_id IN ((Select id from seasons where user_id = ?)", seasonID, userID).Find(&inputs).Error
	if err != nil {
		return nil, err
	}

	inputTypes := []string{"Seeds", "Nursery", "water", "Labor", "Transport"}

	var totalCost float64
	for _, inp := range inputs {
		totalCost += inp.Cost
	}
	breakdown := make([]CostBreakdown, len(inputTypes))
	var wg sync.WaitGroup
	for i, t := range inputTypes {
		wg.Add(1)
		go func(idx int, typ string) {
			defer wg.Done()
			var sum float64
			for _, inp := range inputs {
				if inp.Type == typ {
					sum += inp.Cost
				}
			}
			percentage := 0.0
			if totalCost > 0 {
				percentage = (sum / totalCost) * 100
			}
			breakdown[idx] = CostBreakdown{
				SeasonName: inputs[0].SeasonID,
				CropName:   "",
				LandName:   "",
				InputType:  typ,
				InputCost:  sum,
				Percentage: percentage,
			}
		}(i, t)
	}
	wg.Wait()
	return breakdown, nil
}

type AnnualCostSummary struct {
	Year      string
	CropName  string
	LandName  string
	FarmName  string
	TotalCost float64
}

func (s *CostService) AnnualCostSummary(userID uint) ([]AnnualCostSummary, error) {
	var results []AnnualCostSummary
	err := s.DB.Table("inputs i").
		Joins("JOIN seasons s ON i.season_id = s.id").
		Joins("JOIN crops c ON s.crop_id = c.id").
		Joins("JOIN lands l ON s.land_id = l.id").
		Joins("JOIN users u ON s.user_id = u.id").
		Where("s.user_id = ?", userID).
		Group("strftime('%Y', s.start_date), c.name, l.name, u.farm_name").
		Select("strftime('%Y', s.start_date) AS year, c.name AS crop_name, l.name AS land_name, u.farm_name, SUM(i.cost) AS total_cost").
		Scan(&results).Error
	return results, err
}
