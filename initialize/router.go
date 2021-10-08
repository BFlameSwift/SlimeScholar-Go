package initialize

import (
	"slime-scholar-go/middleware"
	"slime-scholar-go/router"

	"github.com/gin-gonic/gin"
)

// 配置组路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/", v1.Index)
	Group := r.Group("api/v1/")
	{
		router.InitRouter(Group)
	}
	return r
}
