package v1

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/olivere/elastic/v7"

	"github.com/BFlameSwift/SlimeScholar-Go/service"
	"github.com/BFlameSwift/SlimeScholar-Go/utils"
	"github.com/gin-gonic/gin"
	"github.com/wxnacy/wgo/arrays"
	"golang.org/x/net/context"
)

// GetPaper doc
// @description es获取Paper详细信息
// @Tags elasticsearch
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": "该PaperID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/get/paper [POST]
func GetPaper(c *gin.Context) {
	thisId := c.Request.FormValue("id")
	var mapParam map[string]string = make(map[string]string)

	mapParam["index"], mapParam["id"] = "paper", thisId
	_, errorGet := service.Gets(mapParam)
	if errorGet != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在", "status": 404})
		fmt.Printf("this id %s not existed", thisId)
		return
	}
	// GetFullPaper 关键，获取文献的全部详细信息
	paper := service.GetFullPaper(thisId)
	// 补齐paper的social信息
	paper = service.FullPaperSocial(paper)
	//service.CitePaper(thisId)

	//yearList, citationCountList := service.GetCitationPapersGraph(append(make([]string, 0), thisId), 200)
	////fmt.Println(citations)
	////fmt.Println(len(citations))
	//fmt.Println(yearList)
	//fmt.Println(citationCountList)
	//fmt.Println(len(yearList))

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "details": paper})

	return
}

// GetAuthor doc
// @description 获取es作者
// @Tags elasticsearch
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": "该ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/get/author [POST]
func GetAuthor(c *gin.Context) {
	// 根据id获取作者信息，不过已经被GetScholar所取代
	thisId := c.Request.FormValue("id")
	var mapParam map[string]string = make(map[string]string)
	mapParam["index"], mapParam["id"], mapParam["bodyJson"] = "author", thisId, ""

	_, errorGet := service.Gets(mapParam)
	if errorGet != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在", "status": 404})
		fmt.Println("this id %s not existed", thisId)
		return
	}
	ret, _ := service.Gets(mapParam)
	var authorMap map[string]interface{} = make(map[string]interface{})
	bodyByte, _ := json.Marshal(ret.Source)
	err := json.Unmarshal(bodyByte, &authorMap)
	if err != nil {
		panic(err)
	}
	fmt.Println(authorMap)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "details": authorMap, "body": string(bodyByte)})
	return
}

// GetAffiliation doc
// @description 获取es会议详细信息
// @Tags elasticsearch
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": 会议ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/get/affiliation [POST]
func GetAffiliation(c *gin.Context) {
	// 根据id获取机构信息，不过最后未使用
	thisId := c.Request.FormValue("id")
	var mapParam map[string]string = make(map[string]string)
	mapParam["index"], mapParam["id"], mapParam["bodyJson"] = "affiliation", thisId, ""

	_, errorGet := service.Gets(mapParam)
	if errorGet != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在", "status": 404})
		fmt.Println("this id %s not existed", thisId)
		return
	}
	ret, _ := service.Gets(mapParam)
	var affiliationMap map[string]interface{} = make(map[string]interface{})
	bodyByte, _ := json.Marshal(ret.Source)
	err := json.Unmarshal(bodyByte, &affiliationMap)
	if err != nil {
		panic(err)
	}
	//fmt.Println(affiliationMap)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "details": affiliationMap})
	return
}

// GetConference doc
// @description 获取es会议详细信息
// @Tags elasticsearch
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": 会议ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/get/conference [POST]
func GetConference(c *gin.Context) {
	// 根据id获取会议信息，不过已经被GetScholar所取代
	thisId := c.Request.FormValue("id")
	var mapParam map[string]string = make(map[string]string)
	mapParam["index"], mapParam["id"], mapParam["bodyJson"] = "conference", thisId, ""

	_, errorGet := service.Gets(mapParam)
	if errorGet != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在", "status": 404})
		fmt.Println("this id %s not existed", thisId)
		return
	}
	ret, _ := service.Gets(mapParam)
	var conferenceMap map[string]interface{} = make(map[string]interface{})
	bodyByte, _ := json.Marshal(ret.Source)
	err := json.Unmarshal(bodyByte, &conferenceMap)
	if err != nil {
		panic(err)
	}
	fmt.Println(conferenceMap)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "details": conferenceMap})
	return
}

// GetJournal doc
// @description 获取es期刊详细信息
// @Tags elasticsearch
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": 期刊ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/get/journal [POST]
func GetJournal(c *gin.Context) {
	// 根据id获取Journal信息，不过已经被GetScholar所取代
	thisId := c.Request.FormValue("id")
	var mapParam map[string]string = make(map[string]string)
	mapParam["index"], mapParam["id"], mapParam["bodyJson"] = "journal", thisId, ""

	_, errorGet := service.Gets(mapParam)
	if errorGet != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在", "status": 404})
		fmt.Println("this id %s not existed", thisId)
		return
	}
	ret, _ := service.Gets(mapParam)
	var journalMap map[string]interface{} = make(map[string]interface{})
	bodyByte, _ := json.Marshal(ret.Source)
	err := json.Unmarshal(bodyByte, &journalMap)
	if err != nil {
		panic(err)
	}
	fmt.Println(journalMap)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "details": journalMap})
	return
}

// TitleQueryPaper doc
// @description es 根据title查询论文
// @Tags elasticsearch
// @Param title formData string true "title"
// @Param page formData int true "page"
// @Param is_precise formData bool false "is_precise"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 401 {string} string "{"success": false, "message": "page 不是整数"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/title [POST]
func TitleQueryPaper(c *gin.Context) {
	// 多表联查，查id的时候同时查询author，  查个屁（父子文档开销太大，扁平化管理了
	title := c.Request.FormValue("title")
	page, err := strconv.Atoi(c.Request.FormValue("page"))
	isPreciseStr := c.Request.FormValue("is_precise")
	is_precise := true
	if isPreciseStr == "true" {
		is_precise = false
	}
	title = strings.ToLower(title)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "page 不为整数", "status": 401})
		return
	}
	searchResult := service.PaperQueryByField("paper", "paper_title", title, page, 10, is_precise, elastic.NewBoolQuery(), 1, true)

	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this title query %s not existed", title)
		return
	}
	fmt.Println("search title", title, "hits :", searchResult.TotalHits())

	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	// 关键的GetPapers，**批量**的根据10个id获取paper的详细信息，尽量减少开销
	paperSequences = service.GetPapers(paperIds)

	//其他的aggregation都集成起来了，毕竟每一个都查询都十行代码挺臭的
	aggregation := make(map[string]interface{})

	aggregation["doctype"] = service.Paper_Aggregattion(searchResult, "doctype")
	fmt.Println(aggregation["doctype"])
	aggregation["journal"] = service.Paper_Aggregattion(searchResult, "journal")
	aggregation["conference"] = service.Paper_Aggregattion(searchResult, "conference")
	aggregation["fields"] = service.Paper_Aggregattion(searchResult, "fields")
	aggregation["publisher"] = service.Paper_Aggregattion(searchResult, "publisher")
	aggregation["min_year"] = service.GetAggregationYear(searchResult, "min_year", true)
	aggregation["max_year"] = service.GetAggregationYear(searchResult, "max_year", false)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": aggregation})
	return
}

