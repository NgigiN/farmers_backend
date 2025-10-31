// Package models provides structs that form the data layer
// for this application
package models

import (
	"gorm.io/gorm"
)

type Activity struct {
	gorm.Model
	SeasonID uint
	Type     string
	Details  string
	Cost     float32
}
