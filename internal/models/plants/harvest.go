package plants

import (
	users "farm-backend/internal/models/users"
	"time"

	"gorm.io/gorm"
)

type Harvest struct {
	gorm.Model
	UserID    uint        `gorm:"type:bigint;index:idx_harvest_user_id;not null" json:"user_id"`
	User      *users.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	SeasonID  uint        `gorm:"type:bigint;index:idx_harvest_season_id;not null" json:"season_id"`
	Season    *Season     `gorm:"foreignKey:SeasonID;constraint:OnDelete:RESTRICT" json:"-"`
	Quantity  float64     `gorm:"type:decimal(10,2);not null" json:"quantity" validate:"required,gt=0"`
	Unit      string      `gorm:"type:varchar(50);not null" json:"unit" validate:"required,min=1"`
	Date      time.Time   `gorm:"type:date;not null;index:idx_harvest_date" json:"date"`
	Notes     string      `gorm:"type:text" json:"notes"`
	RevenueID *uint       `gorm:"type:bigint;index:idx_harvest_revenue_id" json:"revenue_id,omitempty"`
}