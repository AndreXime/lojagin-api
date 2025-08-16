package main

import (
	"LojaGin/docs"
	"LojaGin/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	docs.SetupDocs(router)
	server.SetupAPI(router)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
