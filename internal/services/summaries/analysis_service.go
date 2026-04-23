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
	Category   string  `json:"category"`
	Type       string  `json:"type"`
	Origin     string  `json:"origin"`
	TotalCost  float64 `json:"total_cost"`
	Percentage float64 `json:"percentage"`
}

func (s *AnalysisService) GetCostBreakdownByCategory(UserID uint) ([]CategoryBreakdown, error) {
	var plantInputs []struct {
		Type   string
		Origin string
		Cost   float64
	}
	err := s.DB.Table("inputs i").
		Joins("JOIN users u ON i.user_id = u.id").
		Where("i.user_id = ? AND i.source_type = ?", UserID, "plant").
		Select("i.type, u.farm_name as origin, SUM(i.cost) as cost").
		Group("i.type, u.farm_name").
		Scan(&plantInputs).Error
	if err != nil {
		return nil, err
	}

	var plantActivities []struct {
		Type   string
		Origin string
		Cost   float64
	}
	err = s.DB.Table("activities a").
		Joins("JOIN users u ON a.user_id = u.id").
		Where("a.user_id = ? AND a.source_type = ?", UserID, "plant").
		Select("a.type, u.farm_name as origin, SUM(a.cost) as cost").
		Group("a.type, u.farm_name").
		Scan(&plantActivities).Error
	if err != nil {
		return nil, err
	}

	var animalInputs []struct {
		Type   string
		Origin string
		Cost   float64
	}
	err = s.DB.Table("inputs i").
		Joins("LEFT JOIN herds h ON i.source_id = h.id").
		Where("i.user_id = ? AND i.source_type = ?", UserID, "animal").
		Select("i.type, COALESCE(h.name, 'General Animal') as origin, SUM(i.cost) as cost").
		Group("i.type, h.name").
		Scan(&animalInputs).Error
	if err != nil {
		return nil, err
	}

	var animalActivities []struct {
		Type   string
		Origin string
		Cost   float64
	}
	err = s.DB.Table("activities a").
		Joins("LEFT JOIN herds h ON a.source_id = h.id").
		Where("a.user_id = ? AND a.source_type = ?", UserID, "animal").
		Select("a.type, COALESCE(h.name, 'General Animal') as origin, SUM(a.cost) as cost").
		Group("a.type, h.name").
		Scan(&animalActivities).Error
	if err != nil {
		return nil, err
	}

	var infrastructures []struct {
		Type   string
		Origin string
		Cost   float64
	}
	err = s.DB.Table("infrastructures").
		Where("user_id = ?", UserID).
		Select("type, name as origin, SUM(cost) as cost").
		Group("type, name").
		Scan(&infrastructures).Error
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string]CategoryBreakdown)

	for _, inp := range plantInputs {
		key := fmt.Sprintf("%s_%s", inp.Type, inp.Origin)
		cat := categoryMap[key]
		cat.Category = inp.Type
		cat.Type = "plant"
		cat.Origin = inp.Origin
		cat.TotalCost += inp.Cost
		categoryMap[key] = cat
	}

	for _, act := range plantActivities {
		key := fmt.Sprintf("%s_%s", act.Type, act.Origin)
		cat := categoryMap[key]
		cat.Category = act.Type
		cat.Type = "plant"
		cat.Origin = act.Origin
		cat.TotalCost += act.Cost
		categoryMap[key] = cat
	}

	for _, inp := range animalInputs {
		key := fmt.Sprintf("%s_%s", inp.Type, inp.Origin)
		cat := categoryMap[key]
		cat.Category = inp.Type
		cat.Type = "animal"
		cat.Origin = inp.Origin
		cat.TotalCost += inp.Cost
		categoryMap[key] = cat
	}

	for _, act := range animalActivities {
		key := fmt.Sprintf("%s_%s", act.Type, act.Origin)
		cat := categoryMap[key]
		cat.Category = act.Type
		cat.Type = "animal"
		cat.Origin = act.Origin
		cat.TotalCost += act.Cost
		categoryMap[key] = cat
	}

	for _, inf := range infrastructures {
		key := fmt.Sprintf("%s_%s", inf.Type, inf.Origin)
		cat := categoryMap[key]
		cat.Category = inf.Type
		cat.Type = "animal"
		cat.Origin = inf.Origin
		cat.TotalCost += inf.Cost
		categoryMap[key] = cat
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
	Month        string           `json:"month"`
	TotalCosts   float64          `json:"total_costs"`
	TotalRevenue float64          `json:"total_revenue"`
	Profit       float64          `json:"profit"`
	Breakdown    MonthlyBreakdown `json:"breakdown"`
}

