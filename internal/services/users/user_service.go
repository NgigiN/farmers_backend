package users

import (
	"errors"
	"farm-backend/internal/middleware"
	models "farm-backend/internal/models/users"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

type UpdateProfileDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FarmName  string `json:"farm_name"`
	Location  string `json:"location"`
}

func (s *UserService) GetProfile(userID uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateProfile(userID uint, req *UpdateProfileDTO) error {
	if err := middleware.ValidateStruct(req); err != nil {
		return err
	}
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.FarmName != "" {
		user.FarmName = req.FarmName
	}
	if req.Location != "" {
		user.Location = req.Location
	}

	return s.DB.Save(&user).Error
}

func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashed)
	return s.DB.Save(&user).Error
}
