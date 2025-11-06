package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string `gorm:"type:varchar(255);uniqueIndex:idx_user_email;not null"`
	FirstName string `gorm:"type:varchar(100);not null"`
	LastName  string `gorm:"type:varchar(100);not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	FarmName  string `gorm:"type:varchar(255)"`
	Location  string `gorm:"type:varchar(255)"`
}
