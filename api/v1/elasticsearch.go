package v1

import (
	"encoding/json"
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/gin-gonic/gin"
)

// Confirm doc
// @description 创建es索引
// @Tags elasticsearch
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "创建成功"}"
// @Failure 401 {string} string "{"success": false, "message": "该ID已存在"}"
// @Failure 500 {string} string "{"success": false, "message": "创建错误500"}"
// @Router /es/create/mytype [POST]
func CreateMyType(c *gin.Context) {

	this_id := c.Request.FormValue("id")
	var mytype service.MyType
	mytype.Id = this_id
	json_str, err := json.Marshal(mytype)
	fmt.Println(json_str)
	if err != nil {
		// panic(err)
		c.JSON(400, gin.H{"success": false, "message": "创建错误"})
	}
	fmt.Println("json:", string(json_str))

	//service.Create("employee", "megacorp", string(this_id), string(json_str))
	return
}
