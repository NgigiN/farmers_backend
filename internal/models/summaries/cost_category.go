package summaries

import (
	"gorm.io/gorm"
)

type CostCategory struct {
	gorm.Model
	UserID    uint   `gorm:"index" json:"user_id"`
	Name      string `gorm:"not null" json:"name"`
	Type      string `gorm:"not null" json:"type"`     // "plant" or "animal"
	Category  string `gorm:"not null" json:"category"` // "input", "activity", or "infrastructure"
	IsDefault bool   `gorm:"not null" json:"is_default"`
}
