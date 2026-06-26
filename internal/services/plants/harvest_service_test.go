package plants

import (
	plantModels "farm-backend/internal/models/plants"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupHarvestTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(
		&plantModels.Plant{},
		&plantModels.Land{},
		&plantModels.Season{},
		&plantModels.Harvest{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestHarvestService_CreateListDelete(t *testing.T) {
	db := setupHarvestTestDB(t)
	service := NewHarvestService(db)
	userID := uint(1)

	land := plantModels.Land{UserID: userID, Name: "Field A", Size: 2, Location: "North", SoilType: "Loam"}
	if err := db.Create(&land).Error; err != nil {
		t.Fatalf("create land: %v", err)
	}
	plant := plantModels.Plant{UserID: userID, Name: "Maize", Variety: "Hybrid"}
	if err := db.Create(&plant).Error; err != nil {
		t.Fatalf("create plant: %v", err)
	}
	season := plantModels.Season{
		UserID:    userID,
		Name:      "Long Rains 2026",
		PlantID:   plant.ID,
		LandID:    land.ID,
		StartDate: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC),
	}
	if err := db.Create(&season).Error; err != nil {
		t.Fatalf("create season: %v", err)
	}

	harvest := &plantModels.Harvest{
		SeasonID: season.ID,
		Quantity: 120,
		Unit:     "sacks",
		Date:     time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
		Notes:    "First picking",
	}
	if err := service.Create(userID, harvest); err != nil {
		t.Fatalf("create harvest: %v", err)
	}

	harvests, err := service.List(userID, season.ID)
	if err != nil {
		t.Fatalf("list harvests: %v", err)
	}
	if len(harvests) != 1 {
		t.Fatalf("expected 1 harvest, got %d", len(harvests))
	}
	if harvests[0].Unit != "sacks" || harvests[0].Quantity != 120 {
		t.Fatalf("unexpected harvest data: %+v", harvests[0])
	}

	if err := service.Delete(userID, harvests[0].ID); err != nil {
		t.Fatalf("delete harvest: %v", err)
	}
	harvests, err = service.List(userID, season.ID)
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(harvests) != 0 {
		t.Fatalf("expected 0 harvests after delete, got %d", len(harvests))
	}
}

func TestHarvestService_RejectsInvalidSeason(t *testing.T) {
	db := setupHarvestTestDB(t)
	service := NewHarvestService(db)

	harvest := &plantModels.Harvest{
		SeasonID: 999,
		Quantity: 10,
		Unit:     "kg",
		Date:     time.Now(),
	}
	if err := service.Create(1, harvest); err == nil {
		t.Fatal("expected error for invalid season")
	}
}