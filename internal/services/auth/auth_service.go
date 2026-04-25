// Package auth provides authentication logic: user registration and JWT-based login.
package auth

import (
	"context"
	"errors"
	"time"

	"farm-backend/internal/config"
	"farm-backend/internal/middleware"
	users "farm-backend/internal/models/users"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
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

func (s *Service) GoogleLogin(idTokenStr string) (string, error) {
	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, idTokenStr, s.Cfg.GoogleClientID)
	if err != nil {
		return "", errors.New("invalid google id token: " + err.Error())
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return "", errors.New("email not found in token claims")
	}

	var user users.User
	err = s.DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User doesn't exist, create a new one
			firstName, _ := payload.Claims["given_name"].(string)
			lastName, _ := payload.Claims["family_name"].(string)
			profilePic, _ := payload.Claims["picture"].(string)

			user = users.User{
				Email:          email,
				FirstName:      firstName,
				LastName:       lastName,
				GoogleID:       payload.Subject,
				ProfilePicture: profilePic,
			}
			if createErr := s.DB.Create(&user).Error; createErr != nil {
				return "", createErr
			}
		} else {
			return "", err
		}
	} else if user.GoogleID == "" {
		// Update user with GoogleID and ProfilePicture if missing
		updates := map[string]interface{}{
			"google_id": payload.Subject,
		}
		if user.ProfilePicture == "" {
			updates["profile_picture"], _ = payload.Claims["picture"].(string)
		}
		if updateErr := s.DB.Model(&user).Updates(updates).Error; updateErr != nil {
			return "", updateErr
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(s.Cfg.JWTSecret))
}
