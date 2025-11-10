package models

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	gorm.Model
	UserID     uint      `gorm:"type:bigint;index:idx_activity_user_id;not null"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	SourceType string    `gorm:"type:varchar(20);not null;index:idx_activity_source_type"` // "plant" or "animal"
	SourceID   uint      `gorm:"type:bigint;index:idx_activity_source_id;not null"`        // SeasonID if plant, HerdID if animal
	AnimalID   uint      `gorm:"type:bigint;index:idx_activity_animal_id"`
	Type       string    `gorm:"type:varchar(100);not null;index:idx_activity_type"`
	Details    string    `gorm:"type:text;not null"`
	Cost       float64   `gorm:"type:decimal(10,2)"`
	Date       time.Time `gorm:"type:date;index:idx_activity_date"`
	Notes      string    `gorm:"type:text"`
}
