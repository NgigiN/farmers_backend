package models

import (
	"time"

	"gorm.io/gorm"
)

type Input struct {
	gorm.Model
	SeasonID uint `gorm:"index:idx_input_seasons"`
	Type     string
	Quantity float64
	Cost     float64
	Date     time.Time `gorm:"not null"`
	Notes    string
}
