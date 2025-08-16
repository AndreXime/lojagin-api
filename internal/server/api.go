package server

import (
	"LojaGin/internal/auth"
	"LojaGin/internal/middleware"
	"LojaGin/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupAuthRoutes(api *gin.RouterGroup, db *gorm.DB) {
	userRepo := user.NewRepository(db)
	authService := auth.NewService(userRepo)
	authHandler := auth.NewHandler(authService)

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}
}

func setupUserRoutes(api *gin.RouterGroup, db *gorm.DB) {
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	userRoutes := api.Group("/users").Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.GET("/:id", userHandler.GetUserByID)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}
}

func SetupAPI(router *gin.Engine) {
	db := InitDB()
	api := router.Group("/api")

	setupAuthRoutes(api, db)
	setupUserRoutes(api, db)
}
