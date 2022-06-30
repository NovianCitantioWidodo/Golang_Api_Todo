package controllers

import (
	"net/http"
	"strconv"

	"golang/helper"
	"golang/models"
	"golang/services"

	"github.com/gin-gonic/gin"
)

type TodoController struct {
	todoService services.TodoService
}

func NewTodoController(todoService services.TodoService) TodoController {
	return TodoController{todoService}
}

func (tc *TodoController) List(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.User)
	userID := currentUser.ID
	var todos []*models.Todo
	todos, err := tc.todoService.All(userID)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}

	response := helper.BuildResponse("OK", todos)
	ctx.JSON(http.StatusOK, response)
}

func (tc *TodoController) FindByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(*models.User)
	userID := currentUser.ID

	var todo models.Todo
	todo, err = tc.todoService.FindByID(id, userID)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}

	if todo.ID != id {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	if tc.todoService.IsAllowed(userID, id) {
		response := helper.BuildResponse("OK", todo)
		ctx.JSON(http.StatusOK, response)
		return
	} else {
		response := helper.BuildResponse("You are not the owner", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
		return
	}
}

func (tc *TodoController) Insert(ctx *gin.Context) {
	var todoCreate models.TodoInput
	errDTO := ctx.ShouldBind(&todoCreate)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	} else {
		currentUser := ctx.MustGet("currentUser").(*models.User)
		todoCreate.UserID = currentUser.ID
		result, err := tc.todoService.Insert(todoCreate)
		if err != nil {
			response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
			return
		}
		response := helper.BuildResponse("OK", result)
		ctx.JSON(http.StatusCreated, response)
		return
	}
}

func (tc *TodoController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(*models.User)
	userID := currentUser.ID

	var todo models.Todo
	todo, err = tc.todoService.FindByID(id, userID)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}
	if todo.ID != id {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	var todoUpdate models.TodoInput
	errDTO := ctx.ShouldBind(&todoUpdate)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if tc.todoService.IsAllowed(userID, id) {
		todoUpdate.UserID = userID
		result, err := tc.todoService.Update(id, todoUpdate)
		if err != nil {
			response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
			return
		}
		response := helper.BuildResponse("OK", result)
		ctx.JSON(http.StatusOK, response)
		return
	} else {
		response := helper.BuildResponse("You dont have permission", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
		return
	}
}

func (tc *TodoController) Delete(ctx *gin.Context) {
	var todo models.Todo
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(*models.User)
	userID := currentUser.ID

	todo, err = tc.todoService.FindByID(id, userID)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}
	if todo.ID != id {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	if tc.todoService.IsAllowed(userID, id) {
		todo.ID = id
		err := tc.todoService.Delete(todo)
		if err != nil {
			response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
			return
		}
		response := helper.BuildResponse("Deleted", helper.EmptyObj{})
		ctx.JSON(http.StatusOK, response)
		return
	} else {
		response := helper.BuildResponse("You dont have permission", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
		return
	}
}
