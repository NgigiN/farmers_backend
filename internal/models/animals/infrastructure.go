package animals

import (
	"time"

	"gorm.io/gorm"
)

type Infrastructure struct {
	gorm.Model
	UserID   uint   `gorm:"index"`
	Type     string `gorm:"not null"` // store, House, Fence
	Name     string `gorm:"not null"`
	Location string
	Cost     float64   `gorm:"not null"`
	Date     time.Time `gorm:"not null"`
	Notes    string
}
