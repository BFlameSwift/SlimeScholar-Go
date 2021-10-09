package v1

import (
	"net/http"
	"strconv"

	"gitee.com/online-publish/slime-scholar-go/model"
	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/gin-gonic/gin"
)

// Index doc
// @description 测试 Index 页
// @Tags 测试
// @Success 200 {string} string "{"success": true, "message": "gcp"}"
// @Router / [GET]
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "gcp"})
}

// Register doc
// @description 注册
// @Tags 用户管理
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Param user_info formData string true "用户个人信息"
// @Param user_type formData string true "用户类型（0: 普通用户，1: 认证机构用户）"
// @Param affiliation formData string false "认证机构名"
// @Success 200 {string} string "{"success": true, "message": "用户创建成功"}"
// @Failure 200 {string} string "{"success": false, "message": "用户已存在"}"
// @Router /user/register [POST]
func Register(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	userInfo := c.Request.FormValue("user_info")
	userType, _ := strconv.ParseUint(c.Request.FormValue("user_type"), 0, 64)
	affiliation := c.Request.FormValue("affiliation")
	user := model.User{Username: username, Password: password, UserInfo: userInfo, UserType: userType, Affiliation: affiliation}
	_, notFound := service.QueryAUserByUsername(username)
	if notFound {
		service.CreateAUser(&user)
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "用户创建成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户已存在"})
	}
}

// Login doc
// @description 登录
// @Tags 用户管理
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {string} string "{"success": true, "message": "登录成功", "detail": user的信息}"
// @Failure 200 {string} string "{"success": false, "message": "密码错误"}"
// @Failure 200 {string} string "{"success": false, "message": "没有该用户"}"
// @Router /user/login [POST]
func Login(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	user, notFound := service.QueryAUserByUsername(username)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有该用户"})
	} else {
		if user.Password != password {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "密码错误"})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "登录成功", "detail": user})
		}
	}
}

// ModifyUser doc
// @description 修改用户信息（支持修改用户名和密码）
// @Tags 用户管理
// @Param user_id formData string true "用户ID"
// @Param username formData string true "用户名"
// @Param user_info formData string true "用户个人信息"
// @Param password_old formData string true "原密码"
// @Param password_new formData string true "新密码"
// @Success 200 {string} string "{"success": true, "message": "修改成功", "data": "model.User的所有信息"}"
// @Failure 200 {string} string "{"success": false, "message": "原密码输入错误"}"
// @Failure 200 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 400 {string} string "{"success": false, "message": "数据库操作时的其他错误"}"
// @Router /user/modify [POST]
func ModifyUser(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	username := c.Request.FormValue("username")
	userInfo := c.Request.FormValue("user_info")
	passwordOld := c.Request.FormValue("password_old")
	passwordNew := c.Request.FormValue("password_new")
	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户ID不存在",
		})
		return
	}
	if passwordOld != user.Password {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "原密码输入错误",
		})
		return
	}
	_, notFoundUserByName := service.QueryAUserByUsername(username)
	if !notFoundUserByName && username != user.Username {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户名已被占用",
		})
		return
	}
	err := service.UpdateAUser(&user, username, passwordNew, userInfo)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	//data, _ := jsoniter.Marshal(&user)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "修改成功",
		"data":    user,
	})
}

// TellUserInfo doc
// @description 查看用户个人信息
// @Tags 用户管理
// @Param user_id formData string true "用户ID"
// @Success 200 {string} string "{"success": true, "message": "查看用户信息成功", "data": "model.User的所有信息"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Router /user/info [POST]
func TellUserInfo(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(404, gin.H{
			"success": false,
			"message": "用户ID不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查看用户信息成功",
		"data":    user,
	})
}
