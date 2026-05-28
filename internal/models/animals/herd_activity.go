package animals

import (
	"time"

	"gorm.io/gorm"
)

type HerdActivity struct {
	gorm.Model
	HerdID       uint      `gorm:"index" json:"herd_id"`
	ActivityType string    `gorm:"not null" json:"activity_type"` // e.g., "birth" or "fatality"
	Count        int       `gorm:"not null" json:"count"`
	Date         time.Time `gorm:"not null" json:"date"`
	Reason       string    `json:"reason"`
}
