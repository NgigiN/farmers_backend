package animals

import (
	"gorm.io/gorm"
)

type Herd struct {
	gorm.Model
	UserID       uint   `gorm:"index"`
	Name         string `gorm:"not null"`
	AnimalTypeID uint
	Location     string
}
