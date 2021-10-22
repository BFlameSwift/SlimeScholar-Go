package main

import (
	"gitee.com/online-publish/slime-scholar-go/docs"
	"gitee.com/online-publish/slime-scholar-go/initialize"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	//"io"
	//"os"
)

// @title hzh txd1 Golang Backend
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
	r := initialize.SetupRouter()
	//f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8000")
}
