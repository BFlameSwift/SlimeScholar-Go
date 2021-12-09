package v1

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"
	// "container/list"
	"fmt"

	"gitee.com/online-publish/slime-scholar-go/model"
	"gitee.com/online-publish/slime-scholar-go/service"
	"gitee.com/online-publish/slime-scholar-go/utils"
	"github.com/gin-gonic/gin"
)

// GetUserTag doc
// @description 查看用户所有标签
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Success 200 {string} string "{"success": true, "message": "查看用户标签成功", "data": "model.User的所有标签"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 403 {string} string "{"success": false, "message": "未查询到结果"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Router /social/get/tags [POST]
func GetUserTag(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	VerifyLogin(userID, authorization, c)

	tags := service.QueryTagList(userID)
	if tags == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  403,
			"message": "未查询到结果",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查看用户标签成功",
		"data":    tags,
	})
}

// GetTagPaper doc
// @description 查看用户标签的文章列表
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param tag_name formData string true "标签名称"
// @Success 200 {string} string "{"success": true, "message": "查看文献成功", "data": "标签下的文章列表"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Failure 402 {string} string "{"success": false, "message": "标签下没有文章"}"
// @Router /social/get/tag/paper [POST]
func GetTagPaper(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	VerifyLogin(userID, authorization, c)
	tagName := c.Request.FormValue("tag_name")
	tag, _ := service.QueryATag(userID, tagName)

	papers := service.QueryTagPaper(tag.TagID)
	fmt.Println(papers)
	if papers == nil || len(papers) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  402,
			"message": "标签下没有文章",
		})
		return
	}

	var paper_ids []string
	for _,paper := range papers{
		paper_ids = append(paper_ids,paper.PaperID)
	}
	var data map[string]interface{}
	data = service.IdsGetItems(paper_ids,"paper")

	var paper_detail []map[string]interface{}

	k := 0
	for _,tmp := range data{
		tmp.(map[string]interface{})["create_time"] = papers[k].CreateTime
		k++
		paper_detail = append(paper_detail,tmp.(map[string]interface{}))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查看文献成功",
		"data":    paper_detail,
	})

}

// GetAllCollect doc
// @description 获取用户收藏的所有文献
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Success 200 {string} string "{"success": true, "message": "查看文献成功", "data": "文章列表"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Failure 402 {string} string "{"success": false, "message": "用户无收藏文章"}"
// @Router /social/get/all/collect [POST]
func GetAllCollect(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	VerifyLogin(userID, authorization, c)

	papers := service.QueryAllPaper()
	fmt.Println(papers)
	if papers == nil || len(papers) == 0{
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  402,
			"message": "用户无收藏文章",
		})
		return
	}

	var paper_ids []string
	for _,paper := range papers{
		paper_ids = append(paper_ids,paper.PaperID)
	}
	var data map[string]interface{}
	data = service.IdsGetItems(paper_ids,"paper")
	var paper_detail []map[string]interface{}

	k := 0
	for _,tmp := range data{
		tmp.(map[string]interface{})["create_time"] = papers[k].CreateTime
		k++
		paper_detail = append(paper_detail,tmp.(map[string]interface{}))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查看文献成功",
		"data":    paper_detail,
	})

}

// CreateATag doc
// @description 新建标签
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param tag_name formData string true "标签名称"
// @Success 200 {string} string "{"success": true, "message": "标签创建成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Failure 402 {string} string "{"success": false, "message": "已创建该标签"}"
// @Router /social/create/tag [POST]
func CreateATag(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	user,_ := VerifyLogin(userID, authorization, c)

	tagName := c.Request.FormValue("tag_name")
	tag, notFoundTag := service.QueryATag(userID, tagName)
	if !notFoundTag {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  402,
			"message": "已创建该标签",
		})
		return
	}
	tag = model.Tag{TagName: tagName, UserID: userID, CreateTime: time.Now(), Username: user.Username}
	service.CreateATag(&tag)
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "标签创建成功"})
}

// DeleteATag doc
// @description 删除标签
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param tag_name formData string true "标签名称"
// @Success 200 {string} string "{"success": true, "message": "标签删除成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Failure 403 {string} string "{"success": false, "message": "标签不存在"}"
// @Router /social/delete/tag [POST]
func DeleteATag(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	VerifyLogin(userID, authorization, c)

	tagName := c.Request.FormValue("tag_name")
	tag, notFoundTag := service.QueryATag(userID, tagName)
	if notFoundTag {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 403, "message": "标签不存在"})
	}
	service.DeleteATag(tag.TagID)
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "标签删除成功"})
}

