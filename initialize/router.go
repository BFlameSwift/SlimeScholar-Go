package initialize

import (
	v1 "gitee.com/online-publish/slime-scholar-go/api/v1"
	"gitee.com/online-publish/slime-scholar-go/middleware"
	"gitee.com/online-publish/slime-scholar-go/router"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r.GET("/backend/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
