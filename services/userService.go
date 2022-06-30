package services

import (
	"golang/models"
	"strings"

	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

type UserService interface {
	FindUserById(string) (*models.User, error)
	FindUserByEmail(string) (*models.User, error)
	Update(userUpdate *models.UserEdit) (*models.User, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db}
}

func (us *userService) FindUserById(id string) (*models.User, error) {
	var user *models.User
	err := us.db.Find(&user, id).Error
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (us *userService) FindUserByEmail(email string) (*models.User, error) {
	var user *models.User
	err := us.db.Where("email = ?", strings.ToLower(email)).First(&user).Error
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (us *userService) Update(userUpdate *models.UserEdit) (*models.User, error) {
	var user *models.User
	err := smapping.FillStruct(&user, smapping.MapFields(&userUpdate))
	if err != nil {
		return &models.User{}, err
	}

	err = us.db.Save(&user).Error
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}
