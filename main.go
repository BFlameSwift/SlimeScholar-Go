package main


import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"slime-scholar-go/initialize"
)


// @title Slime Scholar Golang Backend
// @version 1.0
// @description Scholar Sharing platform
// @schemes https
func main() {
	//docs.SwaggerInfo.Title = "Scholar Sharing platform"
	//docs.SwaggerInfo.Description = "This is Slime Scholar Golang Backend"
	//docs.SwaggerInfo.Version = "1.0"
	//docs.SwaggerInfo.BasePath = "/api/v1"
	//docs.SwaggerInfo.Schemes = []string{"http", "https"}
	initialize.InitMySQL()
	r := initialize.SetupRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
