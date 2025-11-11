package summaries

import (
	"fmt"

	"gorm.io/gorm"
)

type AnalysisService struct {
	DB *gorm.DB
}

func NewAnalysisService(db *gorm.DB) *AnalysisService {
	return &AnalysisService{DB: db}
}

func (s *AnalysisService) GetTotalPlantCosts(UserID uint) (float64, error) {
	var inputTotal float64
	err := s.DB.Table("inputs").
		Where("user_id = ? AND source_type = ?", UserID, "plant").
		Select("COALESCE(SUM(cost), 0) as total").
		Scan(&inputTotal).Error
	if err != nil {
		return 0, err
	}

	var activityTotal float64
	err = s.DB.Table("activities").
		Where("user_id = ? AND source_type = ?", UserID, "plant").
		Select("COALESCE(SUM(cost), 0) as total").
		Scan(&activityTotal).Error
	if err != nil {
		return 0, err
	}

	return inputTotal + activityTotal, nil
}

func (s *AnalysisService) GetTotalAnimalCosts(UserID uint) (float64, error) {
	var inputTotal float64
	err := s.DB.Table("inputs").
		Where("user_id = ? AND source_type = ?", UserID, "animal").
		Select("COALESCE(SUM(cost), 0) as total").
		Scan(&inputTotal).Error
	if err != nil {
		return 0, err
	}

	var activityTotal float64
	err = s.DB.Table("activities").
		Where("user_id = ? AND source_type = ?", UserID, "animal").
		Select("COALESCE(SUM(cost), 0) as total").
		Scan(&activityTotal).Error
	if err != nil {
		return 0, err
	}

	var infrastructureTotal float64
	err = s.DB.Table("infrastructures").
		Where("user_id = ?", UserID).
		Select("COALESCE(SUM(cost), 0) as total").
		Scan(&infrastructureTotal).Error
	if err != nil {
		return 0, err
	}

	return inputTotal + activityTotal + infrastructureTotal, nil
}

func (s *AnalysisService) GetTotalCosts(UserID uint) (float64, error) {
	plantCosts, err := s.GetTotalPlantCosts(UserID)
	if err != nil {
		return 0, err
	}

	animalCosts, err := s.GetTotalAnimalCosts(UserID)
	if err != nil {
		return 0, err
	}

	return plantCosts + animalCosts, nil
}

