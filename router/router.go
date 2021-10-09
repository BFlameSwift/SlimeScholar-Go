package router

import (
	v1 "gitee.com/online-publish/slime-scholar-go/api/v1"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func InitRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("/user")
	{
		UserRouter.POST("/register", v1.Register)
		UserRouter.POST("/login", v1.Login)
		UserRouter.POST("/modify", v1.ModifyUser)
		UserRouter.POST("/info", v1.TellUserInfo)
		UserRouter.POST("/confirm", v1.Confirm)
	}

}