// CollectAPaper doc
// @description 收藏文献
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param paper_id formData string true "文献id"
// @Param tag_name formData string false "标签名称"
// @Success 200 {string} string "{"success": true, "message": "收藏成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Failure 403 {string} string "{"success": false, "message": "文献已收藏"}"
// @Router /social/collect/paper [POST]
func CollectAPaper(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	user,_ := VerifyLogin(userID, authorization, c)

	id := c.Request.FormValue("paper_id")
	tagName := c.Request.FormValue("tag_name")
	if tagName == "" {
		tagName = "默认"
	}
	tag, notFound := service.QueryATag(userID, tagName)
	if notFound {
		tag = model.Tag{TagName: tagName, UserID: userID, CreateTime: time.Now(), Username: user.Username}
		service.CreateATag(&tag)
	}
	_,notFoundPaper := service.QueryATagPaper(tag.TagID, id)
	if !notFoundPaper {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 403, "message": "文献已收藏"})
	}
	tagPaper := model.TagPaper{TagID: tag.TagID, TagName: tag.TagName, PaperID: id, CreateTime: time.Now()}
	service.CreateATagPaper(&tagPaper)
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "收藏成功"})
}

// DeleteATagPaper doc
// @description 删除某标签下的文章
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param paper_id formData string true "文献id"
// @Param tag_name formData string true "标签名称"
// @Success 200 {string} string "{"success": true, "message": "删除成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Router /social/delete/tag/paper [POST]
func DeleteATagPaper(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	VerifyLogin(userID, authorization, c)

	id := c.Request.FormValue("paper_id")
	tagName := c.Request.FormValue("tag_name")

	tag, _ := service.QueryATag(userID, tagName)
	tagPaper, _ := service.QueryATagPaper(tag.TagID, id)
	service.DeleteATagPaper(tagPaper.ID)
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "删除成功"})
}

// DeleteCollectPaper doc
// @description 取消文章收藏
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param paper_id formData string true "文献id"
// @Success 200 {string} string "{"success": true, "message": "删除成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Router /social/delete/collect/paper [POST]
func DeleteCollectPaper(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	VerifyLogin(userID, authorization, c)

	id := c.Request.FormValue("paper_id")

	tags := service.QueryTagList(userID)
	for _,tag := range tags{
		paper,notfound := service.QueryATagPaper(tag.TagID,id)
		if !notfound{
			service.DeleteATagPaper(paper.ID)
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "删除成功"})
}

// CreateAComment doc
// @description 创建评论
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param paper_id formData string true "文献id"
// @Param content formData string true "评论内容"
// @Success 200 {string} string "{"success": true, "message": "评论创建成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Failure 403 {string} string "{"success": false, "message": "评论创建失败"}"
// @Router /social/create/comment [POST]
func CreateAComment(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	user, _ := VerifyLogin(userID, authorization, c)

	id := c.Request.FormValue("paper_id")
	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["id"] = "paper", id
	ret, _ := service.Gets(map_param)
	body_byte, _ := json.Marshal(ret.Source)
	var paper = make(map[string]interface{})
	_ = json.Unmarshal(body_byte, &paper)

	content := c.Request.FormValue("content")

	comment := model.Comment{UserID: user.UserID, Username: user.Username,
		PaperID: paper["paper_id"].(string), PaperTitle: paper["paper_title"].(string),
		CommentTime: time.Now(), Content: content}

	notCreated := service.CreateAComment(&comment)
	if notCreated {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 403, "message": "评论创建失败"})
	} else {
		paperID := comment.PaperID
		// fmt.Println(paperID)
		comments := service.QueryComsByPaperId(paperID)
		// fmt.Println(comments)

		var dataList []map[string]interface{}
		for _, comment := range comments {
			var com = make(map[string]interface{})
			com["id"] = comment.CommentID
			com["like"] = comment.Like
			com["is_animating"] = false
			com["is_like"] = false
			if service.UserLike(userID, comment.CommentID) {
				com["is_like"] = true
			}
			com["user_id"] = comment.UserID
			com["username"] = comment.Username
			com["content"] = comment.Content
			com["time"] = comment.CommentTime
			com["reply_count"] = comment.ReplyCount
			// fmt.Println(com)
			dataList = append(dataList, com)
		}
		// fmt.Println(dataList)

		var data = make(map[string]interface{})
		data["paper_id"] = paperID

		data["paper_title"] = comments[0].PaperTitle

		data["comments"] = dataList
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"status":  200,
			"message": "评论创建成功",
			"data":    data,
		})
	}
}

// LikeComment doc
// @description 点赞评论
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param comment_id formData string true "评论id"
// @Success 200 {string} string "{"success": true, "message": "操作成功"}"
// @Failure 403 {string} string "{"success": false, "message": "评论不存在"}"
// @Router /social/like/comment [POST]
func LikeComment(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	user, _ := VerifyLogin(userID, authorization, c)

	commentID, _ := strconv.ParseUint(c.Request.FormValue("comment_id"), 0, 64)
	comment, notFound := service.QueryAComment(commentID)
	if notFound {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  403,
			"message": "评论不存在",
		})
		return
	}
	service.UpdateCommentLike(comment, user)
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "操作成功"})
}

