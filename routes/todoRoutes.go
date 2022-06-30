package routes

import (
	"golang/controllers"
	"golang/middleware"
	"golang/services"

	"github.com/gin-gonic/gin"
)

type TodoRouteController struct {
	todoController controllers.TodoController
}

func NewRouteTodoController(todoController controllers.TodoController) TodoRouteController {
	return TodoRouteController{todoController}
}

func (tc *TodoRouteController) TodoRoute(rg *gin.RouterGroup, userService services.UserService) {

	router := rg.Group("todo")
	router.Use(middleware.DeserializeUser(userService))
	router.GET("/list", tc.todoController.List)
	router.GET("/detail/:id", tc.todoController.FindByID)
	router.POST("/create", tc.todoController.Insert)
	router.PUT("/edit/:id", tc.todoController.Update)
	router.DELETE("/delete/:id", tc.todoController.Delete)
}
