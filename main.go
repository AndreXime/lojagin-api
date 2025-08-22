package main

import (
	"LojaGin/docs"
	"LojaGin/internal/config"
	"LojaGin/internal/database"
	"LojaGin/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	config.InitEnv()
	db := database.InitDB()
	router := gin.Default()

	docs.SetupDocs(router)
	routes.SetupAPI(router, db)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
