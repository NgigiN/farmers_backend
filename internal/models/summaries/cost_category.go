package summaries

import (
	"gorm.io/gorm"
)

type CostCategory struct {
	gorm.Model
	UserID    uint   `gorm:"index"`
	Name      string `gorm:"not null"`
	Type      string `gorm:"not null"` // plant v Animal
	IsDefault bool
}
