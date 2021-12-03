package v1

import (
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/global"
	"gitee.com/online-publish/slime-scholar-go/model"
	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// Index doc
// @description 用户申请创建，402 用户id不是正忽视，404用户不存在，401 申请创建失败。后端炸了
// @Tags 管理员
// @Param author_name formData string true "作者姓名"
// @Param affiliation_name formData string true "机构姓名"
// @Param work_email formData string true "工作邮箱"
// @Param fields formData string true "领域"
// @Param home_page formData string true "主页"
// @Param author_id formData string true "作者id"
// @Param user_id formData string true "用户id"
// @Success 200 {string} string "{"success": true, "message": "创建成功"}"
// @Router /submit/create [POST]
func CreateSubmit(c *gin.Context) {
	author_name := c.Request.FormValue("author_name")
	affiliation_name := c.Request.FormValue("affiliation_name")
	work_email := c.Request.FormValue("work_email")
	fields := c.Request.FormValue("fields")
	home_page := c.Request.FormValue("home_page")
	author_id := c.Request.FormValue("author_id")
	user_id := c.Request.FormValue("user_id")
	user_id_u64, err := strconv.ParseUint(user_id, 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户ID不为正整数", "status": 402})
		return
	}
	_, notFound := service.QueryAUserByID(user_id_u64)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有该用户", "status": 404})
		return
	}

	submit := model.SubmitScholar{AffiliationName: affiliation_name, AuthorName: author_name, WorkEmail: work_email,
		HomePage: home_page, AuthorID: author_id, Fields: fields, UserID: user_id_u64, Status: 0, Content: "",
		CreatedTime: time.Now()}
	err = service.CreateASubmit(&submit)
	if err != nil {
		panic(err)
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "申请创建失败", "status": 401})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "申请提交成功", "status": 200})
	return
}

// Index doc
// @description 用户申请创建，401 402 用户id，提交id不是正整数，404提交不存在，405 用户不存在
// @Tags 管理员
// @Param submit_id formData string true "提交id"
// @Param user_id formData string true "用户id"
// @Param success formData string true "success"
// @Success 200 {string} string "{"success": true, "message": "创建成功"}"
// @Router /submit/check [POST]
func CheckSubmit(c *gin.Context) {
	submit_id := c.Request.FormValue("submit_id")
	user_id := c.Request.FormValue("user_id")
	success := c.Request.FormValue("success")
	submit_id_u64, err := strconv.ParseUint(submit_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "提交ID不为正整数", "status": 402})
		return
	}
	submit, notFound := service.QueryASubmitByID(submit_id_u64)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有该提交", "status": 404})
		return
	}
	user_id_u64, err := strconv.ParseUint(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户ID不为正整数", "status": 401})
		return
	}
	user, notFound := service.QueryAUserByID(user_id_u64)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有该用户", "status": 405})
		return
	}
	fmt.Println("check user submit", user.UserID)

	if success == "false" {
		submit.Status = 2
		service.SendCheckAnswer(user.Email, false)
	} else if success == "true" {
		submit.Status = 1
		service.SendCheckAnswer(user.Email, true)
		// TODO: 发邮件
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "success 不为true false", "status": 403})
		return
	}
	err = global.DB.Save(submit).Error
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "申请审批成功", "status": 200})
	return
}
