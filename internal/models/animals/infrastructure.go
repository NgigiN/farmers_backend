package animals

import (
	"time"

	"gorm.io/gorm"
)

type Infrastructure struct {
	gorm.Model
	UserID   uint      `gorm:"index" json:"user_id"`
	Type     string    `gorm:"not null"` // store, House, Fence
	Name     string    `gorm:"not null" json:"name"`
	Location string    `gorm:"not null" json:"location"`
	Cost     float64   `gorm:"not null" json:"cost"`
	Date     time.Time `gorm:"not null" json:"date"`
	Notes    string    `gorm:"not null" json:"notes"`
}
