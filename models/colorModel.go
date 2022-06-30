package models

import (
	"time"

	"gorm.io/gorm"
)

type Color struct {
	ID        int            `gorm:"primary_key:auto_increment" json:"id"`
	ColorName *string        `gorm:"text" json:"colorName"`
	ColorCode *string        `gorm:"text" json:"colorCode"`
	CreatedAt time.Time      `gorm:"autoCreateTime; <-:create" json:"createdAt"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime; <-:update" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleteAt"`
}

func (Color) TableName() string {
	return "color"
}

type ColorInput struct {
	ColorName *string `gorm:"text" json:"colorName"`
	ColorCode *string `gorm:"text" json:"colorCode"`
}
