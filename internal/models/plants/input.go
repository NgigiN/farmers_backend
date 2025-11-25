package plants

import (
	users "farm-backend/internal/models/users"
	"time"

	"gorm.io/gorm"
)

type Input struct {
	gorm.Model
	UserID     uint       `gorm:"type:bigint;index:idx_input_user_id;not null" json:"user_id"`
	User       users.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	SourceType string     `gorm:"type:varchar(20);not null;index:idx_input_source_type" json:"source_type"` // "plant" or "animal"
	SourceID   uint       `gorm:"type:bigint;index:idx_input_source_id;not null" json:"source_id"`          // SeasonID if plant, HerdID if animal
	AnimalID   uint       `gorm:"type:bigint;index:idx_input_animal_id" json:"animal_id"`
	Type       string     `gorm:"type:varchar(100);not null;index:idx_input_type" json:"type"`
	Quantity   float64    `gorm:"type:decimal(10,2);not null" json:"quantity"`
	Cost       float64    `gorm:"type:decimal(10,2);not null" json:"cost"`
	Date       time.Time  `gorm:"type:date;not null;index:idx_input_date" json:"date"`
	Notes      string     `gorm:"type:text" json:"notes"`
}
