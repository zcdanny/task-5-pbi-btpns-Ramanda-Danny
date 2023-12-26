package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zcdanny/task-5-pbi-btpns-Ramanda-Danny/controllers"
	"github.com/zcdanny/task-5-pbi-btpns-Ramanda-Danny/database"
	"github.com/zcdanny/task-5-pbi-btpns-Ramanda-Danny/middlewares"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	userController := controllers.UserController{DB: db}
	photoController := controllers.PhotoController{DB: db}

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", userController.RegisterUser)
		userRoutes.POST("/login", userController.LoginUser)
		userRoutes.PUT("/:userId", middlewares.AuthMiddleware(), userController.UpdateUser)
		userRoutes.DELETE("/:userId", middlewares.AuthMiddleware(), userController.DeleteUser)
	}

	photoRoutes := r.Group("/photos")
	{
		photoRoutes.POST("/", middlewares.AuthMiddleware(), photoController.CreatePhoto)
		photoRoutes.GET("/", photoController.GetPhotos)
		photoRoutes.GET("/:photoId", middlewares.AuthMiddleware(), photoController.GetPhoto)
		photoRoutes.PUT("/:photoId", middlewares.AuthMiddleware(), photoController.UpdatePhoto)
		photoRoutes.DELETE("/:photoId", middlewares.AuthMiddleware(), photoController.DeletePhoto)
	}

	return r
}

func Init() *gin.Engine {
	db := database.InitDB()
	return SetupRouter(db)
}
