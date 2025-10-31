package models

import (
	"gorm.io/gorm"
)

type Land struct {
	gorm.Model
	UserID   uint   `gorm:"index:idx_lands_user"`
	Name     string `gorm:"not null"`
	Size     float32
	Location string
	SoilType string
}
