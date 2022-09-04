package v1

import (
	"encoding/json"
	"fmt"
	"github.com/BFlameSwift/SlimeScholar-Go/service"
	"github.com/BFlameSwift/SlimeScholar-Go/utils"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"

	"net/http"
)

// Index doc
// @description 测试 Index 页
// @Tags basic
// @Success 200 {string} string "{"success": true, "message": "gcp"}"
// @Router / [GET]
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "gcp"})
}

// DocumentCount doc
// @description 在线返回所有文档数目
// @Tags basic
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
	if err := json.Unmarshal([]byte(service.GetUrl(utils.ELASTIC_SEARCH_HOST+"/affiliation/_count")), &affiliation_map); err != nil {
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

// UploadPdf doc
// @description  文件下载转换
// @Tags basic
// @Param pdf_url formData string true "文件路径"
// @Success 200 {string} string "{"success": true, "message": "上传成功",}"
// @Router /upload/get/pdf [POST]
func UploadPdf(c *gin.Context) {
	pdfUrl := c.Request.FormValue("pdf_url")
	// Get the data
	resp, err := http.Get(pdfUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "错误", "error": err.Error()})
		return
	}
	defer resp.Body.Close()
	// Create output file
	a := rand.Int()

	path := utils.UPLOAD_PATH + fmt.Sprintf("%d", a)[0:6] + ".pdf"
	out, err := os.Create(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": true, "message": "错误", "error": err.Error()})
		return
		//panic(err)
	}
	defer out.Close()
	// copy stream
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": true, "message": "错误", "error": err.Error()})
		return
		//panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "data": "/upload/media/" + fmt.Sprintf("%d", a)[0:6] + ".pdf"})
}
