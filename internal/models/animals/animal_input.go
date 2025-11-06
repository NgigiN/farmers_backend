package animals

import (
	"time"

	"gorm.io/gorm"
)

type AnimalInput struct {
	gorm.Model
	UserID   uint   `gorm:"index"`
	HerdID   uint   `gorm:"index"`
	AnimalID uint   `gorm:"index"`
	Type     string `gorm:"not null"` // Food, vaccination, Labor
	Quantity float64
	Cost     float64   `gorm:"not null"`
	Date     time.Time `gorm:"not null"`
	Notes    string
}