// TitleSelectPaper doc
// @description es 根据title筛选论文，包括对文章类型journal的筛选，页数的更换,页面大小size的设计, \n 错误码：401 参数格式错误, 排序方式1为默认，2为引用率，3为年份
// @Tags elasticsearch
// @Param title formData string true "title"
// @Param is_precise formData bool false "is_precise"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Param min_year formData int true "min_year"
// @Param max_year formData int true "max_year"
// @Param doctypes formData string true "doctypes"
// @Param conferences formData string true "conferences"
// @Param journals formData string true "journals"
// @Param publishers formData string true "publishers"
// @Param sort_type formData int true "sort_type"
// @Param sort_ascending formData bool true "sort_ascending"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 401 {string} string "{"success": false, "message": "page 不是整数"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/select/paper/title [POST]
func TitleSelectPaper(c *gin.Context) {
	//曾想多表联查，查id的时候同时查询author，  然而查个屁（父子文档开销太大，扁平化管理了

	var sort_ascending bool
	title := c.Request.FormValue("title")

	page_str := c.Request.FormValue("page")
	size_str := c.Request.FormValue("size")
	min_year := c.Request.FormValue("min_year")
	max_year := c.Request.FormValue("max_year")
	doctypesJson, journalsJson, conferenceJson, publisherJson := c.Request.FormValue("doctypes"), c.Request.FormValue("journals"), c.Request.FormValue("conferences"), c.Request.FormValue("publishers")
	doctypes, conferences, journals, publishers := make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100)
	sort_type_str := c.Request.FormValue("sort_type")
	sort_ascending_str := c.Request.FormValue("sort_ascending")
	isPreciseStr := c.Request.FormValue("is_precise")
	is_precise := true

	err := service.CheckSelectPaperParams(c, page_str, size_str, min_year, max_year, doctypesJson, journalsJson, conferenceJson, publisherJson, sort_ascending_str)
	if err != nil {
		// 参数校验401错误
		return
	}
	sort_ascending, _ = strconv.ParseBool(sort_ascending_str)

	page, size, sort_type := service.PureAtoi(page_str), service.PureAtoi(size_str), service.PureAtoi(sort_type_str)
	json.Unmarshal([]byte(doctypesJson), &doctypes)
	json.Unmarshal([]byte(journalsJson), &journals)
	json.Unmarshal([]byte(conferenceJson), &conferences)
	json.Unmarshal([]byte(publisherJson), &publishers)
	if isPreciseStr == "true" {
		is_precise = false
	}
	title = strings.ToLower(title)
	boolQuery := service.SelectTypeQuery(doctypes, journals, conferences, publishers, service.PureAtoi(min_year), service.PureAtoi(max_year))
	//searchResult := service.SearchSort(boolQuery, sort_type, sort_ascending, page, size)
	searchResult := service.PaperQueryByField("paper", "paper_title", title, page, size, is_precise, boolQuery, sort_type, sort_ascending)

	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this title query %s not existed", title)
		return
	}
	fmt.Println("search title", title, "hits :", searchResult.TotalHits())

	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// NameQueryAuthor doc
// @description es 根据姓名查询作者：精确查询,is_precise=0 为模糊匹配，为1为精准匹配
// @Tags elasticsearch
// @Param author_name formData string true "author_name"
// @Param sort_type formData int true "排序方式，1,代表按照论文数量排序，2代表按照被引用书目排序,0 为默认"
// @Param sort_ascending formData bool true "sort_ascending"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Param affiliations formData string true "列表形式，对结果按照机构进行筛选,不筛选传空列表,为机构id的列表"
// @Success 200 {string} string "{"success": true, "message": "获取作者成功"}"
// @Failure 404 {string} string "{"success": false, "message": "作者不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/author/name [POST]
func NameQueryAuthor(c *gin.Context) {
	name := c.Request.FormValue("author_name")
	name = strings.ToLower(name)
	page, _ := strconv.Atoi(c.Request.FormValue("page"))
	size, _ := strconv.Atoi(c.Request.FormValue("size"))
	sort_type, _ := strconv.Atoi(c.Request.FormValue("sort_type"))
	sort_ascending, _ := strconv.ParseBool(c.Request.FormValue("sort_ascending"))
	affiliation_name := c.Request.FormValue("affiliations")
	affiliations := make([]string, 0)
	err := json.Unmarshal([]byte(affiliation_name), &affiliations)
	if err != nil || (sort_type != 0 && sort_type != 2 && sort_type != 1) {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "机构格式错误", "status": 401})
		return
	}

	boolQuery := elastic.NewBoolQuery()
	query := elastic.NewMatchPhraseQuery("name", name)
	boolQuery.Must(query)
	orQuery := elastic.NewBoolQuery()
	for _, affiliation := range affiliations {
		fmt.Println(affiliation)
		orQuery.Should(elastic.NewMatchPhraseQuery("affiliation_id.keyword", affiliation))
	}
	boolQuery.Must(orQuery)

	searchResult := service.AuthorQuery(page, size, sort_type, sort_ascending, "author", boolQuery)

	aggregation := make(map[string]interface{})
	aggregation["affiliations"] = service.Paper_Aggregattion(searchResult, "affiliation")
	//aggregation["author"] = service.Paper_Aggregattion(searchResult, "author")
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this authors query %s not existed", name)
		return
	}
	fmt.Println("search author", name, "hits :", searchResult.TotalHits())
	ids := make([]string, 0)

	for _, paper := range searchResult.Hits.Hits {
		ids = append(ids, paper.Id)
	}
	authors := service.GetAuthors(ids)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": *service.ParseEnterScholarMsg(&authors), "aggregation": aggregation})
	return
}

