package v1

import (
	"encoding/json"
	"gitee.com/online-publish/slime-scholar-go/service"
	"gitee.com/online-publish/slime-scholar-go/utils"
	"github.com/gin-gonic/gin"

	"net/http"
)

// Index doc
// @description 测试 Index 页
// @Tags 测试
// @Success 200 {string} string "{"success": true, "message": "gcp"}"
// @Router / [GET]
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "gcp"})
}

// Index doc
// @description 返回所有文档数目
// @Tags 测试
// @Success 200 {string} string "{"success": true, "message": "gcp"}"
// @Router /count/all [POST]
func DocumentCount(c *gin.Context) {
	paper_map := make(map[string]interface{})
	if err := json.Unmarshal([]byte(service.GetUrl(utils.ELASTIC_SEARCH_HOST+"/paper/_count")), &paper_map); err != nil {
		panic(err)
	}
	author_map := make(map[string]interface{})
	if err := json.Unmarshal([]byte(service.GetUrl(utils.ELASTIC_SEARCH_HOST+"/author/_count")), &author_map); err != nil {
		panic(err)
	}
	coference_map := make(map[string]interface{})
	if err := json.Unmarshal([]byte(service.GetUrl(utils.ELASTIC_SEARCH_HOST+"/conference/_count")), &coference_map); err != nil {
		panic(err)
	}
	journal_map := make(map[string]interface{})
	if err := json.Unmarshal([]byte(service.GetUrl(utils.ELASTIC_SEARCH_HOST+"/journal/_count")), &journal_map); err != nil {
		panic(err)
	}
	conference_map := make(map[string]interface{})
	if err := json.Unmarshal([]byte(service.GetUrl(utils.ELASTIC_SEARCH_HOST+"/conference/_count")), &conference_map); err != nil {
		panic(err)
	}
	fields_map := make(map[string]interface{})
	if err := json.Unmarshal([]byte(service.GetUrl(utils.ELASTIC_SEARCH_HOST+"/fields/_count")), &fields_map); err != nil {
		panic(err)
	}
	affiliation_map := make(map[string]interface{})
	if err := json.Unmarshal([]byte(service.GetUrl(utils.ELASTIC_SEARCH_HOST+"/fields/_count")), &affiliation_map); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200,
		"paper_count":       paper_map["count"],
		"author_count":      author_map["count"],
		"conference_count":  conference_map["count"],
		"journal_count":     journal_map["count"],
		"fields_count":      fields_map["count"],
		"affiliation_count": affiliation_map["count"],
	})

}
