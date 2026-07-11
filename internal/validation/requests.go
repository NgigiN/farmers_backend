package validation

import "time"

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,max=100,safe_name"`
	LastName  string `json:"last_name"  validate:"omitempty,max=100,safe_name"`
	FarmName  string `json:"farm_name"  validate:"required,max=255,safe_name"`
	Location  string `json:"location"   validate:"required,max=255,safe_location"`
}

func (r *UpdateProfileRequest) Sanitize() {
	r.FirstName = SanitizeText(r.FirstName)
	r.LastName = SanitizeText(r.LastName)
	r.FarmName = SanitizeText(r.FarmName)
	r.Location = SanitizeText(r.Location)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type RegisterRequest struct {
	Email     string `json:"email"      validate:"required,email,max=255"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100,safe_name"`
	LastName  string `json:"last_name"  validate:"required,min=1,max=100,safe_name"`
	Password  string `json:"password"   validate:"required,min=8,max=255"`
}

func (r *RegisterRequest) Sanitize() {
	r.FirstName = SanitizeText(r.FirstName)
	r.LastName = SanitizeText(r.LastName)
}

type LandRequest struct {
	Name     string   `json:"name"      validate:"required,min=1,max=255,safe_name"`
	Size     *float32 `json:"size"      validate:"omitempty,gte=0,lte=9999999"`
	Location string   `json:"location"  validate:"omitempty,max=255,safe_location"`
	SoilType string   `json:"soil_type" validate:"omitempty,max=100,safe_name"`
}

func (r *LandRequest) Sanitize() {
	r.Name = SanitizeText(r.Name)
	r.Location = SanitizeOptionalText(r.Location)
	r.SoilType = SanitizeOptionalText(r.SoilType)
}

type PlantRequest struct {
	Name    string `json:"name"    validate:"required,min=1,max=255,safe_name"`
	Variety string `json:"variety" validate:"omitempty,max=255,safe_name"`
}

func (r *PlantRequest) Sanitize() {
	r.Name = SanitizeText(r.Name)
	r.Variety = SanitizeOptionalText(r.Variety)
}

type SeasonRequest struct {
	Name      string    `json:"name"       validate:"required,min=1,max=255,safe_name"`
	PlantID   uint      `json:"plant_id"   validate:"required,gt=0"`
	LandID    uint      `json:"land_id"    validate:"required,gt=0"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date"`
}

func (r *SeasonRequest) Sanitize() {
	r.Name = SanitizeText(r.Name)
}

type HarvestRequest struct {
	SeasonID uint      `json:"season_id" validate:"required,gt=0"`
	Quantity float64   `json:"quantity"  validate:"required,gt=0,lte=9999999"`
	Unit     string    `json:"unit"      validate:"required,min=1,max=50,safe_name"`
	Date     time.Time `json:"date"      validate:"required"`
	Notes    string    `json:"notes"     validate:"omitempty,max=1000"`
}

func (r *HarvestRequest) Sanitize() {
	r.Unit = SanitizeText(r.Unit)
	r.Notes = SanitizeOptionalText(r.Notes)
}

type InputRequest struct {
	SourceType string    `json:"source_type" validate:"required,oneof=plant animal"`
	SourceID   uint      `json:"source_id"   validate:"required,gt=0"`
	AnimalID   uint      `json:"animal_id"`
	Type       string    `json:"type"        validate:"required,min=1,max=100,safe_name"`
	Quantity   float64   `json:"quantity"    validate:"gte=0,lte=9999999"`
	Cost       float64   `json:"cost"        validate:"required,gt=0,lte=9999999"`
	Date       time.Time `json:"date"        validate:"required"`
	Notes      string    `json:"notes"       validate:"omitempty,max=1000"`
}

func (r *InputRequest) Sanitize() {
	r.Type = SanitizeText(r.Type)
	r.Notes = SanitizeOptionalText(r.Notes)
}

type ActivityRequest struct {
	SourceType string    `json:"source_type" validate:"required,oneof=plant animal"`
	SourceID   uint      `json:"source_id"   validate:"required,gt=0"`
	AnimalID   uint      `json:"animal_id"`
	Type       string    `json:"type"        validate:"required,min=1,max=100,safe_name"`
	Details    string    `json:"details"     validate:"omitempty,max=1000"`
	Cost       float64   `json:"cost"        validate:"gte=0,lte=9999999"`
	Date       time.Time `json:"date"        validate:"required"`
	Notes      string    `json:"notes"       validate:"omitempty,max=1000"`
}

func (r *ActivityRequest) Sanitize() {
	r.Type = SanitizeText(r.Type)
	r.Details = SanitizeOptionalText(r.Details)
	r.Notes = SanitizeOptionalText(r.Notes)
}

type AnimalTypeRequest struct {
	Name  string `json:"name"  validate:"required,min=1,max=100,safe_name"`
	Notes string `json:"notes" validate:"omitempty,max=1000"`
}

func (r *AnimalTypeRequest) Sanitize() {
	r.Name = SanitizeText(r.Name)
	r.Notes = SanitizeOptionalText(r.Notes)
}

type HerdRequest struct {
	Name             string `json:"name"               validate:"required,min=1,max=100,safe_name"`
	AnimalTypeID     uint   `json:"animal_type_id"     validate:"required,gt=0"`
	Location         string `json:"location"           validate:"required,min=1,max=255,safe_location"`
	InitialHeadCount int    `json:"initial_head_count" validate:"required,gt=0,lte=999999"`
}

func (r *HerdRequest) Sanitize() {
	r.Name = SanitizeText(r.Name)
	r.Location = SanitizeText(r.Location)
}

type HerdActivityRequest struct {
	ActivityType string    `json:"activity_type" validate:"required,oneof=birth fatality"`
	Count        int       `json:"count"         validate:"required,gt=0,lte=999999"`
	Date         time.Time `json:"date"          validate:"required"`
	Reason       string    `json:"reason"        validate:"omitempty,max=1000"`
}

func (r *HerdActivityRequest) Sanitize() {
	r.Reason = SanitizeOptionalText(r.Reason)
}

type InfrastructureRequest struct {
	Type     string    `json:"type"     validate:"required,oneof=Store House Fence Barn Greenhouse Other"`
	Name     string    `json:"name"     validate:"required,min=1,max=100,safe_name"`
	Location string    `json:"location" validate:"required,min=1,max=255,safe_location"`
	Cost     float64   `json:"cost"     validate:"gte=0,lte=9999999"`
	Date     time.Time `json:"date"     validate:"required"`
	Notes    string    `json:"notes"    validate:"omitempty,max=1000"`
}

func (r *InfrastructureRequest) Sanitize() {
	r.Name = SanitizeText(r.Name)
	r.Location = SanitizeText(r.Location)
	r.Notes = SanitizeOptionalText(r.Notes)
}

type RevenueRequest struct {
	Source    string    `json:"source"     validate:"required,oneof=plant animal"`
	SourceID  uint      `json:"source_id"  validate:"required,gt=0"`
	Type      string    `json:"type"       validate:"required,min=1,max=100,safe_name"`
	Quantity  float64   `json:"quantity"   validate:"required,gt=0,lte=9999999"`
	UnitPrice float64   `json:"unit_price" validate:"required,gt=0,lte=9999999"`
	Total     float64   `json:"total"`
	Date      time.Time `json:"date"       validate:"required"`
	Notes     string    `json:"notes"      validate:"omitempty,max=1000"`
}

func (r *RevenueRequest) Sanitize() {
	r.Type = SanitizeText(r.Type)
	r.Notes = SanitizeOptionalText(r.Notes)
}

type CostCategoryRequest struct {
	Name     string `json:"name"     validate:"required,min=1,max=100,safe_name"`
	Type     string `json:"type"     validate:"required,oneof=plant animal"`
	Category string `json:"category" validate:"required,oneof=input activity infrastructure"`
}

func (r *CostCategoryRequest) Sanitize() {
	r.Name = SanitizeText(r.Name)
}