// AffiliationNameQueryAuthor doc
// @description es 根据机构姓名查询作者：
// @Tags elasticsearch
// @Param affiliation_name formData string true "affiliation_name机构名称"
// @Param sort_type formData int true "排序方式，1,代表按照论文数量排序，2代表按照被引用书目排序,0 为默认"
// @Param sort_ascending formData bool true "sort_ascending"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Param affiliations formData string true "列表形式，对结果按照机构进行筛选,不筛选传空列表,为机构id的列表"
// @Success 200 {string} string "{"success": true, "message": "获取作者成功"}"
// @Failure 404 {string} string "{"success": false, "message": "作者不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/author/affiliation [POST]
func AffiliationNameQueryAuthor(c *gin.Context) {
	name := c.Request.FormValue("affiliation_name")
	name = strings.ToLower(name)
	page, _ := strconv.Atoi(c.Request.FormValue("page"))
	size, _ := strconv.Atoi(c.Request.FormValue("size"))
	sort_type, _ := strconv.Atoi(c.Request.FormValue("sort_type"))
	sort_ascending, _ := strconv.ParseBool(c.Request.FormValue("sort_ascending"))
	affiliation_name := c.Request.FormValue("affiliations")
	affiliations := make([]string, 0)
	err := json.Unmarshal([]byte(affiliation_name), &affiliations)
	if err != nil || (sort_type != 0 && sort_type != 2 && sort_type != 1) {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "机构列表格式错误", "status": 401})
		return
	}

	//affiliationResult := service.QueryByField("affiliation", "name", name, 1, 15)

	boolQuery := elastic.NewBoolQuery()
	query := service.IndexFieldsGetQuery("affiliation", "name", name, 12, "affiliation_id")
	//query := elastic.NewMatchPhraseQuery("name", name)
	boolQuery.Must(query)
	orQuery := elastic.NewBoolQuery()
	for _, affiliation := range affiliations {
		fmt.Println(affiliation)
		orQuery.Should(elastic.NewMatchPhraseQuery("affiliation_id.keyword", affiliation))
	}
	boolQuery.Must(orQuery)
	//affiliationIdQuery := elastic.NewBoolQuery()
	//for _, hit := range affiliationResult.Hits.Hits {
	//	affiliationIdQuery.Should(elastic.NewMatchPhraseQuery("affiliation_id.keyword", hit.Id))
	//}
	//boolQuery.Must(affiliationIdQuery)
	searchResult := service.AuthorQuery(page, size, sort_type, sort_ascending, "author", boolQuery)

	aggregation := make(map[string]interface{})
	aggregation["affiliations"] = service.Paper_Aggregattion(searchResult, "affiliation")
	//aggregation["author"] = service.Paper_Aggregattion(searchResult, "author")
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "作者不存在", "status": 404})
		fmt.Printf("this authors query %s not existed", name)
		return
	}
	fmt.Println("search author", name, "hits :", searchResult.TotalHits())
	ids := make([]string, 0)

	for _, paper := range searchResult.Hits.Hits {
		ids = append(ids, paper.Id)
	}
	authors := service.GetAuthors(ids)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": *service.ParseEnterScholarMsg(&authors), "aggregation": aggregation})
	return
}

// GetAuthorAvatars doc
// @description es 获取一组作者的头像
// @Tags elasticsearch
// @Param author_ids formData string true "作者id组"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Router /es/query/author/avatars [POST]
func GetAuthorAvatars(c *gin.Context) {
	author_ids := c.Request.FormValue("author_ids")
	authorIds := strings.Split(author_ids, `,`)
	pictures := make([]string, 0)
	submits, users := service.QuerySubmitsByAuthor(authorIds)
	for _, authorId := range authorIds {
		find := false
		for _, submit := range submits {
			if authorId == submit.AuthorID {
				for _, user := range users {
					if user.UserID == submit.UserID {
						pictures = append(pictures, user.Avatar)
						find = true
						break
					}
				}
			}
		}
		if !find {
			pictures = append(pictures, utils.PICTURE)
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "data": pictures})
}

// DoiQueryPaper doc
// @description es doi查询论文 精确搜索，结果要么有要么没有
// @Tags elasticsearch
// @Param doi formData string true "doi"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/doi [POST]
func DoiQueryPaper(c *gin.Context) {
	doi := c.Request.FormValue("doi")
	//searchResult := service.MatchPhraseQuery("paper", "doi.keyword", doi, 1, 1)
	searchResult := service.PaperQueryByField("paper", "doi.keyword", doi, 1, 10, true, elastic.NewBoolQuery(), 1, true)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this abstract query %s not existed", doi)
		return
	}
	fmt.Println("search doi", doi, "hits :", searchResult.TotalHits())
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// MainQueryPaper doc
// @description es 根据文章标题 与摘要进行模糊搜索
// @Tags elasticsearch
// @Param main formData string true "main"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 401 {string} string "{"success": false, "message": "参数错误"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Router /es/query/paper/main [POST]
func MainQueryPaper(c *gin.Context) {
	main := c.Request.FormValue("main")
	main = strings.ToLower(main)
	page, _ := strconv.Atoi(c.Request.FormValue("page"))
	size, _ := strconv.Atoi(c.Request.FormValue("size"))
	boolQuery := elastic.NewBoolQuery().Should(elastic.NewMatchPhraseQuery("paper_title", main)).Should(elastic.NewMatchPhraseQuery("abstract", main))
	//searchResult := service.PaperQueryByField("paper")
	doc_type_agg := elastic.NewTermsAggregation().Field("doctype.keyword") // 设置统计字段
	fields_agg := elastic.NewTermsAggregation().Field("fields.keyword")
	conference_agg := elastic.NewTermsAggregation().Field("conference_id.keyword") // 设置统计字段
	journal_id_agg := elastic.NewTermsAggregation().Field("journal_id.keyword")    // 设置统计字段
	publisher_agg := elastic.NewTermsAggregation().Field("publisher.keyword")
	min_year_agg, max_year_agg := elastic.NewMinAggregation().Field("date"), elastic.NewMaxAggregation().Field("date")
	searchResult, err := service.Client.Search().Index("paper").Query(boolQuery).Size(size).TerminateAfter(utils.TERMINATE_AFTER).Aggregation("conference", conference_agg).
		Aggregation("journal", journal_id_agg).Aggregation("doctype", doc_type_agg).Aggregation("fields", fields_agg).Aggregation("publisher", publisher_agg).Aggregation("min_year", min_year_agg).Aggregation("max_year", max_year_agg).
		From((page - 1) * size).Do(context.Background())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "参数错误", "status": 401})
		return
	}
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		return
	}
	fmt.Println("search publisher", main, "hits :", searchResult.TotalHits())

	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// MainSelectPaper doc
