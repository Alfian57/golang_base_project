package router

import (
	"github.com/Alfian57/belajar-golang/internal/di"
	"github.com/Alfian57/belajar-golang/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterV1Route(router *gin.RouterGroup) {

	authHandler := di.InitializeAuthHandler()
	userHandler := di.InitializeUserHandler()

	router.POST("/login", authHandler.Login)
	router.POST("/register", authHandler.Register)
	router.POST("/refresh", middleware.AuthMiddleware(), authHandler.Refresh)
	router.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)

	admin := router.Group("admin", middleware.AuthMiddleware())

	users := admin.Group("users")
	{
		users.GET("/", userHandler.GetAllUsers)
		users.POST("/", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUserByID)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}
