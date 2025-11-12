package animals

import (
	"time"

	"gorm.io/gorm"
)

type AnimalActivity struct {
	gorm.Model
	UserID   uint      `gorm:"index" json:"user_id"`
	HerdID   uint      `gorm:"index" json:"herd_id"`
	AnimalID uint      `gorm:"index" json:"animal_id"`
	Type     string    `gorm:"not null" json:"type"`
	Details  string    `gorm:"not null" json:"details"`
	Cost     float64   `gorm:"not null" json:"cost"`
	Date     time.Time `gorm:"not null" json:"date"`
	Notes    string    `gorm:"not null" json:"notes"`
}
