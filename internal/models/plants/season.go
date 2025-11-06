package models

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	gorm.Model
	UserID    uint      `gorm:"type:bigint;index:idx_season_user_id;not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Name      string    `gorm:"type:varchar(255);not null"`
	CropID    uint      `gorm:"type:bigint;index:idx_season_crop_id;not null"`
	Crop      Crop      `gorm:"foreignKey:CropID;constraint:OnDelete:RESTRICT"`
	LandID    uint      `gorm:"type:bigint;index:idx_season_land_id;not null"`
	Land      Land      `gorm:"foreignKey:LandID;constraint:OnDelete:RESTRICT"`
	StartDate time.Time `gorm:"type:date;not null;index:idx_season_dates"`
	EndDate   time.Time `gorm:"type:date;index:idx_season_dates"`
}