// @description es 根据title筛选论文，包括对文章类型journal的筛选，页数的更换,页面大小size的设计, \n 错误码：401 参数格式错误, 排序方式1为默认，2为引用率，3为年份
// @Tags elasticsearch
// @Param main formData string true "main偏官寨关键词"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Param min_year formData int true "min_year"
// @Param max_year formData int true "max_year"
// @Param doctypes formData string true "doctypes"
// @Param conferences formData string true "conferences"
// @Param journals formData string true "journals"
// @Param publishers formData string true "publishers"
// @Param sort_type formData int true "sort_type"
// @Param sort_ascending formData bool true "sort_ascending"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 401 {string} string "{"success": false, "message": "page 不是整数"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/select/paper/main [POST]
func MainSelectPaper(c *gin.Context) {

	var sort_ascending bool
	main := c.Request.FormValue("main")
	main = strings.ToLower(main)
	page_str := c.Request.FormValue("page")
	size_str := c.Request.FormValue("size")
	min_year := c.Request.FormValue("min_year")
	max_year := c.Request.FormValue("max_year")
	doctypesJson, journalsJson, conferenceJson, publisherJson := c.Request.FormValue("doctypes"), c.Request.FormValue("journals"), c.Request.FormValue("conferences"), c.Request.FormValue("publishers")
	doctypes, conferences, journals, publishers := make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100)
	sort_type_str := c.Request.FormValue("sort_type")
	sort_ascending_str := c.Request.FormValue("sort_ascending")

	err := service.CheckSelectPaperParams(c, page_str, size_str, min_year, max_year, doctypesJson, journalsJson, conferenceJson, publisherJson, sort_ascending_str)
	if err != nil {
		// 参数校验401错误
		return
	}
	sort_ascending, _ = strconv.ParseBool(sort_ascending_str)

	page, size, sort_type := service.PureAtoi(page_str), service.PureAtoi(size_str), service.PureAtoi(sort_type_str)
	json.Unmarshal([]byte(doctypesJson), &doctypes)
	json.Unmarshal([]byte(journalsJson), &journals)
	json.Unmarshal([]byte(conferenceJson), &conferences)
	json.Unmarshal([]byte(publisherJson), &publishers)

	boolQuery := service.SelectTypeQuery(doctypes, journals, conferences, publishers, service.PureAtoi(min_year), service.PureAtoi(max_year))
	boolQuery.Must(elastic.NewBoolQuery().Should(elastic.NewMatchPhraseQuery("paper_title", main)).Should(elastic.NewMatchPhraseQuery("abstract", main)))
	searchResult := service.SearchSort(boolQuery, sort_type, sort_ascending, page, size)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this main  query %s not existed", main)
		return
	}
	fmt.Println("search main", main, "hits :", searchResult.TotalHits())
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregationsYear(searchResult)})
	return
}

// AdvancedSearch doc
// @description es 高级搜索
// @Tags elasticsearch
// @Param conditions formData string true "conditions 为条件，表示字典的列表：type 123表示运算符must or，not，"
// @Param min_date formData string true "min_date"
// @Param max_date formData string true "max_date"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 401 {string} string "{"success": false, "message": "参数错误"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/advanced [POST]
func AdvancedSearch(c *gin.Context) {
	//title := c.Request.FormValue("title")

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
	min_date := c.Request.FormValue("min_date")
	max_date := c.Request.FormValue("max_date")
	if min_date == "0" {
		min_date = "1200-01-01 00:00:00"
	} else {
		min_date += " 00:00:00"
	}
	if max_date == "0" {
		max_date = "2050-01-01 00:00:00"
	} else {
		max_date += " 00:00:00"
	}
	minDate, maxDate := service.TimeStrToTimeDefault(min_date), service.TimeStrToTimeDefault(max_date)

	conditionsJson := c.Request.FormValue("conditions")
	var conditions []interface{}
	err = json.Unmarshal([]byte(conditionsJson), &conditions)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "条件表达式格式错误", "status": 401})
		return
	}
	fmt.Println(conditions)
	fmt.Println(minDate, maxDate)

	advancedQuery := service.AdvancedCondition(conditions)
	boolQuery := elastic.NewBoolQuery().Must(advancedQuery)
	boolQuery.Must(elastic.NewRangeQuery("date").Lte(maxDate).Gte(minDate))
	//boolQuery.Must(elastic.NewQuery)
	//boolQuery.Must(elastic.NewMatchPhraseQuery("paper_title", title))
	doc_type_agg := elastic.NewTermsAggregation().Field("doctype.keyword") // 设置统计字段
	fields_agg := elastic.NewTermsAggregation().Field("fields.keyword")
	conference_agg := elastic.NewTermsAggregation().Field("conference_id.keyword") // 设置统计字段
	journal_id_agg := elastic.NewTermsAggregation().Field("journal_id.keyword")    // 设置统计字段
	publisher_agg := elastic.NewTermsAggregation().Field("publisher.keyword")
	min_year_agg, max_year_agg := elastic.NewMinAggregation().Field("date"), elastic.NewMaxAggregation().Field("date")

	searchResult, err := service.Client.Search().Index("paper").Query(boolQuery).TerminateAfter(utils.TERMINATE_AFTER).Aggregation("conference", conference_agg).
		Aggregation("journal", journal_id_agg).Aggregation("doctype", doc_type_agg).Aggregation("fields", fields_agg).Aggregation("publisher", publisher_agg).Aggregation("journal", journal_id_agg).Aggregation("doctype", doc_type_agg).Aggregation("fields", fields_agg).Aggregation("publisher", publisher_agg).Aggregation("min_year", min_year_agg).Aggregation("max_year", max_year_agg).
		Size(size).
		From((page - 1) * size).Do(context.Background())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "条件表达式存在不支持的字段", "status": 401})
		return
	}

	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this advanced query %s not existed")
		return
	}
	fmt.Println("search title", "hits :", searchResult.TotalHits())

	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)
	//paper_authorMap := service.IdsGetItems(paperIds, "paper_author")
	//for i, paper_map_item := range paperSequences {
	//	paper_map_item.(map[string]interface{})["authors"] = service.ParseRelPaperAuthor(paper_authorMap[paperIds[i]].(map[string]interface{}))["rel"]
	//}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// AdvancedSelectPaper doc
// @description es 高级检索筛选论文，包括对文章类型journal的筛选，页数的更换,页面大小size的设计, \n 错误码：401 参数格式错误, 排序方式1为默认，2为引用率，3为年份
// @Tags elasticsearch
// @Param conditions formData string true "conditions 为条件，表示字典的列表：type 123表示运算符must or，not，"
// @Param min_date formData string true "min_date"
// @Param max_date formData string true "max_date"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Param doctypes formData string true "doctypes"
// @Param conferences formData string true "conferences"
// @Param journals formData string true "journals"
// @Param publishers formData string true "publishers"
// @Param sort_type formData int true "sort_type"
// @Param sort_ascending formData bool true "sort_ascending"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 401 {string} string "{"success": false, "message": "page 不是整数"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/select/paper/advanced [POST]
func AdvancedSelectPaper(c *gin.Context) {

	var sort_ascending bool
	page_str := c.Request.FormValue("page")
	size_str := c.Request.FormValue("size")
	min_date := c.Request.FormValue("min_date")
	max_date := c.Request.FormValue("max_date")
	if min_date == "0" {
		min_date = "1200-01-01 00:00:00"
	} else {
		min_date += " 00:00:00"
	}
	if max_date == "0" {
		max_date = "2050-01-01 00:00:00"
	} else {
		max_date += " 00:00:00"
	}
	minDate, maxDate := service.TimeStrToTimeDefault(min_date), service.TimeStrToTimeDefault(max_date)

	conditionsJson := c.Request.FormValue("conditions")
	var conditions []interface{}
	err := json.Unmarshal([]byte(conditionsJson), &conditions)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "条件表达式格式错误", "status": 401})
		return
	}
	doctypesJson, journalsJson, conferenceJson, publisherJson := c.Request.FormValue("doctypes"), c.Request.FormValue("journals"), c.Request.FormValue("conferences"), c.Request.FormValue("publishers")
	doctypes, conferences, journals, publishers := make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100)
	sort_type_str := c.Request.FormValue("sort_type")
	sort_ascending_str := c.Request.FormValue("sort_ascending")

	err = service.CheckSelectPaperParams(c, page_str, size_str, "1", "1", doctypesJson, journalsJson, conferenceJson, publisherJson, sort_ascending_str)
	if err != nil {
		// 参数校验401错误
		return
	}
	sort_ascending, _ = strconv.ParseBool(sort_ascending_str)

	page, size, sort_type := service.PureAtoi(page_str), service.PureAtoi(size_str), service.PureAtoi(sort_type_str)
	json.Unmarshal([]byte(doctypesJson), &doctypes)
	json.Unmarshal([]byte(journalsJson), &journals)
	json.Unmarshal([]byte(conferenceJson), &conferences)
	json.Unmarshal([]byte(publisherJson), &publishers)

	boolQuery := service.SelectTypeQuery(doctypes, journals, conferences, publishers, 0, 2050)
	advancedQuery := service.AdvancedCondition(conditions)
	boolQuery.Must(advancedQuery)
	boolQuery.Must(elastic.NewRangeQuery("date").Lte(maxDate).Gte(minDate))

	fmt.Println(minDate, maxDate)
	//boolQuery.Must(elastic.NewRangeQuery("date").Lte(max_date).Gte(min_date))
	// boolQuery.Filter(elastic.NewRangeQuery("date").From(minDate).To(maxDate))

	searchResult := service.SearchSort(boolQuery, sort_type, sort_ascending, page, size)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this authors name query %s not existed", conditions)
		return
	}
	fmt.Println("search conditions", conditions, "hits :", searchResult.TotalHits())
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregationsYear(searchResult)})
	return
}

