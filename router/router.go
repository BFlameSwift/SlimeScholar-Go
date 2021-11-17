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
		UserRouter.POST("/tag", v1.GetUserTag)
	}
	EsRouter := Router.Group("/es")
	{
		EsRouter.POST("/create/mytype", v1.CreateMyType)
		EsRouter.POST("/update/mytype", v1.UpdateMyType)
		EsRouter.POST("/get/mytype", v1.GetMyType)
		EsRouter.POST("/get/author", v1.GetAuthor)
		EsRouter.POST("/get/paper", v1.GetPaper)
		EsRouter.POST("/query/paper/title", v1.TitleQueryPaper)
		EsRouter.POST("/query/author/name", v1.NameQueryAuthor)
	}
}