// CancelLike doc
// @description 取消点赞
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param comment_id formData string true "评论id"
// @Success 200 {string} string "{"success": true, "message": "操作成功"}"
// @Failure 403 {string} string "{"success": false, "message": "用户未点赞"}"
// @Router /social/like/cancel [POST]
func CancelLike(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	user, _ := VerifyLogin(userID, authorization, c)

	commentID, _ := strconv.ParseUint(c.Request.FormValue("comment_id"), 0, 64)
	comment, _ := service.QueryAComment(commentID)

	notFound := service.CancelLike(comment, user)
	if notFound {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  403,
			"message": "用户未点赞",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "操作成功"})
}

// ReplyAComment doc
// @description 回复评论
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param comment_id formData string true "评论id"
// @Param content formData string true "回复内容"
// @Success 200 {string} string "{"success": true, "message": "回复成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Router /social/reply/comment [POST]
func ReplyAComment(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	user, _ := VerifyLogin(userID, authorization, c)

	comment_id, _ := strconv.ParseUint(c.Request.FormValue("comment_id"), 0, 64)
	comment, _ := service.QueryAComment(comment_id)

	content := c.Request.FormValue("content")
	reply := model.Comment{UserID: user.UserID, Username: user.Username,
		PaperID: comment.PaperID, PaperTitle: comment.PaperTitle,
		CommentTime: time.Now(), Content: content, RelateID: comment_id}
	service.CreateAComment(&reply)

	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["id"] = "paper", comment.PaperID
	ret, _ := service.Gets(map_param)
	body_byte, _ := json.Marshal(ret.Source)
	var paper = make(map[string]interface{})
	_ = json.Unmarshal(body_byte, &paper)
	paper_url := "https://dx.doi.org/" + paper["doi"].(string)

	be_reply_user, _ := service.QueryAUserByID(comment.UserID)
	utils.SendReplyEmail(be_reply_user.Email, paper_url)

	bas_comment := service.QueryABaseCom(comment)
	replies := service.QueryComReply(bas_comment.CommentID)

	var data = make(map[string]interface{})
	data["paper_id"] = bas_comment.PaperID
	data["paper_title"] = bas_comment.PaperTitle

	var base_comment = make(map[string]interface{})
	base_comment["id"] = bas_comment.CommentID
	base_comment["username"] = bas_comment.Username
	base_comment["time"] = bas_comment.CommentTime
	base_comment["content"] = bas_comment.Content
	data["base_comment"] = base_comment

	var answers []map[string]interface{}
	for _, reply := range replies {
		var answer = make(map[string]interface{})
		answer["reply_id"] = reply.CommentID
		answer["time"] = reply.CommentTime
		answer["content"] = reply.Content
		answer["answerIt"] = false
		answer["myAnswer"] = " "
		answer["reply_username"] = reply.Username
		comment, _ = service.QueryAComment(reply.RelateID)
		answer["be_replied_username"] = comment.Username
		answers = append(answers, answer)
	}
	sort.Slice(answers, func(i, j int) bool {
		return (answers[i]["time"].(time.Time)).Before(answers[j]["time"].(time.Time)) //顺序
	})
	data["answers"] = answers
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "回复成功",
		"data":    data,
	})
}

// GetPaperComment doc
// @description 获取文献所有评论，时间倒序
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param paper_id formData string true "文献id"
// @Success 200 {string} string "{"success": true, "message": "查找成功"}"
// @Failure 403 {string} string "{"success": false, "message": "评论不存在"}"
// @Router /social/get/comments [POST]
func GetPaperComment(c *gin.Context) {

	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	if userID != 0 {
		authorization := c.Request.Header.Get("Authorization")
		VerifyLogin(userID, authorization, c)
	}

	paperID := c.Request.FormValue("paper_id")
	// fmt.Println(paperID)
	comments := service.QueryComsByPaperId(paperID)
	fmt.Println(comments)
	if len(comments) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  403,
			"message": "评论不存在",
		})
		return
	}

	var dataList []map[string]interface{}
	for _, comment := range comments {
		var com = make(map[string]interface{})
		com["id"] = comment.CommentID
		com["like"] = comment.Like
		com["is_animating"] = false
		com["is_like"] = false
		if service.UserLike(userID, comment.CommentID) {
			com["is_like"] = true
		}
		com["user_id"] = comment.UserID
		com["username"] = comment.Username
		com["content"] = comment.Content
		com["time"] = comment.CommentTime
		com["reply_count"] = comment.ReplyCount
		// fmt.Println(com)
		dataList = append(dataList, com)
	}
	// fmt.Println(dataList)

	var data = make(map[string]interface{})
	data["paper_id"] = paperID

	data["paper_title"] = comments[0].PaperTitle

	data["comments"] = dataList

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查找成功",
		"data":    data,
	})
}

