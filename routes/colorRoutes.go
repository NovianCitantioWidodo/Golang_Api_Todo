package routes

import (
	"golang/controllers"
	"golang/middleware"
	"golang/services"

	"github.com/gin-gonic/gin"
)

type ColorRouteController struct {
	colorController controllers.ColorController
}

func NewRouteColorController(colorController controllers.ColorController) ColorRouteController {
	return ColorRouteController{colorController}
}

func (tc *ColorRouteController) ColorRoute(rg *gin.RouterGroup, userService services.UserService) {

	router := rg.Group("color")
	router.Use(middleware.DeserializeUser(userService))
	router.GET("/list", tc.colorController.List)
	router.GET("/detail/:id", tc.colorController.FindByID)
	router.POST("/create", tc.colorController.Insert)
	router.PUT("/edit/:id", tc.colorController.Update)
	router.DELETE("/delete/:id", tc.colorController.Delete)
}