func (s *AnalysisService) GetTotalRevenue(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("revenues").
		Where("user_id = ?", UserID).
		Select("COALESCE(SUM(total), 0) as total").
		Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetPlantRevenue(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("revenues").
		Where("user_id = ? AND source = ?", UserID, "plant").
		Select("COALESCE(SUM(total), 0) as total").
		Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetAnimalRevenue(UserID uint) (float64, error) {
	var total float64
	err := s.DB.Table("revenues").
		Where("user_id = ? AND source = ?", UserID, "animal").
		Select("COALESCE(SUM(total), 0) as total").
		Scan(&total).Error
	return total, err
}

func (s *AnalysisService) GetProfit(UserID uint) (float64, error) {
	revenue, err := s.GetTotalRevenue(UserID)
	if err != nil {
		return 0, err
	}

	costs, err := s.GetTotalCosts(UserID)
	if err != nil {
		return 0, err
	}

	return revenue - costs, nil
}

func (s *AnalysisService) GetPlantProfit(UserID uint) (float64, error) {
	revenue, err := s.GetPlantRevenue(UserID)
	if err != nil {
		return 0, err
	}

	costs, err := s.GetTotalPlantCosts(UserID)
	if err != nil {
		return 0, err
	}

	return revenue - costs, nil
}

func (s *AnalysisService) GetAnimalProfit(UserID uint) (float64, error) {
	revenue, err := s.GetAnimalRevenue(UserID)
	if err != nil {
		return 0, err
	}

	costs, err := s.GetTotalAnimalCosts(UserID)
	if err != nil {
		return 0, err
	}

	return revenue - costs, nil
}

func (s *AnalysisService) GetProfitBySource(UserID uint) (plantProfit, animalProfit float64, err error) {
	plantProfit, err = s.GetPlantProfit(UserID)
	if err != nil {
		return 0, 0, err
	}

	animalProfit, err = s.GetAnimalProfit(UserID)
	if err != nil {
		return 0, 0, err
	}

	return plantProfit, animalProfit, nil
}

func (s *AnalysisService) GetCostBreakdownBySource(UserID uint) (plantCosts, animalCosts float64, err error) {
	plantCosts, err = s.GetTotalPlantCosts(UserID)
	if err != nil {
		return 0, 0, err
	}

	animalCosts, err = s.GetTotalAnimalCosts(UserID)
	if err != nil {
		return 0, 0, err
	}

	return plantCosts, animalCosts, nil
}

type CategoryBreakdown struct {
	Category   string
	Type       string
	TotalCost  float64
	Percentage float64
}

func (s *AnalysisService) GetCostBreakdownByCategory(UserID uint) ([]CategoryBreakdown, error) {
	var plantInputs []struct {
		Type string
		Cost float64
	}
	err := s.DB.Table("inputs").
		Where("user_id = ? AND source_type = ?", UserID, "plant").
		Select("type, SUM(cost) as cost").
		Group("type").
		Scan(&plantInputs).Error
	if err != nil {
		return nil, err
	}

	var plantActivities []struct {
		Type string
		Cost float64
	}
	err = s.DB.Table("activities").
		Where("user_id = ? AND source_type = ?", UserID, "plant").
		Select("type, SUM(cost) as cost").
		Group("type").
		Scan(&plantActivities).Error
	if err != nil {
		return nil, err
	}

	var animalInputs []struct {
		Type string
		Cost float64
	}
	err = s.DB.Table("inputs").
		Where("user_id = ? AND source_type = ?", UserID, "animal").
		Select("type, SUM(cost) as cost").
		Group("type").
		Scan(&animalInputs).Error
	if err != nil {
		return nil, err
	}

	var animalActivities []struct {
		Type string
		Cost float64
	}
	err = s.DB.Table("activities").
		Where("user_id = ? AND source_type = ?", UserID, "animal").
		Select("type, SUM(cost) as cost").
		Group("type").
		Scan(&animalActivities).Error
	if err != nil {
		return nil, err
	}

	var infrastructures []struct {
		Type string
		Cost float64
	}
	err = s.DB.Table("infrastructures").
		Where("user_id = ?", UserID).
		Select("type, SUM(cost) as cost").
		Group("type").
		Scan(&infrastructures).Error
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string]CategoryBreakdown)

	for _, inp := range plantInputs {
		cat := categoryMap[inp.Type]
		cat.Category = inp.Type
		cat.Type = "plant"
		cat.TotalCost += inp.Cost
		categoryMap[inp.Type] = cat
	}

	for _, act := range plantActivities {
		cat := categoryMap[act.Type]
		cat.Category = act.Type
		cat.Type = "plant"
		cat.TotalCost += act.Cost
		categoryMap[act.Type] = cat
	}

	for _, inp := range animalInputs {
		cat := categoryMap[inp.Type]
		cat.Category = inp.Type
		cat.Type = "animal"
		cat.TotalCost += inp.Cost
		categoryMap[inp.Type] = cat
	}

	for _, act := range animalActivities {
		cat := categoryMap[act.Type]
		cat.Category = act.Type
		cat.Type = "animal"
		cat.TotalCost += act.Cost
		categoryMap[act.Type] = cat
	}

	for _, inf := range infrastructures {
		cat := categoryMap[inf.Type]
		cat.Category = inf.Type
		cat.Type = "animal"
		cat.TotalCost += inf.Cost
		categoryMap[inf.Type] = cat
	}

	var totalCost float64
	for _, cat := range categoryMap {
		totalCost += cat.TotalCost
	}

	var breakdown []CategoryBreakdown
	for _, cat := range categoryMap {
		if totalCost > 0 {
			cat.Percentage = (cat.TotalCost / totalCost) * 100
		}
		breakdown = append(breakdown, cat)
	}

	return breakdown, nil
}

type RevenueTypeBreakdown struct {
	Type         string
	Source       string
	TotalRevenue float64
	Percentage   float64
}

func (s *AnalysisService) GetRevenueBreakdownByType(UserID uint) ([]RevenueTypeBreakdown, error) {
	var revenues []struct {
		Type  string
		Total float64
	}
	err := s.DB.Table("revenues").
		Where("user_id = ?", UserID).
		Select("type, SUM(total) as total").
		Group("type").
		Scan(&revenues).Error
	if err != nil {
		return nil, err
	}

	var totalRevenue float64
	for _, rev := range revenues {
		totalRevenue += rev.Total
	}

	var breakdown []RevenueTypeBreakdown
	for _, rev := range revenues {
		percentage := 0.0
		if totalRevenue > 0 {
			percentage = (rev.Total / totalRevenue) * 100
		}

		var source string
		err = s.DB.Table("revenues").
			Where("user_id = ? AND type = ?", UserID, rev.Type).
			Select("source").
			Limit(1).
			Scan(&source).Error
		if err != nil {
			source = "unknown"
		}

		breakdown = append(breakdown, RevenueTypeBreakdown{
			Type:         rev.Type,
			Source:       source,
			TotalRevenue: rev.Total,
			Percentage:   percentage,
		})
	}

	return breakdown, nil
}

