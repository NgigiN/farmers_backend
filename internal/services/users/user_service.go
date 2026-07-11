package users

import (
	"errors"
	"farm-backend/internal/validation"
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

func (s *UserService) GetProfile(userID uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateProfile(userID uint, req *validation.UpdateProfileRequest) error {
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}
	req.Sanitize()

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
	user.FarmName = req.FarmName
	user.Location = req.Location

	return s.DB.Save(&user).Error
}

func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashed)
	return s.DB.Save(&user).Error
}