package animals

import (
	"gorm.io/gorm"
)

type AnimalType struct {
	gorm.Model
	UserID uint   `gorm:"index" json:"user_id"`
	Name   string `gorm:"not null" json:"name"` //Cows, Poultry, Goats
	Notes  string `gorm:"not null" json:"notes"`
}