// AuthorNameQueryPaper doc
// @description es 根据作者姓名查询文献：精确查询,is_precise=0 为模糊匹配，为1为精准匹配
// @Tags elasticsearch
// @Param author_name formData string true "author_name"
// @Param is_precise formData bool true "is_precise"
// @Success 200 {string} string "{"success": true, "message": "获取作者成功"}"
// @Failure 404 {string} string "{"success": false, "message": "作者不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/author_name [POST]
func AuthorNameQueryPaper(c *gin.Context) {
	author_name := c.Request.FormValue("author_name")
	author_name = strings.ToLower(author_name)
	is_precise, err := strconv.ParseBool(c.Request.FormValue("is_precise"))
	if err != nil {
		panic(err)
	}
	searchResult := service.PaperQueryByField("paper", "authors.aname", author_name, 1, 10, is_precise, elastic.NewBoolQuery(), 1, true)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this author_name query %s not existed", author_name)
		return
	}
	fmt.Println("search author_name", author_name, "hits :", searchResult.TotalHits())

	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// AuthorNameSelectPaper doc
// @description es 根据title筛选论文，包括对文章类型journal的筛选，页数的更换,页面大小size的设计, \n 错误码：401 参数格式错误, 排序方式1为默认，2为引用率，3为年份
// @Tags elasticsearch
// @Param author_name formData string true "author_name"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Param min_year formData int true "min_year"
// @Param max_year formData int true "max_year"
// @Param doctypes formData string true "doctypes"
// @Param conferences formData string true "conferences"
// @Param journals formData string true "journals"
// @Param publishers formData string true "publishers"
// @Param sort_type formData int true "sort_type"
// @Param sort_ascending formData bool true "sort_ascending"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 401 {string} string "{"success": false, "message": "page 不是整数"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/select/paper/author_name [POST]
func AuthorNameSelectPaper(c *gin.Context) {

	var sort_ascending bool
	author_name := c.Request.FormValue("author_name")
	author_name = strings.ToLower(author_name)
	page_str := c.Request.FormValue("page")
	size_str := c.Request.FormValue("size")
	min_year := c.Request.FormValue("min_year")
	max_year := c.Request.FormValue("max_year")
	doctypesJson, journalsJson, conferenceJson, publisherJson := c.Request.FormValue("doctypes"), c.Request.FormValue("journals"), c.Request.FormValue("conferences"), c.Request.FormValue("publishers")
	doctypes, conferences, journals, publishers := make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100)
	sort_type_str := c.Request.FormValue("sort_type")
	sort_ascending_str := c.Request.FormValue("sort_ascending")

	err := service.CheckSelectPaperParams(c, page_str, size_str, min_year, max_year, doctypesJson, journalsJson, conferenceJson, publisherJson, sort_ascending_str)
	if err != nil {
		// 参数校验401错误
		return
	}
	sort_ascending, _ = strconv.ParseBool(sort_ascending_str)

	page, size, sort_type := service.PureAtoi(page_str), service.PureAtoi(size_str), service.PureAtoi(sort_type_str)
	json.Unmarshal([]byte(doctypesJson), &doctypes)
	json.Unmarshal([]byte(journalsJson), &journals)
	json.Unmarshal([]byte(conferenceJson), &conferences)
	json.Unmarshal([]byte(publisherJson), &publishers)

	boolQuery := service.SelectTypeQuery(doctypes, journals, conferences, publishers, service.PureAtoi(min_year), service.PureAtoi(max_year))
	//boolQuery.Must(elastic.NewMatchPhraseQuery("authors.aname", author_name))
	searchResult := service.PaperQueryByField("paper", "authors.aname", author_name, page, size, true, boolQuery, sort_type, sort_ascending)
	//searchResult := service.SearchSort(boolQuery, sort_type, sort_ascending, page, size)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this authors name query %s not existed", author_name)
		return
	}
	fmt.Println("search authors name", author_name, "hits :", searchResult.TotalHits())
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// AffiliationNameQueryPaper doc
// @description es 根据作者姓名查询文献：精确查询,is_precise=0 为模糊匹配，为1为精准匹配
// @Tags elasticsearch
// @Param affiliation_name formData string true "affiliation_name"
// @Param is_precise formData bool true "is_precise"
// @Success 200 {string} string "{"success": true, "message": "获取作者成功"}"
// @Failure 404 {string} string "{"success": false, "message": "作者不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/affiliation_name [POST]
func AffiliationNameQueryPaper(c *gin.Context) {
	affiliation_name := c.Request.FormValue("affiliation_name")
	affiliation_name = strings.ToLower(affiliation_name)
	is_precise, err := strconv.ParseBool(c.Request.FormValue("is_precise"))
	if err != nil {
		panic(err)
	}
	searchResult := service.PaperQueryByField("paper", "authors.afname", affiliation_name, 1, 10, is_precise, elastic.NewBoolQuery(), 1, true)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this affiliation_name query %s not existed", affiliation_name)
		return
	}
	fmt.Println("search author_name", affiliation_name, "hits :", searchResult.TotalHits())
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// AffiliationNameSelectPaper doc
// @description es affiliation_name筛选论文，包括对文章类型journal的筛选，页数的更换,页面大小size的设计, \n 错误码：401 参数格式错误, 排序方式1为默认，2为引用率，3为年份
// @Tags elasticsearch
// @Param affiliation_name formData string true "affiliation_name"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Param min_year formData int true "min_year"
// @Param max_year formData int true "max_year"
// @Param doctypes formData string true "doctypes"
// @Param conferences formData string true "conferences"
// @Param journals formData string true "journals"
// @Param publishers formData string true "publishers"
// @Param sort_type formData int true "sort_type"
// @Param sort_ascending formData bool true "sort_ascending"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 401 {string} string "{"success": false, "message": "page 不是整数"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/select/paper/affiliation_name [POST]
func AffiliationNameSelectPaper(c *gin.Context) {

	var sort_ascending bool
	affiliation_name := c.Request.FormValue("affiliation_name")
	affiliation_name = strings.ToLower(affiliation_name)
	page_str := c.Request.FormValue("page")
	size_str := c.Request.FormValue("size")
	min_year := c.Request.FormValue("min_year")
	max_year := c.Request.FormValue("max_year")
	doctypesJson, journalsJson, conferenceJson, publisherJson := c.Request.FormValue("doctypes"), c.Request.FormValue("journals"), c.Request.FormValue("conferences"), c.Request.FormValue("publishers")
	doctypes, conferences, journals, publishers := make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100)
	sort_type_str := c.Request.FormValue("sort_type")
	sort_ascending_str := c.Request.FormValue("sort_ascending")

	err := service.CheckSelectPaperParams(c, page_str, size_str, min_year, max_year, doctypesJson, journalsJson, conferenceJson, publisherJson, sort_ascending_str)
	if err != nil {
		// 参数校验401错误
		return
	}
	sort_ascending, _ = strconv.ParseBool(sort_ascending_str)

	page, size, sort_type := service.PureAtoi(page_str), service.PureAtoi(size_str), service.PureAtoi(sort_type_str)
	json.Unmarshal([]byte(doctypesJson), &doctypes)
	json.Unmarshal([]byte(journalsJson), &journals)
	json.Unmarshal([]byte(conferenceJson), &conferences)
	json.Unmarshal([]byte(publisherJson), &publishers)

	boolQuery := service.SelectTypeQuery(doctypes, journals, conferences, publishers, service.PureAtoi(min_year), service.PureAtoi(max_year))
	//boolQuery.Must(elastic.NewMatchPhraseQuery("authors.afname", affiliation_name))
	//searchResult := service.SearchSort(boolQuery, sort_type, sort_ascending, page, size)
	searchResult := service.PaperQueryByField("paper", "authors.afname", affiliation_name, page, size, true, boolQuery, sort_type, sort_ascending)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this authors name query %s not existed", affiliation_name)
		return
	}
	fmt.Println("search authors name", affiliation_name, "hits :", searchResult.TotalHits())
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// PublisherQueryPaper doc
// @description es 根据出版商查询文献：精确查询,is_precise=0 为模糊匹配，为1为精准匹配
// @Tags elasticsearch
// @Param publisher formData string true "publisher"
// @Param is_precise formData bool true "is_precise"
// @Success 200 {string} string "{"success": true, "message": "获取作者成功"}"
// @Failure 404 {string} string "{"success": false, "message": "作者不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/publisher [POST]
func PublisherQueryPaper(c *gin.Context) {
	publisher := c.Request.FormValue("publisher")
	publisher = strings.ToLower(publisher)
	is_precise, err := strconv.ParseBool(c.Request.FormValue("is_precise"))
	if err != nil {
		panic(err)
	}
	searchResult := service.PaperQueryByField("paper", "publisher", publisher, 1, 10, is_precise, elastic.NewBoolQuery(), 1, true)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this affiliation_name query %s not existed", publisher)
		return
	}
	fmt.Println("search publisher", publisher, "hits :", searchResult.TotalHits())

	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// PublisherSelectPaper doc
