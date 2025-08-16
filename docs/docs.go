package docs

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupDocs(router *gin.Engine) {
	router.StaticFile("/docs/openapi.yaml", "./docs/openapi.yaml")

	url := ginSwagger.URL("http://localhost:8080/docs/openapi.yaml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
