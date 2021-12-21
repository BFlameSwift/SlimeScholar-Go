package v1

import (
	"encoding/json"
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetScholar doc
// @description 获取学者信息
// @Tags 学者门户
// @Param user_id formData int false "user_id"
// @Param author_id formData int false "author_id"
// @Success 200 {string} string "{"success": true, "message": "用户验证邮箱成功"}"
// @Failure 401 {string} string "{"success": false, "message": "userid 不是整数"}"
// @Failure 402 {string} string "{"success": false, "message": "用户不是学者}"
// @Failure 404 {string} string "{"success": false, "message": "用户不存在}"
// @Failure 600 {string} string "{"success": false, "message": "用户待修改，传入false 更新验证码，否则为验证正确}"
// @Router /scholar/info [POST]
func GetScholar(c *gin.Context) {
	// TODO 根据实际的paper维护被引用数目等
	user_id_str := c.Request.FormValue("user_id")
	var ret_author_id string
	var people_msg map[string]interface{}
	var papers []interface{}
	is_user := false
	//var paper_result *elastic.SearchResult
	if user_id_str != "" {
		is_user = true
		user_id, err := strconv.ParseUint(user_id_str, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "user_id 不为整数", "status": 401})
			return
		}
		user, notFound := service.QueryAUserByID(user_id)
		if notFound {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户不存在", "status": 404})
			return
		}
		if user.UserType != 1 {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户不是学者", "status": 402})
			return
		}
		submit, notFound := service.SelectASubmitValid(user_id)
		if notFound {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "找不到用户的申请表", "status": 402})
			return
		}
		ret_author_id = submit.AuthorID
		fmt.Println("!!!!!!!")
		fmt.Println(service.GetAuthorCoAuthorIds(append(make([]string, 0), ret_author_id)))
		papers = service.GetAuthorSomePapers(ret_author_id, 100)
		//paper_result = service.QueryByField("paper", "authors.aid.keyword", submit.AuthorID, 1, 10)
		people_msg = service.UserScholarInfo(service.StructToMap(user), &papers)
		people_msg["follow_count"] = len(service.GetUserFollowedList(user.UserID))

	} else {
		author_id := c.Request.FormValue("author_id")
		var the_user_id uint64
		is_user, the_user_id = service.JudgeAuthorIsSettled(author_id)
		if author_id == "" {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "错误，authorid 与userid都不存在", "status": 401})
			return
		}
		ret_author_id = author_id
		papers = service.GetAuthorSomePapers(ret_author_id, 100)
		//paper_result = service.QueryByField("paper", "authors.aid.keyword", author_id, 1, 10)
		//people_msg = service.GetsByIndexIdWithout("author", author_id)
		if is_user {
			user, _ := service.QueryAUserByID(the_user_id)
			people_msg = service.UserScholarInfo(service.StructToMap(user), &papers)
			fmt.Println(people_msg)
			people_msg["follow_count"] = len(service.GetUserFollowedList(user.UserID))

		} else {
			people_msg = service.GetAuthorMsg(author_id)
			people_msg = service.ProcAuthorMsg(people_msg, &papers)
			people_msg["follow_count"] = 0
		}

	}
	//service.GetAuthorAllPaper(ret_author_id)
	CoauthorIds := service.GetSingleAuthorCoAuthorIds(ret_author_id)
	CoauthorItems := service.GetSimpleAuthors(CoauthorIds[ret_author_id].([]string))

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "status": 200, "is_user": is_user, "papers": papers, "author_id": ret_author_id, "people": people_msg, "coauthors": CoauthorItems})
	return
}

// ScholarManagePaper doc
// @description 学者添加或删除Paper,401 通常表示参数错误，Objw为被转让的人，当为添加或删除时，为零
// @Tags 学者门户
// @Param user_id formData string true "user_id"
// @Param obj_user_id formData string false "obj_user_id"
// @Param paper_id formData string true "paper_id"
// @Param kind formData int true "0添加1删除2转让"
// @Success 200 {string} string "{"success": true, "message": "转移成功"}"
// @Failure 401 {string} string "{"success": false, "message": "参数错误"}"
// @Failure 404 {string} string "{"success": false, "message": "用户不存在}"
// @Failure 600 {string} string "{"success": false, "message": "用户待修改，传入false 更新验证码，否则为验证正确}"
// @Router /scholar/transfer [POST]
func ScholarManagePaper(c *gin.Context) {

	user_id, err := strconv.ParseUint(c.Request.FormValue("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "user_id不为int", "status": 401})
		return
	}
	paper_id := c.Request.FormValue("paper_id")
	kind, err := strconv.Atoi(c.Request.FormValue("kind"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "kind 不是int", "status": 401})
		return
	}
	obj_user_id, err := strconv.ParseUint(c.Request.FormValue("user_id"), 10, 64)
	if err != nil {
		obj_user_id = 0
	}
	user, notFound := service.QueryAUserByID(user_id)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户不存在", "status": 404})
		return
	} else if user.UserType != 1 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户不是入驻学者", "status": 402})
		return
	} else if _, notFound = service.QueryASubmitExist(user_id); notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户提交表不存在", "status": 403})
		return
	}
	//fmt.Println(paper_id, is_add)
	service.TransferPaper(user, user.AuthorID, paper_id, kind, obj_user_id)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建成功", "status": 200})
	return
}

