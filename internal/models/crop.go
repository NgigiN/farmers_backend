package models

import (
	"gorm.io/gorm"
)

type Crop struct {
	gorm.Model
	UserID  uint
	Name    string
	Variety string
}
