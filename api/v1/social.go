package v1

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	_, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

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

// GetCollectPaper doc
// @description 查看用户文章列表
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param tag_name formData string false "标签名称"
// @Success 200 {string} string "{"success": true, "message": "查看文献成功", "data": "标签下的文章列表"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Failure 402 {string} string "{"success": false, "message": "标签下没有文章"}"
// @Router /social/get/collect/paper [POST]
func GetCollectPaper(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	_, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

	tagName := c.Request.FormValue("tag_name")
	papers := make([]model.TagPaper, 0)
	if tagName == "" {
		tags := service.QueryTagList(userID)
		for _, tag := range tags {
			tag_papers := service.QueryTagPaper(tag.TagID)
			for _, tmp := range tag_papers {
				papers = append(papers, tmp)
			}
		}
	} else {
		tag, _ := service.QueryATag(userID, tagName)
		papers = service.QueryTagPaper(tag.TagID)
	}

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
	for _, paper := range papers {
		paper_ids = append(paper_ids, paper.PaperID)
	}
	paper_detail := service.GetPapers(paper_ids)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查看文献成功",
		"data":    paper_detail,
	})
}

// GetCollectPaperByYear doc
// @description 根据年份筛选收藏列表
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param tag_name formData string false "标签名称"
// @Param min_year formData string true "年份下限"
// @Param max_year formData string true "年份上限"
// @Success 200 {string} string "{"success": true, "message": "查看文献成功", "data": "标签下的文章列表"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Failure 402 {string} string "{"success": false, "message": "没有收藏文章"}"
// @Router /social/get/collect/year/paper [POST]
func GetCollectPaperByYear(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	_, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

	minYear, _ := strconv.ParseUint(c.Request.FormValue("min_year"), 0, 64)
	maxYear, _ := strconv.ParseUint(c.Request.FormValue("max_year"), 0, 64)
	tagName := c.Request.FormValue("tag_name")

	papers := make([]model.TagPaper, 0)
	if tagName == "" {
		tags := service.QueryTagList(userID)
		for _, tag := range tags {
			tag_papers := service.QueryTagPaper(tag.TagID)
			for _, tmp := range tag_papers {
				papers = append(papers, tmp)
			}
		}
	} else {
		tag, _ := service.QueryATag(userID, tagName)
		papers = service.QueryTagPaper(tag.TagID)
	}

	if papers == nil || len(papers) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  402,
			"message": "没有收藏文章",
		})
		return
	}

	var paper_ids []string
	for _, paper := range papers {
		paper_ids = append(paper_ids, paper.PaperID)
	}
	var data map[string]interface{}
	data = service.IdsGetItems(paper_ids, "paper")

	var paper_detail []map[string]interface{}

	k := 0
	for _, tmp := range data {
		year, _ := strconv.ParseUint(tmp.(map[string]interface{})["year"].(string), 0, 64)
		if year >= minYear && year <= maxYear {
			tmp.(map[string]interface{})["create_time"] = papers[k].CreateTime
			tmp = service.ComplePaper(tmp.(map[string]interface{}))
			paper_detail = append(paper_detail, tmp.(map[string]interface{}))
		}
		k++
	}
	if len(paper_detail) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  402,
			"message": "没有收藏文章",
		})
		return
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
	user, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

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
	_, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

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
	user, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

	id := c.Request.FormValue("paper_id")
	tagName := c.Request.FormValue("tag_name")
	if tagName == "" || len(tagName) == 0 {
		tagName = "默认"
	}

	tags := strings.Split(tagName, `-<^_^>-`)
	fmt.Println(tags)
	tmp := len(tags)
	fmt.Println(tmp)
	for _, tags_name := range tags {
		if tags_name != "" && len(tags_name) != 0 {
			tag, notFound := service.QueryATag(userID, tags_name)
			if notFound {
				tag = model.Tag{TagName: tags_name, UserID: userID, CreateTime: time.Now(), Username: user.Username}
				service.CreateATag(&tag)
			}
			_, notFoundPaper := service.QueryATagPaper(tag.TagID, id)
			if notFoundPaper {
				tagPaper := model.TagPaper{TagID: tag.TagID, TagName: tag.TagName, PaperID: id, CreateTime: time.Now()}
				service.CreateATagPaper(&tagPaper)
			} else {
				tmp--
			}
		} else {
			tmp--
		}
	}
	fmt.Println(tmp)
	if tmp == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 403, "message": "文献已收藏"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "收藏成功"})
}