// FullPapersSocial doc
// @description 搜索时根据用户与paper_ids 来判断是否具有社交属性并补齐
// @Tags 社交
// @Param user_id formData string true "user_id"
// @Param paper_ids formData string true "paper_ids 文献id列表"
// @Success 200 {string} string "{"success": true, "message": "用户验证邮箱成功"}"
// @Failure 401 {string} string "{"success": false, "message": "参数格式错误"}"
// @Failure 404 {string} string "{"success": false, "message": "用户不存在}"
// @Failure 600 {string} string "{"success": false, "message": "用户待修改，传入false 更新验证码，否则为验证正确}"
// @Router /social/get/paper [POST]
func FullPapersSocial(c *gin.Context) {

	user_id, err := strconv.ParseUint(c.Request.FormValue("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "user_id不为int", "status": 401})
		return
	}
	paperIdsStr := c.Request.FormValue("paper_ids")
	paperIds := make([]string, 0)
	err = json.Unmarshal([]byte(paperIdsStr), &paperIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "paper—ids 格式错误", "status": 401})
		return
	}
	user, notFound := service.QueryAUserByID(user_id)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户不存在", "status": 404})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查询成功", "status": 200, "papers_attribute": service.PapersGetIsCollectedByUser(paperIds, user)})
	return
}

// CitePaper doc
// @description 根据paper_id 获取引用文献格式，返回的是一个字典数组
// @Tags 学者门户
// @Param paper_id formData string true "paper_id"
// @Success 200 {string} string "{"success": true, "message": "获取成共"}"
// @Failure 401 {string} string "{"success": false, "message": "参数格式错误"}"
// @Failure 404 {string} string "{"success": false, "message": "用户不存在}"
// @Failure 600 {string} string "{"success": false, "message": "用户待修改，传入false 更新验证码，否则为验证正确}"
// @Router /scholar/cite_paper [POST]
func CitePaper(c *gin.Context) {

	paper_id := c.Request.FormValue("paper_id")
	_, err := service.GetsByIndexId("paper", paper_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查询成功", "status": 200, "detail": service.CitePaper(paper_id)})
	return
}

// GetAuthorPartialCoAuthors doc
// @description 根據作者id获取作者的合作者
// @Tags 学者门户
// @Param author_id formData string true "author_id"
// @Param level formData string false "level"
// @Success 200 {string} string "{"success": true, "message": "获取成共"}"
// @Failure 401 {string} string "{"success": false, "message": "参数格式错误"}"
// @Failure 404 {string} string "{"success": false, "message": "用户不存在}"
// @Failure 600 {string} string "{"success": false, "message": "用户待修改，传入false 更新验证码，否则为验证正确}"
// @Router /scholar/graph [POST]
func GetAuthorPartialCoAuthors(c *gin.Context) {
	id, level := c.Request.FormValue("author_id"), c.Request.FormValue("level")

	coAuthorMap := service.GetSimpleAuthors(append(make([]string, 0), id))[0].(map[string]interface{})
	firstCoauthorIds := service.GetSingleAuthorCoAuthorIds(id)
	firstCoauthorItems := service.GetSimpleAuthors(firstCoauthorIds[id].([]string))
	if level != "2" {
		coAuthorMap["friends"] = firstCoauthorItems
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "查询成功", "status": 200, "detail": coAuthorMap})
		return
	}
	secondCoauthorIdMap := service.GetAuthorCoAuthorIds(firstCoauthorIds[id].([]string))
	for _, item := range firstCoauthorItems {
		this_id := item.(map[string]interface{})["author_id"].(string)
		if ids, ok := secondCoauthorIdMap[this_id]; ok {
			item.(map[string]interface{})["friends"] = service.GetSimpleAuthors(ids.([]string))
		}

	}
	coAuthorMap["friends"] = firstCoauthorItems
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查询成功", "status": 200, "detail": coAuthorMap})
	return
}
