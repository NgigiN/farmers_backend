package users

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string `gorm:"type:varchar(255);uniqueIndex:idx_user_email;not null" json:"email"`
	FirstName string `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string `gorm:"type:varchar(100);not null" json:"last_name"`
	Password  string `gorm:"type:varchar(255);not null" json:"password"`
	FarmName  string `gorm:"type:varchar(255)" json:"farm_name"`
	Location  string `gorm:"type:varchar(255)" json:"location"`
}
