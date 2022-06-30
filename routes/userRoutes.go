package routes

import (
	"golang/controllers"
	"golang/middleware"
	"golang/services"

	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup, userService services.UserService) {

	router := rg.Group("user")
	router.Use(middleware.DeserializeUser(userService))
	router.GET("/profile", uc.userController.Profile)
	router.PUT("/edit", uc.userController.Update)
}
