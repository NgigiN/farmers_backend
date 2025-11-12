package animals

import (
	"gorm.io/gorm"
)

type Herd struct {
	gorm.Model
	UserID       uint   `gorm:"index" json:"user_id"`
	Name         string `gorm:"not null" json:"name"`
	AnimalTypeID uint   `gorm:"index" json:"animal_type_id"`
	Location     string `gorm:"not null" json:"location"`
}