// @description es affiliation_name筛选论文，包括对文章类型journal的筛选，页数的更换,页面大小size的设计, \n 错误码：401 参数格式错误, 排序方式1为默认，2为引用率，3为年份
// @Tags elasticsearch
// @Param publisher formData string true "publisher"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Param min_year formData int true "min_year"
// @Param max_year formData int true "max_year"
// @Param doctypes formData string true "doctypes"
// @Param conferences formData string true "conferences"
// @Param journals formData string true "journals"
// @Param publishers formData string true "publishers"
// @Param sort_type formData int true "sort_type"
// @Param sort_ascending formData bool true "sort_ascending"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 401 {string} string "{"success": false, "message": "page 不是整数"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/select/paper/publisher [POST]
func PublisherSelectPaper(c *gin.Context) {

	var sort_ascending bool
	publisher := c.Request.FormValue("publisher")
	publisher = strings.ToLower(publisher)
	page_str := c.Request.FormValue("page")
	size_str := c.Request.FormValue("size")
	min_year := c.Request.FormValue("min_year")
	max_year := c.Request.FormValue("max_year")
	doctypesJson, journalsJson, conferenceJson, publisherJson := c.Request.FormValue("doctypes"), c.Request.FormValue("journals"), c.Request.FormValue("conferences"), c.Request.FormValue("publishers")
	doctypes, conferences, journals, publishers := make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100)
	sort_type_str := c.Request.FormValue("sort_type")
	sort_ascending_str := c.Request.FormValue("sort_ascending")

	err := service.CheckSelectPaperParams(c, page_str, size_str, min_year, max_year, doctypesJson, journalsJson, conferenceJson, publisherJson, sort_ascending_str)
	if err != nil {
		// 参数校验401错误
		return
	}
	sort_ascending, _ = strconv.ParseBool(sort_ascending_str)

	page, size, sort_type := service.PureAtoi(page_str), service.PureAtoi(size_str), service.PureAtoi(sort_type_str)
	json.Unmarshal([]byte(doctypesJson), &doctypes)
	json.Unmarshal([]byte(journalsJson), &journals)
	json.Unmarshal([]byte(conferenceJson), &conferences)
	json.Unmarshal([]byte(publisherJson), &publishers)

	boolQuery := service.SelectTypeQuery(doctypes, journals, conferences, publishers, service.PureAtoi(min_year), service.PureAtoi(max_year))
	//boolQuery.Must(elastic.NewMatchPhraseQuery("publisher", publisher))
	searchResult := service.PaperQueryByField("paper", "publisher", publisher, page, size, true, boolQuery, sort_type, sort_ascending)
	//searchResult := service.SearchSort(boolQuery, sort_type, sort_ascending, page, size)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this publisher query %s not existed", publisher)
		return
	}
	fmt.Println("search publisher", publisher, "hits :", searchResult.TotalHits())
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// FieldQueryPaper doc
// @description es doi查询论文 精确搜索，结果要么有要么没有
// @Tags elasticsearch
// @Param field formData string true "field"
// @Param page formData int true "page"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/field [POST]
func FieldQueryPaper(c *gin.Context) {
	field := c.Request.FormValue("field")
	field = strings.ToLower(field)
	page, err := strconv.Atoi(c.Request.FormValue("page"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "page 不为整数", "status": 401})
		return
	}

	fieldIds := service.IndexFieldsQueryGetIds("fields", "name", field, 5)
	fmt.Println("search field", field, "hits :", len(fieldIds))
	if len(fieldIds) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this field query %s not existed", field)
		return
	}

	doc_type_agg := elastic.NewTermsAggregation().Field("doctype.keyword") // 设置统计字段
	fields_agg := elastic.NewTermsAggregation().Field("fields.keyword")
	conference_agg := elastic.NewTermsAggregation().Field("conference_id.keyword") // 设置统计字段
	journal_id_agg := elastic.NewTermsAggregation().Field("journal_id.keyword")    // 设置统计字段
	publisher_agg := elastic.NewTermsAggregation().Field("publisher.keyword")
	min_year_agg, max_year_agg := elastic.NewMinAggregation().Field("date"), elastic.NewMaxAggregation().Field("date")

	boolQuery, query := elastic.NewBoolQuery(), elastic.NewBoolQuery()
	for _, hits := range fieldIds {
		boolQuery.Should(elastic.NewMatchPhraseQuery("fields.keyword", hits))
	}
	query.Filter(boolQuery)
	//boolQuery.Filter(elastic.NewRangeQuery("age").Gt("30"))
	searchResult, err := service.Client.Search("paper").Query(query).Size(10).TerminateAfter(utils.TERMINATE_AFTER/10).Aggregation("conference", conference_agg).
		Aggregation("journal", journal_id_agg).Aggregation("doctype", doc_type_agg).Aggregation("fields", fields_agg).Aggregation("publisher", publisher_agg).Aggregation("min_year", min_year_agg).Aggregation("max_year", max_year_agg).
		From((page - 1) * 10).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(searchResult.TotalHits())
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// FieldSelectPaper doc
// @description es doi查询论文 精确搜索，结果要么有要么没有
// @Tags elasticsearch
// @Param field formData string true "field"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Param min_year formData int true "min_year"
// @Param max_year formData int true "max_year"
// @Param doctypes formData string true "doctypes"
// @Param conferences formData string true "conferences"
// @Param journals formData string true "journals"
// @Param publishers formData string true "publishers"
// @Param sort_type formData int true "sort_type"
// @Param sort_ascending formData bool true "sort_ascending"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/field [POST]
func FieldSelectPaper(c *gin.Context) {
	field := c.Request.FormValue("field")
	field = strings.ToLower(field)
	page, err := strconv.Atoi(c.Request.FormValue("page"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "page 不为整数", "status": 401})
		return
	}
	page_str := c.Request.FormValue("page")
	size_str := c.Request.FormValue("size")
	min_year := c.Request.FormValue("min_year")
	max_year := c.Request.FormValue("max_year")
	doctypesJson, journalsJson, conferenceJson, publisherJson := c.Request.FormValue("doctypes"), c.Request.FormValue("journals"), c.Request.FormValue("conferences"), c.Request.FormValue("publishers")
	doctypes, conferences, journals, publishers := make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100)
	sort_type_str := c.Request.FormValue("sort_type")
	sort_ascending_str := c.Request.FormValue("sort_ascending")

	err = service.CheckSelectPaperParams(c, page_str, size_str, min_year, max_year, doctypesJson, journalsJson, conferenceJson, publisherJson, sort_ascending_str)
	if err != nil {
		// 参数校验401错误
		return
	}
	sort_ascending, _ := strconv.ParseBool(sort_ascending_str)

	page, size, sort_type := service.PureAtoi(page_str), service.PureAtoi(size_str), service.PureAtoi(sort_type_str)
	json.Unmarshal([]byte(doctypesJson), &doctypes)
	json.Unmarshal([]byte(journalsJson), &journals)
	json.Unmarshal([]byte(conferenceJson), &conferences)
	json.Unmarshal([]byte(publisherJson), &publishers)

	fieldIds := service.FieldNameGetSimilarIds(field, 5)
	fmt.Println("search field", field, "hits :", len(fieldIds))
	if len(fieldIds) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this field query %s not existed", field)
		return
	}

	boolQuery, query := elastic.NewBoolQuery(), elastic.NewBoolQuery()
	for _, hits := range fieldIds {
		boolQuery.Should(elastic.NewMatchPhraseQuery("fields.keyword", hits))
	}

	//boolQuery.Filter(elastic.NewRangeQuery("age").Gt("30"))
	boolQuery = service.SelectTypeQuery(doctypes, journals, conferences, publishers, service.PureAtoi(min_year), service.PureAtoi(max_year))
	query.Must(boolQuery)
	searchResult := service.SearchSort(query, sort_type, sort_ascending, page, size)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this field query %s not existed", field)
		return
	}
	fmt.Println("search field", field, "hits :", searchResult.TotalHits())
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregationsYear(searchResult)})
	return
}

