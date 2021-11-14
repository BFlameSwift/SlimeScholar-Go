package v1

import (
	"encoding/json"
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/model"
	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TestCreate doc
// @description 创建es索引
// @Tags elasticsearch
// @Param id formData string true "id"
// @Param id formData string true "intvalue"
// @Success 200 {string} string "{"success": true, "message": "创建成功"}"
// @Failure 401 {string} string "{"success": false, "message": "该ID已存在"}"
// @Failure 500 {string} string "{"success": false, "message": "创建错误500"}"
// @Router /es/create/mytype [POST]
func CreateMyType(c *gin.Context) {
	c.Header("content-type","application/json")
	this_id := c.Request.FormValue("id")
	var mytype model.ValueString
	mytype.Value = this_id; mytype.Stuid = 200
	json_byte, _ := json.Marshal(mytype)
	fmt.Println(string(json_byte))
	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["type"], map_param["id"], map_param["bodyJson"] = "mytype", "mytype", this_id, string(json_byte)

	get1, error_get := service.Gets(map_param)
	if error_get == nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引已存在","status":401})

		obj_byte, _ := json.Marshal(get1.Source)

		fmt.Println("field",get1.Fields)
		fmt.Println("this id "+get1.Id+"has existed",string(obj_byte))
		return
	}
	ret := service.Create(map_param)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建成功"+ret,"status":200})
	return
}

// UpdateMyType doc
// @description 更新es索引
// @Tags elasticsearch
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "更新成功"}"
// @Failure 404 {string} string "{"success": false, "message": "该ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "创建错误500"}"
// @Router /es/update/mytype [POST]
func UpdateMyType(c *gin.Context) {
	this_id := c.Request.FormValue("id")
	var mytype model.ValueString
	mytype.Value = this_id
	json_str, _ := json.Marshal(mytype)
	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["type"], map_param["id"], map_param["bodyJson"] = "mytype", "mytype", this_id, string(json_str)

	_, error_get := service.Gets(map_param)
	if error_get != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在","status":404})
		fmt.Println("this id %s not existed",this_id)
		return
	}
	ret := service.Update(map_param)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "更新成功"+ret,"status":200})
	return
}

// GetMyType doc
// @description 获取es索引
// @Tags elasticsearch
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": "该ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/get/mytype [POST]
func GetMyType(c *gin.Context) {
	this_id := c.Request.FormValue("id")
	var mytype model.ValueString
	mytype.Value = this_id
	json_str, _ := json.Marshal(mytype)
	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["type"], map_param["id"], map_param["bodyJson"] = "mytype", "mytype", this_id, string(json_str)

	_, error_get := service.Gets(map_param)
	if error_get != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在","status":404})
		fmt.Println("this id %s not existed",this_id)
		return
	}
	ret,_ := service.Gets(map_param)
	body_byte,_ := json.Marshal(ret.Source)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功","status":200,"details":string(body_byte)})
	return
}

// GetMyType doc
// @description 获取es作者
// @Tags elasticsearch
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": "该ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/get/author [POST]
func GetAuthor(c *gin.Context) {
	this_id := c.Request.FormValue("id")

	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["type"], map_param["id"], map_param["bodyJson"] = "author", "mytype", this_id, ""

	_, error_get := service.Gets(map_param)
	if error_get != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在","status":404})
		fmt.Println("this id %s not existed",this_id)
		return
	}
	ret,_ := service.Gets(map_param)
	var author_map map[string]interface{} = make(map[string]interface{})
	body_byte,_ := json.Marshal(ret.Source)
	err := json.Unmarshal(body_byte,&author_map)
	if err != nil {panic(err)}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功","status":200,"details":author_map})
	return
}