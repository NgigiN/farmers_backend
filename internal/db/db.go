// Package db provides the db methods to connect and use the database.
package db

import (
	"farm-backend/internal/config"
	plantModels "farm-backend/internal/models/plants"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&plantModels.User{}, &plantModels.Crop{}, &plantModels.Land{}, &plantModels.Season{}, &plantModels.Input{}, &plantModels.Activity{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
