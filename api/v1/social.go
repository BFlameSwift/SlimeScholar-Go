package v1

import (
	"net/http"
	"strconv"
	"time"
	"encoding/json"
	// "container/list"
	"fmt"

	"gitee.com/online-publish/slime-scholar-go/service"
	"gitee.com/online-publish/slime-scholar-go/model"
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
	VerifyLogin(userID,authorization,c)
	// if err{
	// 	return
	// }
	
	tags,notFoundTags := service.QueryTagList(userID)
	if notFoundTags{
		c.JSON(403, gin.H{
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
func GetTagPaper(c *gin.Context){
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	VerifyLogin(userID,authorization,c)

	tagName := c.Request.FormValue("tag_name")
	tag,_ := service.QueryATag(userID,tagName)

	papers := service.QueryTagPaper(tag.TagID)
	if papers == nil {
		c.JSON(402, gin.H{
			"success": false,
			"status":  402,
			"message": "标签下没有文章",
		})
		return
	}

	var data []map[string]interface{}
	for _, tag_paper := range papers{
		var tmp = make(map[string]interface{})
		tmp["collect_time"] = tag_paper.CreateTime

		var map_param map[string]string = make(map[string]string)
		map_param["index"], map_param["id"] = "paper", tag_paper.PaperID
		ret, _ := service.Gets(map_param)
		body_byte, _ := json.Marshal(ret.Source)
		var paper = make(map[string]interface{})
		_ = json.Unmarshal(body_byte, &paper)

		tmp["paper"] = paper
		data = append(data,tmp)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查看文献成功",
		"data":    data,
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
func CreateATag (c *gin.Context){
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	user,err := VerifyLogin(userID,authorization,c)
	if err{
		return
	}
	
	tagName := c.Request.FormValue("tag_name")
	tag, notFoundTag := service.QueryATag(userID,tagName)
	if !notFoundTag{
		c.JSON(402, gin.H{
			"success": false,
			"status":  402,
			"message": "已创建该标签",
		})
		return
	}
	tag = model.Tag{TagName:tagName, UserID: userID, CreateTime: time.Now(), Username:user.Username}
	service.CreateATag(&tag)
	c.JSON(http.StatusOK, gin.H{"success": true,"status":  200, "message": "标签创建成功"})
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
func DeleteATag (c *gin.Context){
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	_, err := VerifyLogin(userID,authorization,c)
	if err{
		return
	}
	
	tagName := c.Request.FormValue("tag_name")
	tag, notFoundTag := service.QueryATag(userID,tagName)
	if notFoundTag{
		c.JSON(403, gin.H{"success": false,"status":  403, "message": "标签不存在"})
	}
	service.DeleteATag(tag.TagID)
	c.JSON(http.StatusOK, gin.H{"success": true,"status":  200, "message": "标签删除成功"})
}

// CollectAPaper doc
// @description 收藏文献
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param id formData string true "id"
// @Param tag_name formData string false "标签名称"
// @Success 200 {string} string "{"success": true, "message": "收藏成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Router /social/collect/paper [POST]
func CollectAPaper(c *gin.Context){
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	user, err := VerifyLogin(userID,authorization,c)
	if err{
		return
	}

	id := c.Request.FormValue("id")
	tagName := c.Request.FormValue("tag_name")
	if tagName == ""{
		tagName = "默认"
	}
	tag,notFound := service.QueryATag(userID,tagName)
	if notFound{
		tag = model.Tag{TagName:tagName, UserID: userID, CreateTime: time.Now(), Username:user.Username}
		service.CreateATag(&tag)
	}
	tagPaper := model.TagPaper{TagID:tag.TagID, TagName:tag.TagName, PaperID:id, CreateTime:time.Now()}
	service.CreateATagPaper(&tagPaper)
	c.JSON(http.StatusOK, gin.H{"success": true,"status":  200, "message": "收藏成功"})
}

// DeleteATagPaper doc
// @description 删除某标签下的文章
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param id formData string true "id"
// @Param tag_name formData string true "标签名称"
// @Success 200 {string} string "{"success": true, "message": "删除成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Router /social/delete/tag/paper [POST]
func DeleteATagPaper(c *gin.Context){
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	_,err := VerifyLogin(userID,authorization,c)
	if err{
		return
	}

	id := c.Request.FormValue("id")
	tagName := c.Request.FormValue("tag_name")
	
	tag, _ := service.QueryATag(userID,tagName)
	tagPaper, _ := service.QueryATagPaper(tag.TagID,id)
	service.DeleteATagPaper(tagPaper.ID)
	c.JSON(http.StatusOK, gin.H{"success": true,"status":  200, "message": "删除成功"})
}

// CreateAComment doc
// @description 创建评论
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param id formData string true "id"
// @Param content formData string true "评论内容"
// @Success 200 {string} string "{"success": true, "message": "评论创建成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Failure 403 {string} string "{"success": false, "message": "评论创建失败"}"
// @Router /social/create/comment [POST]
func CreateAComment (c *gin.Context){
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	_,err := VerifyLogin(userID,authorization,c)
	if err{
		return
	}

	id := c.Request.FormValue("id")
	content := c.Request.FormValue("content")
	comment := model.Comment{UserID:userID, PaperID: id, CommentTime: time.Now(), Content:content}
	notCreated := service.CreateAComment(&comment)
	if notCreated{
		c.JSON(403, gin.H{"success": false,"status":  403, "message": "评论创建失败"})
	}else{
		c.JSON(http.StatusOK, gin.H{"success": true,"status":  200, "message": "评论创建成功"})
	}
}

// LikeorUnlike doc
// @description 赞或踩评论
// @Tags 社交
// @Param comment_id formData string true "评论id"
// @Param option formData string true "赞或踩,0-赞,1-踩" 
// @Success 200 {string} string "{"success": true, "message": "操作成功"}"
// @Failure 403 {string} string "{"success": false, "message": "评论不存在"}"
// @Router /social/like/comment [POST]
func LikeorUnlike (c *gin.Context){
	commentID, _ := strconv.ParseUint(c.Request.FormValue("comment_id"), 0, 64)
	option, _ := strconv.ParseUint(c.Request.FormValue("option"), 0, 64)
	comment,notFound := service.QueryAComment(commentID)
	if notFound{
		c.JSON(403, gin.H{
			"success": false,
			"status":  403,
			"message": "评论不存在",
		})
		return
	}
	service.UpdateCommentLike(comment,option)
	c.JSON(http.StatusOK, gin.H{"success": true,"status":  200, "message": "操作成功"})
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
func ReplyAComment(c *gin.Context)  {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	user,err := VerifyLogin(userID,authorization,c)
	if err{
		return
	}

	comment_id, _ := strconv.ParseUint(c.Request.FormValue("comment_id"), 0, 64)
	comment,_ := service.QueryAComment(comment_id)

	content := c.Request.FormValue("content")
	reply := model.Comment{UserID:userID, PaperID:comment.PaperID, 
		CommentTime:time.Now(), Content:content, RelateID:comment_id}
	service.CreateAComment(&reply)
	
	utils.SendReplyEmail(user.Email)
	c.JSON(http.StatusOK, gin.H{"success": true,"status":  200, "message": "回复成功"})
}

// GetPaperComment doc
// @description 获取文献所有评论
// @Tags 社交
// @Param paper_id formData string true "文献id"
// @Success 200 {string} string "{"success": true, "message": "查找成功"}"
// @Failure 403 {string} string "{"success": false, "message": "评论不存在"}"
// @Router /social/get/comments [POST]
func GetPaperComment(c *gin.Context){

	paperID := c.Request.FormValue("paper_id")
	// fmt.Println(paperID)
	comments := service.QueryComsByPaperId(paperID)
	// fmt.Println(comments)
	if comments == nil{
		c.JSON(403, gin.H{
			"success": false,
			"status":  403,
			"message": "评论不存在",
		})
		return
	}

	var dataList []map[string]interface{}
	for _, comment := range comments{
		var com = make(map[string]interface{})
		com["id"] = comment.CommentID
		com["like"] = comment.Like
		user, _ := service.QueryAUserByID(comment.UserID)
		com["user_id"] = user.UserID
		com["username"] = user.Username
		com["content"] = comment.Content
		com["time"] = comment.CommentTime
		com["reply_count"] = comment.ReplyCount
		// fmt.Println(com)
		dataList = append(dataList,com)
	}
	fmt.Println(dataList)

	var data = make(map[string]interface{})
	data["paper_id"] = paperID

	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["id"] = "paper", paperID
	ret, _ := service.Gets(map_param)
	body_byte, _ := json.Marshal(ret.Source)
	var paper = make(map[string]interface{})
	_ = json.Unmarshal(body_byte, &paper)
	data["paper_title"] = paper["paper_title"]
	
	data["comments"] = dataList

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查找成功",
		"data": data,
	})
}

// GetComReply doc
// @description 获取某条评论的回复
// @Tags 社交
// @Param comment_id formData string true "评论id"
// @Success 200 {string} string "{"success": true, "message": "查找成功"}"
// @Failure 403 {string} string "{"success": false, "message": "回复不存在"}"
// @Router /social/get/replies [POST]
func GetComReply(c *gin.Context){
	ComID,_ := strconv.ParseUint(c.Request.FormValue("comment_id"), 0, 64)
	comment,_ := service.QueryAComment(ComID)
	replies := service.QueryComReply(ComID)
	if replies == nil{
		c.JSON(403, gin.H{
			"success": false,
			"status":  403,
			"message": "回复不存在",
		})
		return
	}

	var data = make(map[string]interface{})
	data["paper_id"] = comment.PaperID

	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["id"] = "paper", comment.PaperID
	ret, _ := service.Gets(map_param)
	body_byte, _ := json.Marshal(ret.Source)
	var paper = make(map[string]interface{})
	_ = json.Unmarshal(body_byte, &paper)
	data["paper_title"] = paper["paper_title"]

	var base_comment = make(map[string]interface{})
	base_comment["id"] = comment.CommentID
	user, _ := service.QueryAUserByID(comment.UserID)
	base_comment["username"] = user.Username
	base_comment["time"] = comment.CommentTime
	base_comment["content"] = comment.Content
	data["base_comment"] = base_comment

	var answers []map[string]interface{}
	for _, reply := range replies{
		var answer = make(map[string]interface{})
		answer["reply_id"] = reply.CommentID
		answer["time"] = reply.CommentTime
		answer["content"] = reply.Content
		answer["answerIt"] = false
		answer["myAnswer"] = " "
		reply_user, _ := service.QueryAUserByID(reply.UserID)
		answer["reply_username"] = reply_user.Username
		comment,_ = service.QueryAComment(reply.RelateID)
		replied_user,_ := service.QueryAUserByID(comment.UserID)
		answer["be_replied_username"] = replied_user.Username
		answers = append(answers,answer)
	}
	data["answers"] = answers

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查找成功",
		"data": data,
	})
}

func VerifyLogin(userID uint64,authorization string,c *gin.Context)(user model.User, err bool){
	user, notFoundUserByID := service.QueryAUserByID(userID)
	verify_answer, _ := service.VerifyAuthorization(authorization, userID, user.Username, user.Password)

	if authorization == "" || !verify_answer {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 400, "message": "用户未登录"})
		return user,true
	}

	if notFoundUserByID {
		c.JSON(404, gin.H{
			"success": false,
			"status":  404,
			"message": "用户ID不存在",
		})
		return user,true
	}
	return user,false
}
