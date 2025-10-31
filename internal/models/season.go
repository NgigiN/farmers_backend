package models

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	gorm.Model
	UserID    uint
	Name      string
	CropID    uint
	LandID    uint
	StartDate time.Time
	EndDate   time.Time `gorm:"default:null"`
}
