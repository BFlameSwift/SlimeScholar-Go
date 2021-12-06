package v1

import (
	"encoding/json"
	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Confirm doc
// @description 获取学者信息
// @Tags 学者门户
// @Param user_id formData int true "user_id"
// @Success 200 {string} string "{"success": true, "message": "用户验证邮箱成功"}"
// @Failure 401 {string} string "{"success": false, "message": "userid 不是整数"}"
// @Failure 402 {string} string "{"success": false, "message": "用户不是学者}"
// @Failure 404 {string} string "{"success": false, "message": "用户不存在}"
// @Failure 600 {string} string "{"success": false, "message": "用户待修改，传入false 更新验证码，否则为验证正确}"
// @Router /scholar/info [POST]
func GetScholar(c *gin.Context) {

	user_id, err := strconv.ParseUint(c.Request.FormValue("user_id"), 10, 64)
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
	paper_ids := make([]string, 0, 10000)
	paper_result := service.QueryByField("paper_author", "rel.aid.keyword", submit.AuthorID, 1, 10)
	authors_map := make(map[string]interface{})
	for _, hit := range paper_result.Hits.Hits {
		hit_map := make(map[string]interface{})
		err := json.Unmarshal([]byte(hit.Source), &hit_map)
		if err != nil {
			panic(err)
		}
		paper_ids = append(paper_ids, hit_map["paper_id"].(string))
		authors_map[hit_map["paper_id"].(string)] = service.ParseRelPaperAuthor(hit_map)
	}
	paper_id_map := service.IdsGetItems(paper_ids, "paper")
	paper_list := make([]interface{}, 0, 1000)
	for _, id := range paper_ids {
		if paper_id_map[id] != nil {
			paper_id_map[id].(map[string]interface{})["authors"] = authors_map[id]
			paper_list = append(paper_list, paper_id_map[id])
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "status": 200, "user": user, "papers": paper_list})
	return
}
