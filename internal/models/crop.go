package models

import (
	"gorm.io/gorm"
)

type Crop struct {
	gorm.Model
	UserID  uint   `gorm:"index:idx_user_crops"`
	Name    string `gorm:"not null"`
	Variety string
}