// AbstractQueryPaper doc
// @description es 根据摘要查询文献：精确查询,is_precise=0 为模糊匹配，为1为精准匹配
// @Tags elasticsearch
// @Param abstract formData string true "abstract"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Param is_precise formData bool true "is_precise"
// @Success 200 {string} string "{"success": true, "message": "获取作者成功"}"
// @Failure 404 {string} string "{"success": false, "message": "作者不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/abstract [POST]
func AbstractQueryPaper(c *gin.Context) {
	abstract := c.Request.FormValue("abstract")
	abstract = strings.ToLower(abstract)
	page, _ := strconv.Atoi(c.Request.FormValue("page"))
	size, _ := strconv.Atoi(c.Request.FormValue("size"))
	is_precise, err := strconv.ParseBool(c.Request.FormValue("is_precise"))
	if err != nil {
		panic(err)
	}
	searchResult := service.PaperQueryByField("paper", "abstract", abstract, page, size, is_precise, elastic.NewBoolQuery(), 1, true)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this affiliation_name query %s not existed", abstract)
		return
	}
	fmt.Println("search abstract", abstract, "hits :", searchResult.TotalHits())

	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// AbstractSelectPaper doc
// @description es abstract筛选论文，包括对文章类型journal的筛选，页数的更换,页面大小size的设计, \n 错误码：401 参数格式错误, 排序方式1为默认，2为引用率，3为年份
// @Tags elasticsearch
// @Param abstract formData string true "abstract"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Param min_year formData int true "min_year"
// @Param max_year formData int true "max_year"
// @Param doctypes formData string true "doctypes"
// @Param conferences formData string true "conferences"
// @Param journals formData string true "journals"
// @Param publishers formData string true "publishers"
// @Param sort_type formData int true "sort_type"
// @Param sort_ascending formData bool true "sort_ascending"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 401 {string} string "{"success": false, "message": "page 不是整数"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/select/paper/abstract [POST]
func AbstractSelectPaper(c *gin.Context) {

	var sort_ascending bool
	abstract := c.Request.FormValue("abstract")
	abstract = strings.ToLower(abstract)
	page_str := c.Request.FormValue("page")
	size_str := c.Request.FormValue("size")
	min_year := c.Request.FormValue("min_year")
	max_year := c.Request.FormValue("max_year")
	doctypesJson, journalsJson, conferenceJson, publisherJson := c.Request.FormValue("doctypes"), c.Request.FormValue("journals"), c.Request.FormValue("conferences"), c.Request.FormValue("publishers")
	doctypes, conferences, journals, publishers := make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100)
	sort_type_str := c.Request.FormValue("sort_type")
	sort_ascending_str := c.Request.FormValue("sort_ascending")

	err := service.CheckSelectPaperParams(c, page_str, size_str, min_year, max_year, doctypesJson, journalsJson, conferenceJson, publisherJson, sort_ascending_str)
	if err != nil {
		// 参数校验401错误
		return
	}
	sort_ascending, _ = strconv.ParseBool(sort_ascending_str)

	page, size, sort_type := service.PureAtoi(page_str), service.PureAtoi(size_str), service.PureAtoi(sort_type_str)
	json.Unmarshal([]byte(doctypesJson), &doctypes)
	json.Unmarshal([]byte(journalsJson), &journals)
	json.Unmarshal([]byte(conferenceJson), &conferences)
	json.Unmarshal([]byte(publisherJson), &publishers)

	boolQuery := service.SelectTypeQuery(doctypes, journals, conferences, publishers, service.PureAtoi(min_year), service.PureAtoi(max_year))

	//boolQuery.Must(elastic.NewMatchPhraseQuery("abstract", abstract))
	searchResult := service.PaperQueryByField("paper", "abstract", abstract, page, size, true, boolQuery, sort_type, sort_ascending)
	//searchResult := service.SearchSort(boolQuery, sort_type, sort_ascending, page, size)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this abstract query %s not existed", abstract)
		return
	}
	fmt.Println("search abstract", abstract, "hits :", searchResult.TotalHits())
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}