// DeleteCollectPaper doc
// @description 删除收藏的文章
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param paper_id formData string true "文献id"
// @Param tag_name formData string false "标签名称"
// @Success 200 {string} string "{"success": true, "message": "删除成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Router /social/delete/collect/paper [POST]
func DeleteCollectPaper(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	_, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

	id := c.Request.FormValue("paper_id")
	tagName := c.Request.FormValue("tag_name")

	if tagName != "" {
		tag, _ := service.QueryATag(userID, tagName)
		tagPaper, _ := service.QueryATagPaper(tag.TagID, id)
		service.DeleteATagPaper(tagPaper.ID)
	} else {
		tags := service.QueryTagList(userID)
		for _, tag := range tags {
			paper, notfound := service.QueryATagPaper(tag.TagID, id)
			if !notfound {
				service.DeleteATagPaper(paper.ID)
			}
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
	user, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

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
	user, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

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
	user, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

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
	user, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

	comment_id, _ := strconv.ParseUint(c.Request.FormValue("comment_id"), 0, 64)
	comment, _ := service.QueryAComment(comment_id)

	content := c.Request.FormValue("content")
	reply := model.Comment{UserID: user.UserID, Username: user.Username,
		PaperID: comment.PaperID, PaperTitle: comment.PaperTitle,
		CommentTime: time.Now(), Content: content, RelateID: comment_id}
	service.CreateAComment(&reply)

	// var map_param map[string]string = make(map[string]string)
	// map_param["index"], map_param["id"] = "paper", comment.PaperID
	// ret, _ := service.Gets(map_param)
	// body_byte, _ := json.Marshal(ret.Source)
	// var paper = make(map[string]interface{})
	// _ = json.Unmarshal(body_byte, &paper)
	// paper_url := "https://dx.doi.org/" + paper["doi"].(string)

	be_reply_user, _ := service.QueryAUserByID(comment.UserID)
	utils.SendReplyEmail(be_reply_user.Email, comment.PaperID)

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
// @Param user_id formData string false "用户ID"
// @Param paper_id formData string true "文献id"
// @Success 200 {string} string "{"success": true, "message": "查找成功"}"
// @Failure 403 {string} string "{"success": false, "message": "评论不存在"}"
// @Router /social/get/comments [POST]
func GetPaperComment(c *gin.Context) {

	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	err := false
	if userID != 0 {
		authorization := c.Request.Header.Get("Authorization")
		user, notFoundUserByID := service.QueryAUserByID(userID)
		verify_answer, _ := service.VerifyAuthorization(authorization, userID, user.Username, user.Password)

		if authorization == "" || !verify_answer || notFoundUserByID {
			err = true
		}
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
		if !err && service.UserLike(userID, comment.CommentID) {
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

// FollowUser doc
// @description 关注用户
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData int true "用户ID"
// @Param be_user_id formData int true "被关注用户ID"
// @Success 200 {string} string "{"success": true, "message": "标签创建成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Router /social/follow/user [POST]
func FollowUser(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	beFollower, _ := strconv.ParseUint(c.Request.FormValue("be_user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	_, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}
	follower, notFound := service.QueryAUserByID(beFollower)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": true, "status": 404, "message": "被关注用户不存在"})
		return
	}
	service.FollowUser(userID, follower.UserID)
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "标签创建成功"})
}

// UnFollowUser doc
// @description 取消关注用户
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData int true "用户ID"
// @Param be_user_id formData int true "被取消关注用户ID"
// @Success 200 {string} string "{"success": true, "message": "标签创建成功"}"
// @Failure 402 {string} string "{"success": false, "message": "用户未关注该被关注用户"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Router /social/unfollow/user [POST]
func UnFollowUser(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	beFollower, _ := strconv.ParseUint(c.Request.FormValue("be_user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	_, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}
	_, notFound := service.QueryAUserByID(beFollower)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": true, "status": 404, "message": "被关注用户不存在"})
		return
	}
	success := service.CanCelFollowUser(userID, beFollower)
	if !success {
		c.JSON(http.StatusOK, gin.H{"success": true, "status": 402, "message": "用户未关注该被关注用户"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "标签创建成功"})
}

// GetUserFollowingList doc
// @description 新建
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData int true "用户ID"
// @Success 200 {string} string "{"success": true, "message": "查找成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Router /social/get/user/following [POST]
func GetUserFollowingList(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	_, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

	followingList := service.GetUserFollowingList(userID)
	userList := make([]model.User, 0)
	for _, id := range followingList {
		thisUser, _ := service.QueryAUserByID(id)
		userList = append(userList, thisUser)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "查找成功", "user_list": userList})
}

// GetUserFollowedList doc
// @description 获取用户的粉丝列表
// @Tags 社交
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData int true "用户ID"
// @Success 200 {string} string "{"success": true, "message": "查找成功"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Router /social/get/user/followed [POST]
func GetUserFollowedList(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	_, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

	followingList := service.GetUserFollowedList(userID)
	userList := make([]model.User, 0)
	for _, id := range followingList {
		thisUser, _ := service.QueryAUserByID(id)
		userList = append(userList, thisUser)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "message": "查找成功", "user_list": userList})
}
