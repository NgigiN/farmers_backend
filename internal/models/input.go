package models

import (
	"time"

	"gorm.io/gorm"
)

type Input struct {
	gorm.Model
	SeasonID uint
	Type     string
	Quantity float64
	Cost     float64
	Date     time.Time
	Notes    string
}
