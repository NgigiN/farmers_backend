package validation

import (
	animalModels "farm-backend/internal/models/animals"
	plantModels "farm-backend/internal/models/plants"
	summaryModels "farm-backend/internal/models/summaries"
	userModels "farm-backend/internal/models/users"
)

func LandFromRequest(req *LandRequest) *plantModels.Land {
	return &plantModels.Land{
		Name:     req.Name,
		Size:     DerefFloat32(req.Size),
		Location: req.Location,
		SoilType: req.SoilType,
	}
}

func PlantFromRequest(req *PlantRequest) *plantModels.Plant {
	return &plantModels.Plant{
		Name:    req.Name,
		Variety: req.Variety,
	}
}

func SeasonFromRequest(req *SeasonRequest) *plantModels.Season {
	return &plantModels.Season{
		Name:      req.Name,
		PlantID:   req.PlantID,
		LandID:    req.LandID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
}

func HarvestFromRequest(req *HarvestRequest) *plantModels.Harvest {
	return &plantModels.Harvest{
		SeasonID: req.SeasonID,
		Quantity: req.Quantity,
		Unit:     req.Unit,
		Date:     req.Date,
		Notes:    req.Notes,
	}
}

func InputFromRequest(req *InputRequest) *plantModels.Input {
	return &plantModels.Input{
		SourceType: req.SourceType,
		SourceID:   req.SourceID,
		AnimalID:   req.AnimalID,
		Type:       req.Type,
		Quantity:   req.Quantity,
		Cost:       req.Cost,
		Date:       req.Date,
		Notes:      req.Notes,
	}
}

func ActivityFromRequest(req *ActivityRequest) *plantModels.Activity {
	return &plantModels.Activity{
		SourceType: req.SourceType,
		SourceID:   req.SourceID,
		AnimalID:   req.AnimalID,
		Type:       req.Type,
		Details:    req.Details,
		Cost:       req.Cost,
		Date:       req.Date,
		Notes:      req.Notes,
	}
}

func AnimalTypeFromRequest(req *AnimalTypeRequest) *animalModels.AnimalType {
	return &animalModels.AnimalType{
		Name:  req.Name,
		Notes: req.Notes,
	}
}

func HerdFromRequest(req *HerdRequest) *animalModels.Herd {
	return &animalModels.Herd{
		Name:             req.Name,
		AnimalTypeID:     req.AnimalTypeID,
		Location:         req.Location,
		InitialHeadCount: req.InitialHeadCount,
	}
}

func HerdActivityFromRequest(req *HerdActivityRequest) *animalModels.HerdActivity {
	return &animalModels.HerdActivity{
		ActivityType: req.ActivityType,
		Count:        req.Count,
		Date:         req.Date,
		Reason:       req.Reason,
	}
}

func InfrastructureFromRequest(req *InfrastructureRequest) *animalModels.Infrastructure {
	return &animalModels.Infrastructure{
		Type:     req.Type,
		Name:     req.Name,
		Location: req.Location,
		Cost:     req.Cost,
		Date:     req.Date,
		Notes:    req.Notes,
	}
}

func RevenueFromRequest(req *RevenueRequest) *summaryModels.Revenue {
	total := req.Total
	if total == 0 && req.Quantity > 0 && req.UnitPrice > 0 {
		total = req.Quantity * req.UnitPrice
	}
	return &summaryModels.Revenue{
		Source:    req.Source,
		SourceID:  req.SourceID,
		Type:      req.Type,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
		Total:     total,
		Date:      req.Date,
		Notes:     req.Notes,
	}
}

func CostCategoryFromRequest(req *CostCategoryRequest) *summaryModels.CostCategory {
	return &summaryModels.CostCategory{
		Name:     req.Name,
		Type:     req.Type,
		Category: req.Category,
	}
}

func UserFromRegisterRequest(req *RegisterRequest) *userModels.User {
	return &userModels.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}
}

func DerefFloat32(ptr *float32) float32 {
	if ptr == nil {
		return 0
	}
	return *ptr
}