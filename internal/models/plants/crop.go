package plants

import (
	users "farm-backend/internal/models/users"

	"gorm.io/gorm"
)

type Plant struct {
	gorm.Model
	UserID  uint       `gorm:"type:bigint;index:idx_plant_user_id;not null" json:"user_id"`
	User    users.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Name    string     `gorm:"type:varchar(255);not null" json:"name"`
	Variety string     `gorm:"type:varchar(255)" json:"variety"`
}
