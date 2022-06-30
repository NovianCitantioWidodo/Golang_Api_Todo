package controllers

import (
	"net/http"

	"golang/helper"
	"golang/models"
	"golang/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService}
}

func (uc *UserController) Profile(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.User)
	response := helper.BuildResponse("OK", models.FilteredResponse(currentUser))
	ctx.JSON(http.StatusOK, response)
}

func (uc *UserController) Update(ctx *gin.Context) {
	var userUpdateDTO *models.UserEdit
	errDTO := ctx.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	result, err := uc.userService.Update(userUpdateDTO)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}

	response := helper.BuildResponse("OK", result)
	ctx.JSON(http.StatusOK, response)
}
