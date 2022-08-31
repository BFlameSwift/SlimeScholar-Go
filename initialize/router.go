package initialize

import (
	v1 "github.com/BFlameSwift/SlimeScholar-Go/api/v1"
	"github.com/BFlameSwift/SlimeScholar-Go/middleware"
	"github.com/BFlameSwift/SlimeScholar-Go/router"
	"github.com/gin-gonic/gin"
)

// 配置组路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(middleware.LoggerToFile())
	r.GET("/", v1.Index)
	Group := r.Group("api/v1/")
	{
		router.InitRouter(Group)
	}

	return r
}