type MonthlyBreakdown struct {
	Costs   MonthlyCostBreakdown    `json:"costs"`
	Revenue MonthlyRevenueBreakdown `json:"revenue"`
}

type MonthlyCostBreakdown struct {
	Plant          float64 `json:"plant"`
	Animal         float64 `json:"animal"`
	Infrastructure float64 `json:"infrastructure"`
}

type MonthlyRevenueBreakdown struct {
	Plant  float64 `json:"plant"`
	Animal float64 `json:"animal"`
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

	plantInputMap := make(map[string]float64)
	for _, c := range plantInputCosts {
		plantInputMap[c.Month] = c.Total
	}
	plantActivityMap := make(map[string]float64)
	for _, c := range plantActivityCosts {
		plantActivityMap[c.Month] = c.Total
	}
	animalInputMap := make(map[string]float64)
	for _, c := range animalInputCosts {
		animalInputMap[c.Month] = c.Total
	}
	animalActivityMap := make(map[string]float64)
	for _, c := range animalActivityCosts {
		animalActivityMap[c.Month] = c.Total
	}
	infrastructureMap := make(map[string]float64)
	for _, c := range infrastructureCosts {
		infrastructureMap[c.Month] = c.Total
	}

	// Since we need to split revenue by source, we need another query or filter the existing one.
	// But let's look at the existing revenue query: it doesn't distinguish source.
	// Let's add specific queries for revenue sources too.
	var plantRevenueResults []struct {
		Month string
		Total float64
	}
	err = s.DB.Raw("SELECT strftime('%Y-%m', date) as month, SUM(total) as total FROM revenues WHERE user_id = ? AND source = 'plant' AND strftime('%Y', date) = ? GROUP BY strftime('%Y-%m', date)", UserID, yearStr).Scan(&plantRevenueResults).Error
	if err != nil {
		return nil, err
	}
	var animalRevenueResults []struct {
		Month string
		Total float64
	}
	err = s.DB.Raw("SELECT strftime('%Y-%m', date) as month, SUM(total) as total FROM revenues WHERE user_id = ? AND source = 'animal' AND strftime('%Y', date) = ? GROUP BY strftime('%Y-%m', date)", UserID, yearStr).Scan(&animalRevenueResults).Error
	if err != nil {
		return nil, err
	}

	plantRevMap := make(map[string]float64)
	for _, r := range plantRevenueResults {
		plantRevMap[r.Month] = r.Total
	}
	animalRevMap := make(map[string]float64)
	for _, r := range animalRevenueResults {
		animalRevMap[r.Month] = r.Total
	}

	monthSet := make(map[string]bool)
	maps := []map[string]float64{plantInputMap, plantActivityMap, animalInputMap, animalActivityMap, infrastructureMap, plantRevMap, animalRevMap}
	for _, m := range maps {
		for month := range m {
			monthSet[month] = true
		}
	}

	var summaries []MonthlySummary
	for month := range monthSet {
		pInp := plantInputMap[month]
		pAct := plantActivityMap[month]
		aInp := animalInputMap[month]
		aAct := animalActivityMap[month]
		inf := infrastructureMap[month]

		pRev := plantRevMap[month]
		aRev := animalRevMap[month]

		totalCosts := pInp + pAct + aInp + aAct + inf
		totalRevenue := pRev + aRev

		summaries = append(summaries, MonthlySummary{
			Month:        month,
			TotalCosts:   totalCosts,
			TotalRevenue: totalRevenue,
			Profit:       totalRevenue - totalCosts,
			Breakdown: MonthlyBreakdown{
				Costs: MonthlyCostBreakdown{
					Plant:          pInp + pAct,
					Animal:         aInp + aAct,
					Infrastructure: inf,
				},
				Revenue: MonthlyRevenueBreakdown{
					Plant:  pRev,
					Animal: aRev,
				},
			},
		})
	}

	return summaries, nil
}
