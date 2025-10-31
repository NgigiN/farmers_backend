package models

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	gorm.Model
	UserID    uint      `gorm:"index:idx_season_user"`
	Name      string    `gorm:"not null"`
	CropID    uint      `gorm:"idx_season_user"`
	LandID    uint      `gorm:"idx_season_user"`
	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"default:null"`
}
