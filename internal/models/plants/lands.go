package plants

import (
	users "farm-backend/internal/models/users"

	"gorm.io/gorm"
)

type Land struct {
	gorm.Model
	UserID   uint       `gorm:"type:bigint;index:idx_land_user_id;not null"`
	User     users.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Name     string     `gorm:"type:varchar(255);not null"`
	Size     float32    `gorm:"type:decimal(10,2)"`
	Location string     `gorm:"type:varchar(255)"`
	SoilType string     `gorm:"type:varchar(100)"`
}
