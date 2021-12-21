package main

import (
	"gitee.com/online-publish/slime-scholar-go/docs"
	"gitee.com/online-publish/slime-scholar-go/initialize"
)

// @title Slime Scholar Golang Backend
// @version 1.0
// @description hzh company
// @schemes https
func main() {
	docs.SwaggerInfo.Title = "Slime scholar !"
	docs.SwaggerInfo.Description = "This is a Scholar sharing platform"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	initialize.Init()
	r := initialize.SetupRouter()
	r.Run(":8000")
}
