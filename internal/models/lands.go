package models

import (
	"gorm.io/gorm"
)

type Land struct {
	gorm.Model
	UserID   uint
	Name     string
	Size     float32
	Location string
	SoilType string
}
