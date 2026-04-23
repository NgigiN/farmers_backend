package summaries

import (
	"sync"
	"time"

	plantModels "farm-backend/internal/models/plants"

	"gorm.io/gorm"
)

type CostService struct {
	DB *gorm.DB
}

func NewCostService(db *gorm.DB) *CostService {
	return &CostService{DB: db}
}

type TotalCostBySeason struct {
	ID         uint      `json:"id"`
	SeasonName string    `json:"season_name"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	PlantName  string    `json:"plant_name"`
	LandName   string    `json:"land_name"`
	FarmName   string    `json:"farm_name"`
	TotalCost  float64   `json:"total_cost"`
}

type DetailedCostItem struct {
	Type         string     `json:"type"` // "plant" or "animal"
	ID           uint       `json:"id"`
	Name         string     `json:"name"`           // Season name or Herd name
	Category     string     `json:"category"`       // Plant name or Animal type name
	Location     string     `json:"location"`       // Land name or Herd location
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	InputCost    float64    `json:"input_cost"`
	ActivityCost float64    `json:"activity_cost"`
	TotalCost    float64    `json:"total_cost"`
}

type TotalCostResponse struct {
	TotalOverallCost float64            `json:"total_overall_cost"`
	Details          []DetailedCostItem `json:"details"`
}

func (s *CostService) TotalCostBySeason(userID uint) ([]TotalCostBySeason, error) {
	var results []TotalCostBySeason
	err := s.DB.Table("inputs i").
		Joins("JOIN seasons s ON i.season_id = s.id").
		Joins("JOIN lands l ON s.land_id = l.id").
		Joins("JOIN plants c ON s.plant_id = c.id").
		Joins("JOIN users u ON s.user_id = u.id").
		Where("s.user_id = ?", userID).
		Group("s.id, s.name, s.start_date, s.end_date, c.name, l.name, u.farm_name").
		Select("s.id, s.name AS season_name, s.start_date, s.end_date, c.name AS plant_name, l.name AS land_name, u.farm_name, SUM(i.cost) AS total_cost").
		Scan(&results).Error

	return results, err
}

func (s *CostService) GetUnifiedDetailedCosts(userID uint) (*TotalCostResponse, error) {
	var details []DetailedCostItem

	// 1. Fetch Plant Costs (grouped by Season)
	var plantDetails []DetailedCostItem
	err := s.DB.Table("seasons s").
		Joins("JOIN plants p ON s.plant_id = p.id").
		Joins("JOIN lands l ON s.land_id = l.id").
		Where("s.user_id = ?", userID).
		Select(`
			'plant' as type,
			s.id,
			s.name,
			p.name as category,
			l.name as location,
			s.start_date,
			s.end_date,
			(SELECT COALESCE(SUM(cost), 0) FROM inputs WHERE source_type = 'plant' AND source_id = s.id) as input_cost,
			(SELECT COALESCE(SUM(cost), 0) FROM activities WHERE source_type = 'plant' AND source_id = s.id) as activity_cost
		`).
		Scan(&plantDetails).Error
	if err != nil {
		return nil, err
	}

	// 2. Fetch Animal Costs (grouped by Herd)
	var animalDetails []DetailedCostItem
	err = s.DB.Table("herds h").
		Joins("JOIN animal_types at ON h.animal_type_id = at.id").
		Where("h.user_id = ?", userID).
		Select(`
			'animal' as type,
			h.id,
			h.name,
			at.name as category,
			h.location,
			h.created_at as start_date,
			NULL as end_date,
			(SELECT COALESCE(SUM(cost), 0) FROM inputs WHERE source_type = 'animal' AND source_id = h.id) as input_cost,
			(SELECT COALESCE(SUM(cost), 0) FROM activities WHERE source_type = 'animal' AND source_id = h.id) as activity_cost
		`).
		Scan(&animalDetails).Error
	if err != nil {
		return nil, err
	}

	// Combine and calculate totals
	var overallTotal float64
	for i := range plantDetails {
		plantDetails[i].TotalCost = plantDetails[i].InputCost + plantDetails[i].ActivityCost
		overallTotal += plantDetails[i].TotalCost
		details = append(details, plantDetails[i])
	}

	for i := range animalDetails {
		animalDetails[i].TotalCost = animalDetails[i].InputCost + animalDetails[i].ActivityCost
		overallTotal += animalDetails[i].TotalCost
		details = append(details, animalDetails[i])
	}

	return &TotalCostResponse{
		TotalOverallCost: overallTotal,
		Details:          details,
	}, nil
}

type CostBreakdown struct {
	SeasonID   uint      `json:"season_id"`
	SeasonName string    `json:"season_name"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	PlantName  string    `json:"plant_name"`
	LandName   string    `json:"land_name"`
	InputType  string    `json:"input_type"`
	InputCost  float64   `json:"input_cost"`
	Percentage float64   `json:"percentage"`
}

func (s *CostService) CostBreakdownByInputType(userID, seasonID uint) ([]CostBreakdown, error) {
	var season plantModels.Season
	err := s.DB.Where("id = ? AND user_id = ?", seasonID, userID).
		Preload("Plant").
		Preload("Land").
		First(&season).Error
	if err != nil {
		return nil, err
	}

	var inputs []plantModels.Input
	err = s.DB.Where("season_id = ?", seasonID).Find(&inputs).Error
	if err != nil {
		return nil, err
	}

	if len(inputs) == 0 {
		return []CostBreakdown{}, nil
	}

	inputTypes := []string{"Seeds", "Nursery", "Water", "Labor", "Transport"}

	var totalCost float64
	for _, inp := range inputs {
		totalCost += inp.Cost
	}

	plantName := season.Plant.Name
	landName := season.Land.Name
	seasonName := season.Name
	startDate := season.StartDate
	endDate := season.EndDate

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
				SeasonID:   seasonID,
				SeasonName: seasonName,
				StartDate:  startDate,
				EndDate:    endDate,
				PlantName:  plantName,
				LandName:   landName,
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
	PlantName string
	LandName  string
	FarmName  string
	TotalCost float64
}

func (s *CostService) AnnualCostSummary(userID uint) ([]AnnualCostSummary, error) {
	var results []AnnualCostSummary
	err := s.DB.Table("inputs i").
		Joins("JOIN seasons s ON i.season_id = s.id").
		Joins("JOIN plants c ON s.plant_id = c.id").
		Joins("JOIN lands l ON s.land_id = l.id").
		Joins("JOIN users u ON s.user_id = u.id").
		Where("s.user_id = ?", userID).
		Group("strftime('%Y', s.start_date), c.name, l.name, u.farm_name").
		Select("strftime('%Y', s.start_date) AS year, c.name AS plant_name, l.name AS land_name, u.farm_name, SUM(i.cost) AS total_cost").
		Scan(&results).Error
	return results, err
}
