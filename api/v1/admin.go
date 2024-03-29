package v1

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/BFlameSwift/SlimeScholar-Go/global"
	"github.com/BFlameSwift/SlimeScholar-Go/model"
	"github.com/BFlameSwift/SlimeScholar-Go/service"
	"github.com/BFlameSwift/SlimeScholar-Go/utils"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Msg struct {
	time, msg string
}

// SubmitCount doc
// @description 获取统计信息
// @Tags 管理员
// @Success 200 {string} string "{"success": true, "message": "执行成功"}"
// @Router /submit/count [POST]
func SubmitCount(c *gin.Context) {
	data := make(map[string]interface{}, 0)

	paper_map := make(map[string]interface{})
	if err := json.Unmarshal([]byte(service.GetUrl(utils.ELASTIC_SEARCH_HOST+"/paper/_count")), &paper_map); err != nil {
		panic(err)
	}

	author_map := make(map[string]interface{})
	if err := json.Unmarshal([]byte(service.GetUrl(utils.ELASTIC_SEARCH_HOST+"/author/_count")), &author_map); err != nil {
		panic(err)
	}
	data["literCount"] = paper_map["count"]
	// fmt.Println(paper_map["count"])
	data["authorCount"] = author_map["count"]
	// fmt.Println(author_map["count"])

	userCount, memberCount := service.QueryUserCount()
	data["userCount"] = userCount
	// fmt.Println(data["userCount"])
	data["memberCount"] = memberCount
	// fmt.Println(data["memberCount"])

	filename := utils.LOG_FILE_PATH + utils.LOG_FILE_NAME
	activeIndex, responseTime := LogAnalize(filename)
	data["activeIndex"] = activeIndex
	fmt.Println(data["activeIndex"])

	data["responseTime"] = responseTime
	fmt.Println(data["responseTime"])

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "data": data})
}

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
	// 创建申请称为学者的表，等待平台管理申请
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
	if _, notFound := service.QueryUserIsScholar(user_id_u64); !notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "您已经是认证学者，请勿重复申请", "status": 406})
		return
	}
	//后续对papers可能需要处理
	//papers := service.GetAuthorAllPaper(author_id)
	author := service.GetAuthors(append(make([]string, 0), author_id))[0].(map[string]interface{})
	submit := model.SubmitScholar{AffiliationName: affiliation_name, AuthorName: author_name, WorkEmail: work_email,
		HomePage: home_page, AuthorID: author_id, Fields: fields, UserID: user_id_u64, Status: 0, Content: "", PaperCount: int(author["paper_count"].(float64)),
		CreatedTime: time.Now()}

	err = service.CreateASubmit(&submit)
	if err != nil {
		panic(err)
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "申请创建失败", "status": 401})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "申请提交成功", "status": 200, "papers": service.GetAuthorSomePapers(author_id, 100)})
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
	// 审核申请
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

	if submit.Status != 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "已审核过该申请", "status": 406})
		return
	}
	// 根据审批的同意与否发送对应信息
	if success == "false" {
		submit.Status = 2
		submit.Content = content
		service.SendCheckAnswer(user.Email, false, content)
	} else if success == "true" {
		submit.Status = 1
		submit.Content = content
		// 使用户成为学者，补全user信息
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
// @description 通过或拒绝多条申请。406-没有需要审批的申请
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

	if success != "false" && success != "true" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "success 不为true false", "status": 403})
		return
	}

	submit_ids := strings.Split(submit_ids_str, `,`)
	len := len(submit_ids)
	fmt.Println(len)
	fmt.Println(submit_ids)

	for _, tmp := range submit_ids {
		submit_id, _ := strconv.ParseUint(tmp, 0, 64)
		submit, notFound := service.QueryASubmitByID(submit_id)
		if notFound || submit.Status != 0 {
			len--
			continue
		}
		fmt.Println(len)
		user, _ := service.QueryAUserByID(submit.UserID)
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
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有需要审批的申请", "status": 406})
		return
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

	submits := make([]model.SubmitScholar, 0)
	if mytype_str == "" || len(mytype_str) == 0 {
		submits = service.QueryAllSubmit()
	} else {
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
		// accept_time 是sql.Nulltime 的格式，以下的操作只是为了将这个格式转化为要求的格式罢了
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
// @description 根据作者姓名返回姓名相近的作者并返回文献组，申请学者后，需要认领一个作者，此处返回对应文献组供选择某一个作者来认领
// @Tags 管理员
// @Param author_name formData string true "author_name"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Success 200 {string} string "{"success": true, "message": "创建成功"}"
// @Router /submit/get/papers [POST]
func PaperGetAuthors(c *gin.Context) {
	//
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
	searchResult, err := service.Client.Search().Index("author").Query(elastic.NewMatchPhraseQuery("name", author_name)).From((page - 1) * size).Size(size).Do(context.Background())
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
		papers := service.GetAuthorSomePapers(author_map["author_id"].(string), 30)
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

// GetSubmitDetail doc
// @description 获取入驻申请详细信息
// @Tags 管理员
// @Param submit_id formData string true "提交id"
// @Success 200 {string} string "{"success": true, "message": "信息获取成功", "data": data}"
// @Router /submit/get/detail [POST]
func GetSubmitDetail(c *gin.Context) {
	submit_id_u64, err := strconv.ParseUint(c.Request.FormValue("submit_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "提交ID不为正整数", "status": 402})
		return
	}
	submit, notFound := service.QueryASubmitByID(submit_id_u64)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有该提交", "status": 404})
		return
	}

	data := make(map[string]interface{})
	data["author_id"] = submit.AuthorID
	data["real_name"] = submit.AuthorName
	data["work_email"] = submit.WorkEmail
	data["affiliation_name"] = submit.AffiliationName
	data["homepage"] = submit.HomePage
	data["papers"] = make([]map[string]interface{}, 0)
	data["fields"] = strings.Split(submit.Fields, `,`)

	author_id := submit.AuthorID
	papers := service.GetAuthorSomePapers(author_id, 30)
	data_papers := make([]map[string]interface{}, 0)
	// fmt.Println(papers)
	for _, tmp := range papers {
		paper := make(map[string]interface{})
		fmt.Println(tmp.(map[string]interface{})["paper_id"])
		paper["paper_id"] = tmp.(map[string]interface{})["paper_id"]
		paper["paper_title"] = tmp.(map[string]interface{})["paper_title"]
		paper["publisher"] = tmp.(map[string]interface{})["publisher"]
		paper["year"] = tmp.(map[string]interface{})["year"]
		paper["authors"] = tmp.(map[string]interface{})["authors"]
		data_papers = append(data_papers, paper)
	}
	data["papers"] = data_papers

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "获取成功", "status": 200, "data": data})

}

// AdminLogin doc
// @description 登录 200-成功	401-用户不存在	402-密码错误	403-用户尚未确认邮箱	405-用户不是管理员
// @Tags 管理员
// @Param username formData string false "用户名"
// @Param email formData string false "用户邮箱"
// @Param password formData string true "密码"
// @Success 200 {string} string "{"success": true, "message": "登录成功", "detail": user的信息}"
// @Router /submit/login [POST]
func AdminLogin(c *gin.Context) {
	username := c.Request.FormValue("username")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	user, notFound := model.User{}, true
	if username != "" {
		user, notFound = service.QueryAUserByUsername(username)
	} else {
		user, notFound = service.QueryAUserByEmail(email)
	}
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户不存在", "status": 401})
	} else {
		if user.Password != password {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "密码错误", "status": 402})
		} else {
			if user.UserType != 2 {
				c.JSON(http.StatusOK, gin.H{"success": false, "message": "该用户不是管理员", "status": 405})
			} else {
				if user.HasConfirmed == false {
					c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户尚未确认邮箱", "status": 403})
				} else {
					claims := &model.JWTClaims{
						UserID:   user.UserID,
						Username: user.Username,
						Password: password,
					}
					claims.IssuedAt = time.Now().Unix()
					claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(utils.ExpireTime)).Unix()
					signedToken, err := service.GetToken(claims)
					if err != nil {
						c.String(http.StatusNotFound, err.Error())
						return
					}
					c.JSON(http.StatusOK, gin.H{
						"success":       true,
						"message":       "登录成功",
						"detail":        user,
						"status":        200,
						"Authorization": signedToken})
				}
			}
		}
	}
}

