package models

import (
	"time"

	"gorm.io/gorm"
)

type Input struct {
	gorm.Model
	UserID     uint      `gorm:"type:bigint;index:idx_input_user_id;not null"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	SourceType string    `gorm:"type:varchar(20);not null;index:idx_input_source_type"` // "plant" or "animal"
	SourceID   uint      `gorm:"type:bigint;index:idx_input_source_id;not null"`        // SeasonID if plant, HerdID if animal
	AnimalID   uint      `gorm:"type:bigint;index:idx_input_animal_id"`
	Type       string    `gorm:"type:varchar(100);not null;index:idx_input_type"`
	Quantity   float64   `gorm:"type:decimal(10,2);not null"`
	Cost       float64   `gorm:"type:decimal(10,2);not null"`
	Date       time.Time `gorm:"type:date;not null;index:idx_input_date"`
	Notes      string    `gorm:"type:text"`
}
