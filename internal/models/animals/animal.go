package animals

import (
	"time"

	"gorm.io/gorm"
)

type Animal struct {
	gorm.Model
	UserID       uint `gorm:"index"`
	AnimalTypeID uint `gorm:"index"`
	HerdID       uint `gorm:"index"`
	Name         string
	BirthDate    time.Time
}
