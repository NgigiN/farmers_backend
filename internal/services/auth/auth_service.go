// Package auth provides authentication logic: user registration and JWT-based login.
package auth

import (
	"time"

	"farm-backend/internal/config"
	"farm-backend/internal/middleware"
	users "farm-backend/internal/models/users"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewService(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{DB: db, Cfg: cfg}
}

func (s *Service) Register(user *users.User) error {
	if err := middleware.ValidateStruct(user); err != nil {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return s.DB.Create(user).Error
}

func (s *Service) Login(email, password string) (string, error) {
	var user users.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(s.Cfg.JWTSecret))
}
