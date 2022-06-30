package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang/config"
	"golang/controllers"
	"golang/routes"
	"golang/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	server         *gin.Engine
	ctx            context.Context
	postgresclient *sql.DB

	userService         services.UserService
	userController      controllers.UserController
	userRouteController routes.UserRouteController

	authService         services.AuthService
	authController      controllers.AuthController
	authRouteController routes.AuthRouteController

	todoService         services.TodoService
	todoController      controllers.TodoController
	todoRouteController routes.TodoRouteController

	colorService         services.ColorService
	colorController      controllers.ColorController
	colorRouteController routes.ColorRouteController
)

func init() {
	err := godotenv.Load("app.env")
	if err != nil {
		fmt.Println("Failed to load env file")
	}

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
	}

	db, err := gorm.Open(postgres.Open(config.DBUri), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	postgresclient, err = db.DB()
	if err != nil {
		log.Fatal("Failed to connect mysql")
	}

	// db.AutoMigrate(&models.DBResponse{})
	userService = services.NewUserService(db)
	userController = controllers.NewUserController(userService)
	userRouteController = routes.NewRouteUserController(userController)

	authService = services.NewAuthService(db)
	authController = controllers.NewAuthController(authService, userService, ctx, db)
	authRouteController = routes.NewAuthRouteController(authController)

	todoService = services.NewTodoService(db)
	todoController = controllers.NewTodoController(todoService)
	todoRouteController = routes.NewRouteTodoController(todoController)

	colorService = services.NewColorService(db)
	colorController = controllers.NewColorController(colorService)
	colorRouteController = routes.NewRouteColorController(colorController)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
	}

	defer postgresclient.Close()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))
	server.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "welcome"})
	})

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "status ok"})
	})

	authRouteController.AuthRoute(router, userService)
	userRouteController.UserRoute(router, userService)
	todoRouteController.TodoRoute(router, userService)
	colorRouteController.ColorRoute(router, userService)

	log.Fatal(server.Run(":" + config.Port))
}
