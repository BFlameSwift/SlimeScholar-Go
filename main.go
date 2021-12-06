package main

import (
	"gitee.com/online-publish/slime-scholar-go/docs"
	"gitee.com/online-publish/slime-scholar-go/initialize"
	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/gin-gonic/gin"
	"io"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Slime Scholar Golang Backend
// @version 1.0
// @description hzh company
// @schemes https
func main() {
	docs.SwaggerInfo.Title = "Slime scholar"
	docs.SwaggerInfo.Description = "This is a Scholar sharing platform"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	initialize.InitMySQL()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r := initialize.SetupRouter()
	service.Init()
	r.GET("/backend/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8000")
}
