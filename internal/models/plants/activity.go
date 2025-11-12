package plants

import (
	users "farm-backend/internal/models/users"
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	gorm.Model
	UserID     uint       `gorm:"type:bigint;index:idx_activity_user_id;not null" json:"user_id"`
	User       users.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	SourceType string     `gorm:"type:varchar(20);not null;index:idx_activity_source_type" json:"source_type"` // "plant" or "animal"
	SourceID   uint       `gorm:"type:bigint;index:idx_activity_source_id;not null" json:"source_id"`          // SeasonID if plant, HerdID if animal
	AnimalID   uint       `gorm:"type:bigint;index:idx_activity_animal_id" json:"animal_id"`
	Type       string     `gorm:"type:varchar(100);not null;index:idx_activity_type" json:"type"`
	Details    string     `gorm:"type:text;not null" json:"details"`
	Cost       float64    `gorm:"type:decimal(10,2)" json:"cost"`
	Date       time.Time  `gorm:"type:date;index:idx_activity_date" json:"date"`
	Notes      string     `gorm:"type:text" json:"notes"`
}