type MonthlySummary struct {
	Month        string
	TotalCosts   float64
	TotalRevenue float64
	Profit       float64
}

func (s *AnalysisService) GetMonthlySummary(UserID uint, year int) ([]MonthlySummary, error) {
	yearStr := fmt.Sprintf("%d", year)
	if year <= 0 {
		return nil, fmt.Errorf("invalid year: %d", year)
	}

	var revenueResults []struct {
		Month string
		Total float64
	}
	err := s.DB.Raw("SELECT strftime('%Y-%m', date) as month, SUM(total) as total FROM revenues WHERE user_id = ? AND strftime('%Y', date) = ? GROUP BY strftime('%Y-%m', date)", UserID, yearStr).Scan(&revenueResults).Error
	if err != nil {
		return nil, err
	}

	var plantInputCosts []struct {
		Month string
		Total float64
	}
	err = s.DB.Raw("SELECT strftime('%Y-%m', date) as month, SUM(cost) as total FROM inputs WHERE user_id = ? AND source_type = ? AND strftime('%Y', date) = ? GROUP BY strftime('%Y-%m', date)", UserID, "plant", yearStr).Scan(&plantInputCosts).Error
	if err != nil {
		return nil, err
	}

	var plantActivityCosts []struct {
		Month string
		Total float64
	}
	err = s.DB.Raw("SELECT strftime('%Y-%m', date) as month, SUM(cost) as total FROM activities WHERE user_id = ? AND source_type = ? AND strftime('%Y', date) = ? GROUP BY strftime('%Y-%m', date)", UserID, "plant", yearStr).Scan(&plantActivityCosts).Error
	if err != nil {
		return nil, err
	}

	var animalInputCosts []struct {
		Month string
		Total float64
	}
	err = s.DB.Raw("SELECT strftime('%Y-%m', date) as month, SUM(cost) as total FROM inputs WHERE user_id = ? AND source_type = ? AND strftime('%Y', date) = ? GROUP BY strftime('%Y-%m', date)", UserID, "animal", yearStr).Scan(&animalInputCosts).Error
	if err != nil {
		return nil, err
	}

	var animalActivityCosts []struct {
		Month string
		Total float64
	}
	err = s.DB.Raw("SELECT strftime('%Y-%m', date) as month, SUM(cost) as total FROM activities WHERE user_id = ? AND source_type = ? AND strftime('%Y', date) = ? GROUP BY strftime('%Y-%m', date)", UserID, "animal", yearStr).Scan(&animalActivityCosts).Error
	if err != nil {
		return nil, err
	}

	var infrastructureCosts []struct {
		Month string
		Total float64
	}
	err = s.DB.Raw("SELECT strftime('%Y-%m', date) as month, SUM(cost) as total FROM infrastructures WHERE user_id = ? AND strftime('%Y', date) = ? GROUP BY strftime('%Y-%m', date)", UserID, yearStr).Scan(&infrastructureCosts).Error
	if err != nil {
		return nil, err
	}

	costMap := make(map[string]float64)
	for _, c := range plantInputCosts {
		costMap[c.Month] += c.Total
	}
	for _, c := range plantActivityCosts {
		costMap[c.Month] += c.Total
	}
	for _, c := range animalInputCosts {
		costMap[c.Month] += c.Total
	}
	for _, c := range animalActivityCosts {
		costMap[c.Month] += c.Total
	}
	for _, c := range infrastructureCosts {
		costMap[c.Month] += c.Total
	}

	revenueMap := make(map[string]float64)
	for _, r := range revenueResults {
		revenueMap[r.Month] = r.Total
	}

	monthSet := make(map[string]bool)
	for month := range costMap {
		monthSet[month] = true
	}
	for month := range revenueMap {
		monthSet[month] = true
	}

	var summaries []MonthlySummary
	for month := range monthSet {
		costs := costMap[month]
		revenue := revenueMap[month]
		summaries = append(summaries, MonthlySummary{
			Month:        month,
			TotalCosts:   costs,
			TotalRevenue: revenue,
			Profit:       revenue - costs,
		})
	}

	return summaries, nil
}
