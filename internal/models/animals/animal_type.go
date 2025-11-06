package animals

import (
	"gorm.io/gorm"
)

type AnimalType struct {
	gorm.Model
	UserID uint   `gorm:"index"`
	Name   string `gorm:"not null"` //Cows, Poultry, Goats
	Notes  string
}
