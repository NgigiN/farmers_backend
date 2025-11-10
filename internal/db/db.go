// Package db provides the db methods to connect and use the database.
package db

import (
	"farm-backend/internal/config"

	animalModels "farm-backend/internal/models/animals"
	plantModels "farm-backend/internal/models/plants"
	summaryModels "farm-backend/internal/models/summaries"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&plantModels.User{}, &plantModels.Crop{}, &plantModels.Land{}, &plantModels.Season{}, &plantModels.Input{}, &plantModels.Activity{}, &animalModels.AnimalType{}, &animalModels.Animal{}, &animalModels.Herd{}, &animalModels.Infrastructure{}, &summaryModels.CostCategory{}, &summaryModels.Revenue{})

	if err != nil {
		return nil, err
	}
	return db, nil
}
