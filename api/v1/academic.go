package v1

import (
	"encoding/json"
	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"net/http"
	"strconv"
)

// Index doc
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
	var people_msg interface{}
	is_user := false
	var paper_result *elastic.SearchResult
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
		paper_result = service.QueryByField("paper", "authors.aid.keyword", submit.AuthorID, 1, 10)
		people_msg = service.UserScholarInfo(service.StructToMap(user))

	} else {
		author_id := c.Request.FormValue("author_id")
		if author_id == "" {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "错误，authorid 与userid都不存在", "status": 401})
			return
		}
		ret_author_id = author_id
		paper_result = service.QueryByField("paper", "authors.aid.keyword", author_id, 1, 10)
		//people_msg = service.GetsByIndexIdWithout("author", author_id)
		people_msg = "不是服务器es暂未存储"
	}

	paper_ids := make([]string, 0, 10000)
	//authors_map := make(map[string]interface{})
	for _, hit := range paper_result.Hits.Hits {
		hit_map := make(map[string]interface{})
		err := json.Unmarshal([]byte(hit.Source), &hit_map)
		if err != nil {
			panic(err)
		}
		paper_ids = append(paper_ids, hit_map["paper_id"].(string))
		//authors_map[hit_map["paper_id"].(string)] = service.ParseRelPaperAuthor(hit_map)
	}
	paper_id_map := service.IdsGetItems(paper_ids, "paper")
	paper_list := make([]interface{}, 0, 1000)
	for _, id := range paper_ids {
		if paper_id_map[id] != nil {
			authors_map := make(map[string]interface{})
			authors_map["rel"] = paper_id_map[id].(map[string]interface{})["authors"]
			paper_id_map[id].(map[string]interface{})["authors"] = service.ParseRelPaperAuthor(authors_map)["rel"]
			//paper_id_map[id].(map[string]interface{})["authors"] = authors_map[id].(map[string]interface{})["rel"]
			paper_list = append(paper_list, paper_id_map[id])
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "status": 200, "is_user": is_user, "papers": paper_list, "author_id": ret_author_id, "people": people_msg})
	return
}

// Index doc
// @description 学者添加或删除Paper,401 通常表示参数错误，Objw为被转让的人，当为添加或删除时，为零
// @Tags 学者门户
// @Param user_id formData string true "user_id"
// @Param obj_user_id formData string false "obj_user_id"
// @Param paper_id formData string true "paper_id"
// @Param kind formData int true "0添加1删除2转让"
// @Success 200 {string} string "{"success": true, "message": "用户验证邮箱成功"}"
// @Failure 401 {string} string "{"success": false, "message": "用户已验证邮箱"}"
// @Failure 402 {string} string "{"success": false, "message": "用户输入验证码错误}"
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
