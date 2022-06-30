package services

import (
	"strings"
	"time"

	"golang/models"
	"golang/utils"

	"gorm.io/gorm"
)

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.User, error)
}

type AuthServiceImpl struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return &AuthServiceImpl{db}
}

func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.User, error) {
	hashedPassword, _ := utils.HashPassword(user.Password)

	newUser := models.User{
		Name:             user.Name,
		Email:            strings.ToLower(user.Email),
		Verified:         false,
		VerificationCode: user.VerificationCode,
		Role:             "user",
		Password:         hashedPassword,
		CreatedAt:        time.Now(),
		UpdatedAt:        user.CreatedAt,
	}

	err := uc.db.Create(&newUser).Error
	if err != nil {
		return nil, err
	}

	err = uc.db.Find(&newUser).Error
	if err != nil {
		return nil, err
	}
	return &newUser, err
}
