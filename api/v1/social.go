package v1

import (
	"net/http"
	"strconv"
	"time"

	"gitee.com/online-publish/slime-scholar-go/service"
	"gitee.com/online-publish/slime-scholar-go/model"
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
// @Failure 403 {string} string "{"success": false, "message": "用户未设置该标签"}"
// @Failure 402 {string} string "{"success": false, "message": "标签下没有文章"}"
// @Router /social/get/tag/paper [POST]
func GetTagPaper(c *gin.Context){
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	VerifyLogin(userID,authorization,c)

	tagName := c.Request.FormValue("tag_name")
	tag, notFoundTag := service.QueryATag(userID,tagName)
	if notFoundTag{
		c.JSON(403, gin.H{
			"success": false,
			"status":  403,
			"message": "用户未设置该标签",
		})
		return
	}

	papers,notFoundpaper := service.QueryTagPaper(tag.TagID)
	if notFoundpaper{
		c.JSON(402, gin.H{
			"success": false,
			"status":  402,
			"message": "标签下没有文章",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查看文献成功",
		"data":    papers,
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
	user := VerifyLogin(userID,authorization,c)
	
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
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "标签创建成功"})
}


func VerifyLogin(userID uint64,authorization string,c *gin.Context)(user model.User){
	user, notFoundUserByID := service.QueryAUserByID(userID)
	verify_answer, _ := service.VerifyAuthorization(authorization, userID, user.Username, user.Password)

	if authorization == "" || !verify_answer {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 400, "message": "用户未登录"})
		return user
	}

	if notFoundUserByID {
		c.JSON(404, gin.H{
			"success": false,
			"status":  404,
			"message": "用户ID不存在",
		})
		return user
	}
	return user
}