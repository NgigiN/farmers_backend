package animals

import (
	"time"

	"gorm.io/gorm"
)

type AnimalActivity struct {
	gorm.Model
	UserID   uint `gorm:"index"`
	HerdID   uint `gorm:"index"`
	AnimalID uint `gorm:"index"`
	Type     string
	Details  string
	Cost     float64
	Date     time.Time `gorm:"not null"`
}
