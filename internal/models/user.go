package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string
	FirstName string
	LastName  string
	Password  []byte
	FarmName  string
	Location  string
}
