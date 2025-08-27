package main

import (
	"LojaGin/docs"
	"LojaGin/internal/config"
	"LojaGin/internal/database"
	"LojaGin/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Recovery())

	config.InitEnv()
	db := database.InitDB()

	docs.SetupDocs(router)
	routes.SetupAPI(router, db)

	if config.ENV_MODE != "production" {
		config.PrintRoutes(router)
	}

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