// QueryHotPaper doc
// @description es 获取热门文献,根据收藏数判定,返回前10篇文章
// @Tags elasticsearch
// @Success 200 {string} string "{"success": true, "message": "获取文献成功"}"
// @Router /es/query/paper/hot [POST]
func QueryHotPaper(c *gin.Context) {
	collects := service.QueryCollectTop10()
	var paper_ids []string
	for _, collect := range collects {
		paper_ids = append(paper_ids, collect.PaperId)
	}
	paper_detail := service.GetPapers(paper_ids)
	for i, paper := range paper_detail {
		paper.(map[string]interface{})["collect_num"] = collects[i].Num
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "data": paper_detail})
}

// QueryRecommendPaper doc
// @description es 获取推荐文献,根据被引用数量数判定,返回10篇文章
// @Tags elasticsearch
// @Success 200 {string} string "{"success": true, "message": "获取文献成功"}"
// @Router /es/query/paper/recommend [POST]
func QueryRecommendPaper(c *gin.Context) {
	ids := service.GetMost1000CitationPaperIds()
	fmt.Println(ids, len(ids))
	var paper_ids []string
	i := 0
	for true {
		b := rand.Intn(1000)
		if arrays.ContainsString(paper_ids, ids[b]) == -1 {
			paper_ids = append(paper_ids, ids[b])
			i++
			if i >= 10 {
				break
			}
		}
	}
	paper_detail := service.GetPapers(paper_ids)
	for _, paper := range paper_detail {
		collects := service.QueryPaperCollect(paper.(map[string]interface{})["paper_id"].(string))
		paper.(map[string]interface{})["collect_num"] = len(collects)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "data": paper_detail})
}

// GetPaperCitationGraph doc
// @description 获取es期刊详细信息
// @Tags 学者门户
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": 期刊ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /scholar/get/citation/paper [POST]
func GetPaperCitationGraph(c *gin.Context) {
	thisId := c.Request.FormValue("id")
	yearList, citationCountList := service.GetCitationPapersGraph(append(make([]string, 0), thisId))
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "years": yearList, "citations": citationCountList})
	return
}

// GetAuthorCitationGraph doc
// @description 获取es期刊详细信息
// @Tags 学者门户
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": 期刊ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /scholar/get/citation/author [POST]
func GetAuthorCitationGraph(c *gin.Context) {
	thisId := c.Request.FormValue("id")
	yearList, citationCountList := service.GetCitationPapersGraph(service.GetAuthorAllPapersIds(thisId))
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "years": yearList, "citations": citationCountList})
	return
}

// GetPaperCitationPaper doc
// @description 获取es引用论文的论文
// @Tags elasticsearch
// @Param id formData string true "id"
// @Param page formData int true "page"
// @Param size formData int true "size"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": 期刊ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/get/citation/paper [POST]
func GetPaperCitationPaper(c *gin.Context) {
	thisId := c.Request.FormValue("id")
	page, _ := strconv.Atoi(c.Request.FormValue("page"))
	size, _ := strconv.Atoi(c.Request.FormValue("size"))
	citationsIds, total_hits := service.GetPaperCitationIds(append(make([]string, 0), thisId), size, page)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "citations": service.GetPapers(citationsIds), "total_hits": total_hits})
	return
}

// GetRelatedPaper doc
// @description 根据paper—id 获取论文的相关文献
// @Tags elasticsearch
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": 期刊ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/get/related/paper [POST]
func GetRelatedPaper(c *gin.Context) {
	thisId := c.Request.FormValue("id")
	simplePaper := service.GetSimplePaper(thisId)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "related": service.GetPapers(service.GetRelatedPapers(simplePaper["paper_title"].(string)))})
	return
}

// PrefixGetInfo doc
// @description 根据前缀得到搜索建议，返回results 字符串数组
// @Tags elasticsearch
// @Param name formData string true "name 表示字段名"
// @Param content formData string true "content，表示字段的内容"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": 参数错误"}"
// @Failure 404 {string} string "{"success": false, "message": 期刊ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/get/prefix [POST]
func PrefixGetInfo(c *gin.Context) {
	name, content := c.Request.FormValue("name"), c.Request.FormValue("content")
	field, index := "", "paper"
	//content = strings.ToLower(content)
	switch name {
	case "title":
		field = "paper_title"
	case "publisher":
		field = "publisher"
	case "main":
		field = "paper_title"
	case "author_name":
		field = "name"
		index = "author"
	case "affiliation_name":
		field = "name"
		index = "affiliation"
	case "field":
		field = "name"
		index = "fields"
	default:
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "参数不支持前缀搜索", "status": 401})
		return
	}
	prefixResult := service.PrefixSearch(index, field, content, 5)
	results := make([]string, 0)
	for _, hit := range prefixResult.Hits.Hits {
		item := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			panic(err)
		}
		results = append(results, item[field].(string))
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "results": results})
	return
}
