package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        uint   `gorm:"PrimaryKey"`
	Email     string `gorm:"unique;not null"`
	FirstName string `gorm:"unique;not null"`
	LastName  string `gorm:"not null"`
	Password  string
	FarmName  string
	Location  string
}
