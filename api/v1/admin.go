package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/global"
	"gitee.com/online-publish/slime-scholar-go/model"
	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CreateSubmit doc
// @description 用户申请创建，402 用户id不是正整数，404用户不存在，401 申请创建失败。后端炸了，405！！！该作者已被成功认领,并直接返回认领了该作者的学者姓名，406 该用户已经提交过对该学者的认领
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
	if the_submit, notFound := service.QueryASubmitByAuthor(author_id); !notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "该作者已被认领", "status": 405, "the_authorname": the_submit.AuthorName})
		return
	}
	if _, notFound := service.QueryASubmitExist(user_id_u64); !notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "您已经是认证学者，请勿重复申请", "status": 406})
		return
	}
	// TODO Paper COunt根据实际输入。：
	submit := model.SubmitScholar{AffiliationName: affiliation_name, AuthorName: author_name, WorkEmail: work_email,
		HomePage: home_page, AuthorID: author_id, Fields: fields, UserID: user_id_u64, Status: 0, Content: "", PaperCount: 12,
		CreatedTime: time.Now()}

	err = service.CreateASubmit(&submit)
	if err != nil {
		panic(err)
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "申请创建失败", "status": 401})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "申请提交成功", "status": 200, "papers": service.GetAuthorAllPaper(author_id)})
	return
}

// CheckSubmit doc
// @description 通过或拒绝某一条申请，401 402 用户id，提交id不是正整数，404提交不存在，405 用户不存在，406-已审核过该申请
// @Tags 管理员
// @Param submit_id formData string true "提交id"
// @Param user_id formData string true "用户id"
// @Param success formData string true "success"
// @Param content formData string false "content"
// @Success 200 {string} string "{"success": true, "message": "创建成功"}"
// @Router /submit/check [POST]
func CheckSubmit(c *gin.Context) {
	submit_id := c.Request.FormValue("submit_id")
	user_id := c.Request.FormValue("user_id")
	success := c.Request.FormValue("success")
	content := c.Request.FormValue("content")
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

	if submit.Status != 0{
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "已审核过该申请", "status": 406})
		return
	}

	if success == "false" {
		submit.Status = 2
		submit.Content = content
		service.SendCheckAnswer(user.Email, false, content)
	} else if success == "true" {
		submit.Status = 1
		submit.Content = content
		service.MakeUserScholar(user, submit)
		service.SendCheckAnswer(user.Email, true, content)
		submit.AcceptTime = sql.NullTime{Time: time.Now(), Valid: true}

	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "success 不为true false", "status": 403})
		return
	}

	err = global.DB.Save(submit).Error
	fmt.Println(submit.AcceptTime)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "申请审批成功", "status": 200})
	return
}

// CheckSubmits doc
// @description 通过或拒绝某一条申请。402-没有需要审批的申请
// @Tags 管理员
// @Param submit_ids formData string true "提交id"
// @Param success formData string true "success"
// @Param content formData string false "content"
// @Success 200 {string} string "{"success": true, "message": "创建成功"}"
// @Router /submit/check/more [POST]
func CheckSubmits(c *gin.Context) {
	submit_ids_str := c.Request.FormValue("submit_ids")
	success := c.Request.FormValue("success")
	content := c.Request.FormValue("content")

	if success != "false" && success != "true"{
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "success 不为true false", "status": 403})
		return
	}

	submit_ids := strings.Split(submit_ids_str, `,`)
	len := len(submit_ids)
	fmt.Println(len)
	fmt.Println(submit_ids)

	for _,tmp := range submit_ids{
		submit_id,_ := strconv.ParseUint(tmp, 0, 64)
		submit, notFound := service.QueryASubmitByID(submit_id)
		if notFound || submit.Status != 0 {
			len--
			continue
		}
		fmt.Println(len)
		user,_ := service.QueryAUserByID(submit.UserID)
		if success == "false" {
			submit.Status = 2
			submit.Content = content
			service.SendCheckAnswer(user.Email, false, content)
		} else if success == "true" {
			submit.Status = 1
			submit.Content = content
			service.MakeUserScholar(user, submit)
			service.SendCheckAnswer(user.Email, true, content)
			submit.AcceptTime = sql.NullTime{Time: time.Now(), Valid: true}
		}

		err := global.DB.Save(submit).Error
		fmt.Println(submit.AcceptTime)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(len)
	if len == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有需要审批的申请", "status": 402})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "申请审批成功", "status": 200})
	return
}



// ListAllSubmit doc
// @description 列举出所有type类型的submit，0表示未审批的，1表示审批成功的，2表示审批失败的；不输入type，则返回所有申请
// @Tags 管理员
// @Param type formData int false "提交id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Router /submit/list [POST]
func ListAllSubmit(c *gin.Context) {
	mytype_str := c.Request.FormValue("type")

	submits := make([]model.SubmitScholar,0)
	if mytype_str == "" || len(mytype_str) == 0{
		submits = service.QueryAllSubmit()
	}else{
		mytype, err := strconv.Atoi(mytype_str)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "type不为正整数", "status": 401})
			return
		}
		submits, _ = service.QuerySubmitByType(mytype)
	}

	submits_arr := make([]interface{}, 0)
	var err error
	for _, obj := range submits {
		// accept_time 是sql。Nulltime h格式，一下的操作只是为了将这个格式转化为要求的格式罢了
		obj_json, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		submit_map := make(map[string]interface{})
		err = json.Unmarshal(obj_json, &submit_map)
		submit_map["accept_time"] = submit_map["accept_time"].(map[string]interface{})["Time"]
		if strings.Index(submit_map["accept_time"].(string), "000") == 0 {
			submit_map["accept_time"] = ""
		}
		submits_arr = append(submits_arr, submit_map)
	}
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "获取成功", "status": 200, "submits": submits_arr, "submit_count": len(submits)})
	return
}

// PaperGetAuthors doc
// @description 根据作者姓名返回姓名相近的作者并返回文献组
// @Tags 管理员
// @Param author_name formData string true "author_name"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Success 200 {string} string "{"success": true, "message": "创建成功"}"
// @Router /submit/get/papers [POST]
func PaperGetAuthors(c *gin.Context) {
	author_name := c.Request.FormValue("author_name")
	page, err := strconv.Atoi(c.Request.FormValue("page"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "page 不为整数", "status": 401})
		return
	}
	size, err := strconv.Atoi(c.Request.FormValue("size"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "size 不为整数", "status": 401})
		return
	}
	searchResult, err := service.Client.Search().Index("author").Query(elastic.NewMatchQuery("name", author_name)).From((page - 1) * size).Size(size).Do(context.Background())
	if err != nil {
		panic(err)
	}
	author_maps := make([]map[string]interface{}, 0, 10)
	for _, hit := range searchResult.Hits.Hits {
		author_map := make(map[string]interface{})
		err = json.Unmarshal(hit.Source, &author_map)
		if err != nil {
			panic(err)
		}
		papers := service.GetAuthorAllPaper(author_map["author_id"].(string))
		if papers == nil {
			author_map["papers"] = make([]string, 0)
		} else {
			author_map["papers"] = papers
		}

		author_maps = append(author_maps, author_map)
	}
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "未找到该作者", "status": 404})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "获取成功", "status": 200, "authors": author_maps, "author_count": searchResult.TotalHits()})
	return
}
