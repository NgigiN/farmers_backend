// Package models provides structs that form the data layer
// for this application
package models

import (
	"gorm.io/gorm"
)

type Activity struct {
	gorm.Model
	UserID   uint
	SeasonID uint `gorm:"index:idx_season_activities"`
	Type     string
	Details  string `gorm:"not null"`
	Cost     float32
}
