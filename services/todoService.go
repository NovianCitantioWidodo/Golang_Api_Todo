package services

import (
	"golang/models"

	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

type TodoService interface {
	All(int) ([]*models.Todo, error)
	FindByID(todoID int, userID int) (models.Todo, error)
	Insert(t models.TodoInput) (models.Todo, error)
	Update(todoID int, t models.TodoInput) (models.Todo, error)
	Delete(t models.Todo) error
	IsAllowed(userID int, todoID int) bool
}

type todoService struct {
	db *gorm.DB
}

func NewTodoService(db *gorm.DB) TodoService {
	return &todoService{db}
}

func (ts *todoService) All(userID int) ([]*models.Todo, error) {
	var todos []*models.Todo
	err := ts.db.Preload("User").Preload("Color").Where("user_id = ?", userID).Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (ts *todoService) FindByID(todoID int, userID int) (models.Todo, error) {
	var todo models.Todo
	err := ts.db.Preload("User").Preload("Color").Find(&todo, todoID).Error
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (ts *todoService) Insert(t models.TodoInput) (models.Todo, error) {
	var todo models.Todo
	err := smapping.FillStruct(&todo, smapping.MapFields(&t))
	if err != nil {
		return models.Todo{}, err
	}
	err = ts.db.Save(&todo).Error
	if err != nil {
		return models.Todo{}, err
	}

	err = ts.db.Preload("User").Preload("Color").Find(&todo).Error
	if err != nil {
		return models.Todo{}, err
	}

	return todo, nil
}

func (ts *todoService) Update(todoID int, t models.TodoInput) (models.Todo, error) {
	var todo models.Todo
	err := smapping.FillStruct(&todo, smapping.MapFields(&t))
	if err != nil {
		return models.Todo{}, err
	}

	err = ts.db.Preload("User").Preload("Color").Find(&todo, todoID).Error
	if err != nil {
		return models.Todo{}, err
	}

	todo.ID = todoID
	err = ts.db.Save(&todo).Error
	if err != nil {
		return models.Todo{}, err
	}

	err = ts.db.Preload("User").Preload("Color").Find(&todo, todoID).Error
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (ts *todoService) Delete(t models.Todo) error {
	err := ts.db.Delete(&t).Error
	if err != nil {
		return err
	}
	return nil
}

func (ts *todoService) IsAllowed(userID int, todoID int) bool {
	var todo models.Todo
	ts.db.Debug().Find(&todo, todoID)
	return userID == todo.UserID
}
