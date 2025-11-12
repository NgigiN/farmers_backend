package animals

import (
	"time"

	"gorm.io/gorm"
)

type Animal struct {
	gorm.Model
	UserID       uint      `gorm:"index" json:"user_id"`
	AnimalTypeID uint      `gorm:"index" json:"animal_type_id"`
	HerdID       uint      `gorm:"index" json:"herd_id"`
	Name         string    `gorm:"not null" json:"name"`
	BirthDate    time.Time `gorm:"not null" json:"birth_date"`
}
