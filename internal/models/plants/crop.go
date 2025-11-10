package plants

import (
	users "farm-backend/internal/models/users"

	"gorm.io/gorm"
)

type Plant struct {
	gorm.Model
	UserID  uint       `gorm:"type:bigint;index:idx_plant_user_id;not null"`
	User    users.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Name    string     `gorm:"type:varchar(255);not null"`
	Variety string     `gorm:"type:varchar(255)"`
}
