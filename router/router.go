package router

import (
	v1 "gitee.com/online-publish/slime-scholar-go/api/v1"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func InitRouter(Router *gin.RouterGroup) {
	BasicRouter := Router.Group("/count")
	{
		BasicRouter.POST("/all", v1.DocumentCount)
	}

	SubmitRouter := Router.Group("/submit")
	{
		SubmitRouter.POST("/create", v1.CreateSubmit)
		SubmitRouter.POST("/check", v1.CheckSubmit)
		SubmitRouter.POST("/list", v1.ListAllSubmit)
		SubmitRouter.POST("/get/papers", v1.PaperGetAuthors)
	}
	SchoalrRouter := Router.Group("/scholar")
	{
		SchoalrRouter.POST("/info", v1.GetScholar)
		SchoalrRouter.POST("/transfer", v1.ScholarManagePaper)
	}
	UserRouter := Router.Group("/user")
	{
		UserRouter.POST("/register", v1.Register)
		UserRouter.POST("/login", v1.Login)
		UserRouter.POST("/modify", v1.ModifyUser)
		UserRouter.POST("/info", v1.TellUserInfo)
		UserRouter.POST("/confirm", v1.Confirm)
	}
	EsRouter := Router.Group("/es")
	{
		EsRouter.POST("/create/mytype", v1.CreateMyType)
		EsRouter.POST("/update/mytype", v1.UpdateMyType)
		EsRouter.POST("/get/mytype", v1.GetMyType)
		EsRouter.POST("/get/author", v1.GetAuthor)
		EsRouter.POST("/get/paper", v1.GetPaper)
		EsRouter.POST("/get/conference", v1.GetConference)
		EsRouter.POST("/get/journal", v1.GetJournal)
		EsRouter.POST("/query/paper/title", v1.TitleQueryPaper)
		EsRouter.POST("/query/paper/advanced", v1.AdvancedSearch)
		EsRouter.POST("/select/paper/title", v1.TitleSelectPaper)
		EsRouter.POST("/query/paper/author_name", v1.AuthorNameQueryPaper)
		EsRouter.POST("/query/paper/affiliation_name", v1.AffiliationNameQueryPaper)
		EsRouter.POST("/query/paper/publisher", v1.PublisherQueryPaper)
		EsRouter.POST("/query/author/name", v1.NameQueryAuthor)
		EsRouter.POST("/query/paper/doi", v1.DoiQueryPaper)
		//EsRouter.POST("/query/paper/abstract", v1.AbstractQueryPaper)
		//EsRouter.POST("/query/paper/main", v1.MainQueryPaper)
	}

	SocialRouter := Router.Group("/social")
	{
		SocialRouter.POST("/get/tags", v1.GetUserTag)
		// SocialRouter.POST("/get/tag/paper", v1.GetTagPaper)
		SocialRouter.POST("/create/tag", v1.CreateATag)
		SocialRouter.POST("/delete/tag", v1.DeleteATag)
		SocialRouter.POST("/collect/paper", v1.CollectAPaper)
		// SocialRouter.POST("/get/all/collect", v1.GetAllCollect)
		SocialRouter.POST("/get/collect/paper", v1.GetCollectPaper)
		SocialRouter.POST("/delete/collect/paper", v1.DeleteCollectPaper)
		SocialRouter.POST("/get/collect/year/paper", v1.GetCollectPaperByYear)

		SocialRouter.POST("/create/comment", v1.CreateAComment)
		SocialRouter.POST("/like/comment", v1.LikeComment)
		SocialRouter.POST("/like/cancel", v1.CancelLike)
		SocialRouter.POST("/reply/comment", v1.ReplyAComment)
		SocialRouter.POST("/get/comments", v1.GetPaperComment)
		SocialRouter.POST("/get/replies", v1.GetComReply)
	}
}
