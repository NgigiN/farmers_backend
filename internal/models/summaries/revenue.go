package summaries

import (
	"time"

	"gorm.io/gorm"
)

type Revenue struct {
	gorm.Model
	UserID    uint      `gorm:"index" json:"user_id"`
	Source    string    `gorm:"not null" json:"source"`
	SourceID  uint      `gorm:"index" json:"source_id"`
	Type      string    `gorm:"not null" json:"type"`
	Quantity  float64   `gorm:"not null" json:"quantity"`
	UnitPrice float64   `gorm:"not null" json:"unit_price"`
	Total     float64   `gorm:"not null" json:"total"`
	Date      time.Time `gorm:"not null" json:"date"`
	Notes     string    `gorm:"not null" json:"notes"`
}
