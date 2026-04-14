package users

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string `gorm:"type:varchar(255);uniqueIndex:idx_user_email;not null" json:"email"      validate:"required,email"`
	FirstName string `gorm:"type:varchar(100);not null"                           json:"first_name" validate:"required,min=1"`
	LastName  string `gorm:"type:varchar(100);not null"                           json:"last_name"  validate:"required,min=1"`
	// Password is excluded from all JSON responses with json:"-"
	Password string `gorm:"type:varchar(255);not null" json:"-" validate:"required,min=8"`
	FarmName string `gorm:"type:varchar(255)"          json:"farm_name,omitempty"`
	Location string `gorm:"type:varchar(255)"          json:"location,omitempty"`
}
