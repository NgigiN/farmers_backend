package models

import (
	"gorm.io/gorm"
)

type Crop struct {
	gorm.Model
	UserID  uint   `gorm:"type:bigint;index:idx_crop_user_id;not null"`
	User    User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Name    string `gorm:"type:varchar(255);not null"`
	Variety string `gorm:"type:varchar(255)"`
}
