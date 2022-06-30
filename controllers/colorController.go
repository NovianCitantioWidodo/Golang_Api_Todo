package controllers

import (
	"net/http"
	"strconv"

	"golang/helper"
	"golang/models"
	"golang/services"

	"github.com/gin-gonic/gin"
)

type ColorController struct {
	colorService services.ColorService
}

func NewColorController(colorService services.ColorService) ColorController {
	return ColorController{colorService}
}

func (cc *ColorController) List(ctx *gin.Context) {
	var colors []*models.Color
	colors, err := cc.colorService.All()
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}
	response := helper.BuildResponse("OK", colors)
	ctx.JSON(http.StatusOK, response)
	return
}

func (cc *ColorController) FindByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	var color models.Color
	color, err = cc.colorService.FindByID(id)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}

	if color.ID != id {
		response := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.BuildResponse("OK", color)
	ctx.JSON(http.StatusOK, response)
}

func (cc *ColorController) Insert(ctx *gin.Context) {
	var colorCreate models.ColorInput
	errDTO := ctx.ShouldBind(&colorCreate)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	} else {
		result, err := cc.colorService.Insert(colorCreate)
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

func (cc *ColorController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	var color models.Color
	color, err = cc.colorService.FindByID(id)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}
	if color.ID != id {
		response := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	var colorUpdate models.ColorInput
	errDTO := ctx.ShouldBind(&colorUpdate)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := cc.colorService.Update(id, colorUpdate)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}
	response := helper.BuildResponse("OK", result)
	ctx.JSON(http.StatusOK, response)
	return
}

func (cc *ColorController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var color models.Color
	color, err = cc.colorService.FindByID(id)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}
	if color.ID != id {
		response := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	color.ID = id
	err = cc.colorService.Delete(color)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}

	response := helper.BuildResponse("Deleted", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}