// 分析日志，获得管理员需要的统计数据
func LogAnalize(filename string) (data []interface{}, resTime float64) {

	f, e := os.Open(filename)
	var msgList []Msg
	if e != nil {
		fmt.Println("File error.")
	} else {
		buf := bufio.NewScanner(f)
		for {
			if !buf.Scan() {
				break
			}
			line := buf.Text()
			line = strings.TrimSpace(line) //去掉前后空格
			// line_list := strings.Split(line, ` level=info `)
			var tmp Msg
			line_map := make(map[string]string)
			_ = json.Unmarshal([]byte(line), &line_map)
			time, _ := time.ParseInLocation("2006-01-02 15:04:05", line_map["time"], time.Local)
			tmp.time = time.Format("2006-01-02")
			tmp.msg = line_map["msg"]
			msgList = append(msgList, tmp)
		}
	}

	//获取日活跃数
	map_tmp := make(map[string]interface{}, 0)
	dest := make([]map[string]interface{}, 0)
	for i, _ := range msgList {
		ai := msgList[i]
		if _, ok := map_tmp[ai.time]; !ok {
			tmp := make(map[string]interface{}, 0)
			tmp["time"] = ai.time
			tmp["count"] = 1
			dest = append(dest, tmp)
			map_tmp[ai.time] = ai
		} else {
			for j, _ := range dest {
				var dj = dest[j]
				if dj["time"].(string) == ai.time {
					count := dj["count"].(int)
					dj["count"] = count + 1
					dest[j] = dj
					break
				}
			}
		}
	}
	data = make([]interface{}, 0)
	for _, tmp := range dest {
		var a [2]interface{}
		a[0] = tmp["time"]
		a[1] = tmp["count"]
		data = append(data, a)
	}
	//

	//获取响应时间
	Reverse(&msgList)
	msg_count := 0 //记录最近100条POST信息
	resTime = 0.0
	for _, tmp := range msgList {
		if strings.Contains(tmp.msg, "POST") && msg_count < 100 {
			a := tmp.msg
			fmt.Println(a)
			end := strings.Index(a, "s")
			b := strings.TrimSpace(a[7:end]) //去掉前后空格
			if a[end-1] == 'm' {
				len := len(b)
				for i := len - 1; i >= 0; i-- {
					if !(b[i] >= '0' && b[i] <= '9') {
						len--
					} else {
						break
					}
				}
				c, _ := strconv.ParseFloat(b[:len], 64)
				if c > 10.0 {
					fmt.Println(c)
					resTime = resTime + c
					msg_count++
					if msg_count >= 100 {
						break
					}
				}
			}
		}
	}
	count, _ := strconv.ParseFloat(strconv.Itoa(msg_count), 64)
	resTime = resTime / count
	resTime, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", resTime), 64)
	//

	return data, resTime
}
func Reverse(arr *[]Msg) {
	var temp Msg
	length := len(*arr)
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}