// GetComReply doc
// @description 获取某条评论的回复
// @Tags 社交
// @Param comment_id formData string true "评论id"
// @Success 200 {string} string "{"success": true, "message": "查找成功"}"
// @Failure 403 {string} string "{"success": false, "message": "回复不存在"}"
// @Router /social/get/replies [POST]
func GetComReply(c *gin.Context) {
	ComID, _ := strconv.ParseUint(c.Request.FormValue("comment_id"), 0, 64)
	comment, _ := service.QueryAComment(ComID)

	var data = make(map[string]interface{})
	data["paper_id"] = comment.PaperID

	data["paper_title"] = comment.PaperTitle

	var base_comment = make(map[string]interface{})
	base_comment["id"] = comment.CommentID
	base_comment["username"] = comment.Username
	base_comment["time"] = comment.CommentTime
	base_comment["content"] = comment.Content
	data["base_comment"] = base_comment

	replies := service.QueryComReply(ComID)
	if len(replies) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  403,
			"message": "回复不存在",
			"data":    data,
		})
		return
	}

	fmt.Println(replies)

	var answers []map[string]interface{}
	for _, reply := range replies {
		var answer = make(map[string]interface{})
		answer["reply_id"] = reply.CommentID
		answer["time"] = reply.CommentTime
		answer["content"] = reply.Content
		answer["answerIt"] = false
		answer["myAnswer"] = " "
		answer["reply_username"] = reply.Username
		comment, _ = service.QueryAComment(reply.RelateID)
		answer["be_replied_username"] = comment.Username
		answers = append(answers, answer)
	}
	// answers = MapSort(answers,"time")
	sort.Slice(answers, func(i, j int) bool {
		// if answer[i]["no"] == answer[j]["no"] {
		//     return s1[i]["score"] < s1[j]["score"]
		// }
		// return answers[i]["reply_id"].(uint64) < answers[j]["reply_id"].(uint64)
		return (answers[i]["time"].(time.Time)).Before(answers[j]["time"].(time.Time)) //顺序
	})
	data["answers"] = answers

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查找成功",
		"data":    data,
	})
}

func VerifyLogin(userID uint64, authorization string, c *gin.Context) (user model.User, err bool) {
	user, notFoundUserByID := service.QueryAUserByID(userID)
	verify_answer, _ := service.VerifyAuthorization(authorization, userID, user.Username, user.Password)

	if authorization == "" || !verify_answer {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 400, "message": "用户未登录"})
		return user, true
	}

	if notFoundUserByID {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  404,
			"message": "用户ID不存在",
		})
		return user, true
	}
	return user, false
}

type MapsSort struct {
	Key     string
	MapList []map[string]interface{}
}

// Len 为集合内元素的总数
func (m *MapsSort) Len() int {
	return len(m.MapList)
}

//如果index为i的元素小于index为j的元素，则返回true，否则返回false
func (m *MapsSort) Less(i, j int) bool {
	var ivalue float64
	var jvalue float64
	var err error
	fmt.Println(m.Key)
	switch m.MapList[i][m.Key].(type) {
	case string:
		ivalue, err = strconv.ParseFloat(m.MapList[i][m.Key].(string), 64)
		fmt.Println(ivalue)
		if err != nil {
			//   logger.Error("map数组排序string转float失败：%v",err)
			return true
		}
		// case int:
		//    ivalue = float64(m.MapList[i][m.Key].(int))
		// case float64:
		//    ivalue = m.MapList[i][m.Key].(float64)
		// case int64:
		//    ivalue = float64(m.MapList[i][m.Key].(int64))
	}
	switch m.MapList[j][m.Key].(type) {
	case string:
		jvalue, err = strconv.ParseFloat(m.MapList[j][m.Key].(string), 64)
		fmt.Println(jvalue)
		if err != nil {
			//   logger.Error("map数组排序string转float失败：%v",err)
			return true
		}
		// case int:
		//    jvalue = float64(m.MapList[j][m.Key].(int))
		// case float64:
		//    jvalue = m.MapList[j][m.Key].(float64)
		// case int64:
		//    jvalue = float64(m.MapList[j][m.Key].(int64))
	}
	return ivalue > jvalue
}

//Swap 交换索引为 i 和 j 的元素
func (m *MapsSort) Swap(i, j int) {
	m.MapList[i], m.MapList[j] = m.MapList[j], m.MapList[i]
}

func MapSort(ms []map[string]interface{}, key string) []map[string]interface{} {
	mapsSort := MapsSort{}
	mapsSort.Key = key
	mapsSort.MapList = ms
	sort.Sort(&mapsSort)
	return mapsSort.MapList
}
