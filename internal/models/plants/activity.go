package models

import (
	"gorm.io/gorm"
)

type Activity struct {
	gorm.Model
	UserID   uint    `gorm:"type:bigint;index:idx_activity_user_id;not null"`
	User     User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	SeasonID uint    `gorm:"type:bigint;index:idx_activity_season_id;not null"`
	Season   Season  `gorm:"foreignKey:SeasonID;constraint:OnDelete:CASCADE"`
	Type     string  `gorm:"type:varchar(100);not null;index:idx_activity_type"`
	Details  string  `gorm:"type:text;not null"`
	Cost     float32 `gorm:"type:decimal(10,2)"`
}
