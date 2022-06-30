package services

import (
	"golang/models"

	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

type ColorService interface {
	All() ([]*models.Color, error)
	FindByID(id int) (models.Color, error)
	Insert(colorInput models.ColorInput) (models.Color, error)
	Update(id int, colorInput models.ColorInput) (models.Color, error)
	Delete(color models.Color) error
}

type colorService struct {
	db *gorm.DB
}

func NewColorService(db *gorm.DB) ColorService {
	return &colorService{db}
}

func (cs *colorService) All() ([]*models.Color, error) {
	var colors []*models.Color
	err := cs.db.Find(&colors).Error
	if err != nil {
		return nil, err
	}

	return colors, nil
}

func (cs *colorService) FindByID(id int) (models.Color, error) {
	var color models.Color
	err := cs.db.Find(&color, id).Error
	if err != nil {
		return models.Color{}, err
	}
	return color, nil
}

func (cs *colorService) Insert(colorInput models.ColorInput) (models.Color, error) {
	var color models.Color
	err := smapping.FillStruct(&color, smapping.MapFields(&colorInput))
	if err != nil {
		return models.Color{}, err
	}
	err = cs.db.Save(&color).Error
	if err != nil {
		return models.Color{}, err
	}

	err = cs.db.Find(&color).Error
	if err != nil {
		return models.Color{}, err
	}
	return color, nil
}

func (cs *colorService) Update(id int, colorInput models.ColorInput) (models.Color, error) {
	var color models.Color
	err := smapping.FillStruct(&color, smapping.MapFields(&colorInput))
	if err != nil {
		return models.Color{}, err
	}

	color.ID = id
	err = cs.db.Debug().Save(&color).Error
	if err != nil {
		return models.Color{}, err
	}

	err = cs.db.Find(&color, id).Error
	if err != nil {
		return models.Color{}, err
	}
	return color, nil
}

func (cs *colorService) Delete(color models.Color) error {
	err := cs.db.Delete(&color).Error
	if err != nil {
		return err
	}
	return nil
}
