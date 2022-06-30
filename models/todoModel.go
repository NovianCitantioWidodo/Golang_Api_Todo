package models

import (
	// "database/sql"
	// "encoding/json"

	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID        int            `gorm:"primary_key:auto_increment" json:"id"`
	Title     string         `gorm:"text" json:"title"`
	Isi       string         `gorm:"text" json:"isi"`
	Reminder  *time.Time     `json:"reminder"`
	CreatedAt time.Time      `gorm:"autoCreateTime; <-:create" json:"createdAt"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime; <-:update" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleteAt"`
	ColorID   *int           `json:"-"`
	Color     *Color         `gorm:"foreignkey:ColorID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"color"`
	UserID    int            `gorm:"not null" json:"userId"`
	User      User           `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
}

func (Todo) TableName() string {
	return "todo"
}

type TodoInput struct {
	Title    string     `json:"title" form:"title" binding:"required"`
	Isi      string     `json:"isi" form:"title" binding:"required"`
	Reminder *time.Time `json:"reminder,omitempty" form:"reminder,omitempty"`
	ColorID  *int       `json:"colorId,omitempty"  form:"colorId,omitempty"`
	UserID   int        `json:"userId"  form:"userId"`
}
