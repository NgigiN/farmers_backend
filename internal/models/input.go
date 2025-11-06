package models

import (
	"time"

	"gorm.io/gorm"
)

type Input struct {
	gorm.Model
	SeasonID uint      `gorm:"type:bigint;index:idx_input_season_id;not null"`
	Season   Season    `gorm:"foreignKey:SeasonID;constraint:OnDelete:CASCADE"`
	Type     string    `gorm:"type:varchar(100);not null"`
	Quantity float64   `gorm:"type:decimal(10,2);not null"`
	Cost     float64   `gorm:"type:decimal(10,2);not null"`
	Date     time.Time `gorm:"type:date;not null;index:idx_input_date"`
	Notes    string    `gorm:"type:text"`
}
