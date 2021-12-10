package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/olivere/elastic/v7"

	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/gin-gonic/gin"
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
	paper := service.GetFullPaper(thisId)
	paper = service.FullPaperSocial(paper)

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

// GetConference doc
// @description 获取es会议详细信息
// @Tags elasticsearch
// @Param id formData string true "id"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": 会议ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/get/conference [POST]
func GetConference(c *gin.Context) {
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
// @Param page formData int true "title"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 401 {string} string "{"success": false, "message": "page 不是整数"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/title [POST]
func TitleQueryPaper(c *gin.Context) {
	//TODO 多表联查，查id的时候同时查询author，  查个屁（父子文档开销太大，扁平化管理了
	title := c.Request.FormValue("title")
	page, err := strconv.Atoi(c.Request.FormValue("page"))

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "page 不为整数", "status": 401})
		return
	}
	searchResult := service.PaperQueryByField("paper", "paper_title", title, page, 10, false)

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

	//其他的aggregation都集成起来了，毕竟每一个都查询都十行代码挺臭的
	aggregation := make(map[string]interface{})

	aggregation["doctype"] = service.Paper_Aggregattion(searchResult, "doctype")
	fmt.Println(aggregation["doctype"])
	aggregation["journal"] = service.Paper_Aggregattion(searchResult, "journal")
	aggregation["conference"] = service.Paper_Aggregattion(searchResult, "conference")
	aggregation["fields"] = service.Paper_Aggregattion(searchResult, "fields")
	aggregation["publisher"] = service.Paper_Aggregattion(searchResult, "publisher")

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences, "aggregation": aggregation})
	return
}

// TitleSelectPaper doc
// @description es 根据title筛选论文，包括对文章类型journal的筛选，页数的更换,页面大小size的设计, \n 错误码：401 参数格式错误, 排序方式1为默认，2为引用率，3为年份
// @Tags elasticsearch
// @Param title formData string true "title"
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
	//TODO 多表联查，查id的时候同时查询author，  查个屁（父子文档开销太大，扁平化管理了
	var searchResult *elastic.SearchResult
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

	err := service.CheckSelectPaperParams(c, page_str, size_str, min_year, max_year, doctypesJson, journalsJson, conferenceJson, publisherJson, sort_ascending_str)
	if err != nil {
		// 参数校验401错误
		return
	}

	if sort_ascending_str == "true" {
		sort_ascending = true
	} else if sort_ascending_str == "false" {
		sort_ascending = false
	}
	page, size, sort_type := service.PureAtoi(page_str), service.PureAtoi(size_str), service.PureAtoi(sort_type_str)
	json.Unmarshal([]byte(doctypesJson), &doctypes)
	json.Unmarshal([]byte(journalsJson), &journals)
	json.Unmarshal([]byte(conferenceJson), &conferences)
	json.Unmarshal([]byte(publisherJson), &publishers)

	boolQuery := service.SelectTypeQuery(doctypes, journals, conferences, publishers, service.PureAtoi(min_year), service.PureAtoi(max_year))
	boolQuery.Must(elastic.NewMatchQuery("paper_title", title))
	if sort_type == 1 {
		searchResult, err = service.Client.Search("paper").Query(boolQuery).Size(size).
			From((page - 1) * size).Do(context.Background())
	} else if sort_type == 2 {
		searchResult, err = service.Client.Search("paper").Query(boolQuery).Size(size).Sort("citation_count", sort_ascending).
			From((page - 1) * size).Do(context.Background())
	} else if sort_type == 3 {
		searchResult, err = service.Client.Search("paper").Query(boolQuery).Size(size).Sort("date", sort_ascending).
			From((page - 1) * size).Do(context.Background())
	}

	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this title query %s not existed", title)
		return
	}
	fmt.Println("search title", title, "hits :", searchResult.TotalHits())
	// TODO 会议与journal信息补全，一次mget替换10此mget
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	paperIds := make([]string, 0, 1000)
	for _, hit := range searchResult.Hits.Hits {
		paperIds = append(paperIds, hit.Id)
	}
	paperSequences = service.GetPapers(paperIds)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences})
	return
}

// NameQueryAuthor doc
// @description es 根据姓名查询作者：精确查询,is_precise=0 为模糊匹配，为1为精准匹配
// @Tags elasticsearch
// @Param name formData string true "name"
// @Param is_precise formData bool flase "is_precise"
// @Success 200 {string} string "{"success": true, "message": "获取作者成功"}"
// @Failure 404 {string} string "{"success": false, "message": "作者不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/author/name [POST]
func NameQueryAuthor(c *gin.Context) {
	name := c.Request.FormValue("name")
	is_precise, err := strconv.ParseBool(c.Request.FormValue("is_precise"))
	if err != nil {
		panic(err)
	}
	boolQuery := elastic.NewBoolQuery()
	if is_precise == true {
		query := elastic.NewMatchPhraseQuery("name.keyword", name)
		boolQuery.Must(query)
	} else {
		query := elastic.NewMatchQuery("name", name)
		boolQuery.Must(query)
	}
	searchResult, err := service.Client.Search().Index("author").Query(boolQuery).From(0).Size(10).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this authors query %s not existed", name)
		return
	}
	fmt.Println("search author", name, "hits :", searchResult.TotalHits())
	var paperSequences []interface{} = make([]interface{}, 0, 1000)
	for _, paper := range searchResult.Hits.Hits {
		paperSequences = append(paperSequences, paper.Source)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paperSequences})
	return
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
	searchResult := service.PaperQueryByField("paper", "doi.keyword", doi, 1, 10, true)
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

//// MainQueryPaper doc
//// @description es 根据文章标题 与摘要进行模糊搜索
//// @Tags elasticsearch
//// @Param main formData string true "main"
//// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
//// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
//// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
//// @Router /es/query/paper/main [POST]
//func MainQueryPaper(c *gin.Context) {
//	main := c.Request.FormValue("main")
//	boolQuery := elastic.NewBoolQuery()
//	queryAbstract := elastic.NewMatchQuery("paperAbstract", main)
//	queryTitle := elastic.NewMatchQuery("title", main)
//	boolQuery.Should(queryAbstract, queryTitle)
//	searchResult, err := service.Client.Search().Index("paper").Query(boolQuery).From(0).Size(10).Do(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	if searchResult.TotalHits() == 0 {
//		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
//		fmt.Printf("this title or abstract query %s not existed", main)
//		return
//	}
//	fmt.Println("search title or abstract", main, "hits :", searchResult.TotalHits())
//	var paperSequences []interface{} = make([]interface{}, 0, 1000)
//	for _, paper := range searchResult.Hits.Hits {
//		paperSequences = append(paperSequences, paper.Source)
//	}
//	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
//		"details": paperSequences})
//	return
//}

// AdvancedSearch doc
// @description es 高级搜索
// @Tags elasticsearch
// @Param musts formData string true "musts"
// @Param nots formData string true "nots"
// @Param ors formData string true "ors 至少是其中之一"
// @Param min_year formData int true "min_year"
// @Param max_year formData int true "max_year"
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
	min_year, err := strconv.Atoi(c.Request.FormValue("min_year"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "min_year 不为整数", "status": 401})
		return
	}
	max_year, err := strconv.Atoi(c.Request.FormValue("max_year"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "max_year 不为整数", "status": 401})
		return
	}

	mustsJson, shouldsJson, notsJson := c.Request.FormValue("musts"), c.Request.FormValue("ors"), c.Request.FormValue("nots")
	musts, nots, shoulds := make(map[string]([]string)), make(map[string]([]string)), make(map[string]([]string))
	err = json.Unmarshal([]byte(mustsJson), &musts)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "musts格式错误", "status": 401})
		return
	}
	err = json.Unmarshal([]byte(shouldsJson), &shoulds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "shoulds 格式错误", "status": 401})
		return
	}
	err = json.Unmarshal([]byte(notsJson), &nots)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "nots格式错误", "status": 401})
		return
	}

	boolQuery := service.AdvancedSearch(min_year, max_year, musts, shoulds, nots)
	//boolQuery.Must(elastic.NewMatchQuery("paper_title", title))
	doc_type_agg := elastic.NewTermsAggregation().Field("doctype.keyword") // 设置统计字段
	fields_agg := elastic.NewTermsAggregation().Field("fields.keyword")
	conference_agg := elastic.NewTermsAggregation().Field("conference_id.keyword") // 设置统计字段
	journal_id_agg := elastic.NewTermsAggregation().Field("journal_id.keyword")    // 设置统计字段
	publisher_agg := elastic.NewTermsAggregation().Field("publisher.keyword")

	searchResult, err := service.Client.Search().Index("paper").Query(boolQuery).Aggregation("conference", conference_agg).
		Aggregation("journal", journal_id_agg).Aggregation("doctype", doc_type_agg).Aggregation("fields", fields_agg).Aggregation("publisher", publisher_agg).
		Size(size).
		From((page - 1) * size).Do(context.Background())

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
	is_precise, err := strconv.ParseBool(c.Request.FormValue("is_precise"))
	if err != nil {
		panic(err)
	}
	searchResult := service.PaperQueryByField("paper", "authors.aname", author_name, 1, 10, is_precise)
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
	is_precise, err := strconv.ParseBool(c.Request.FormValue("is_precise"))
	if err != nil {
		panic(err)
	}
	searchResult := service.PaperQueryByField("paper", "authors.afname", affiliation_name, 1, 10, is_precise)
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
	is_precise, err := strconv.ParseBool(c.Request.FormValue("is_precise"))
	if err != nil {
		panic(err)
	}
	searchResult := service.PaperQueryByField("paper", "publisher", publisher, 1, 10, is_precise)
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
	page, err := strconv.Atoi(c.Request.FormValue("page"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "page 不为整数", "status": 401})
		return
	}

	searchResult := service.PaperQueryByField("fields", "name", field, 1, 5, true)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this field query %s not existed", field)
		return
	}
	fmt.Println("search field", field, "hits :", searchResult.TotalHits())

	doc_type_agg := elastic.NewTermsAggregation().Field("doctype.keyword") // 设置统计字段
	fields_agg := elastic.NewTermsAggregation().Field("fields.keyword")
	conference_agg := elastic.NewTermsAggregation().Field("conference_id.keyword") // 设置统计字段
	journal_id_agg := elastic.NewTermsAggregation().Field("journal_id.keyword")    // 设置统计字段
	publisher_agg := elastic.NewTermsAggregation().Field("publisher.keyword")

	boolQuery := elastic.NewBoolQuery()
	for _, hits := range searchResult.Hits.Hits {
		boolQuery.Should(elastic.NewMatchPhraseQuery("fields.keyword", hits.Id))
	}
	//boolQuery.Filter(elastic.NewRangeQuery("age").Gt("30"))
	searchResult, err = service.Client.Search("paper").Query(boolQuery).Size(10).Aggregation("conference", conference_agg).
		Aggregation("journal", journal_id_agg).Aggregation("doctype", doc_type_agg).Aggregation("fields", fields_agg).Aggregation("publisher", publisher_agg).
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

// AbstractQueryPaper doc
// @description es 根据摘要查询文献：精确查询,is_precise=0 为模糊匹配，为1为精准匹配
// @Tags elasticsearch
// @Param abstract formData string true "abstract"
// @Param is_precise formData bool true "is_precise"
// @Success 200 {string} string "{"success": true, "message": "获取作者成功"}"
// @Failure 404 {string} string "{"success": false, "message": "作者不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/abstract [POST]
func AbstractQueryPaper(c *gin.Context) {
	abstract := c.Request.FormValue("abstract")
	is_precise, err := strconv.ParseBool(c.Request.FormValue("is_precise"))
	if err != nil {
		panic(err)
	}
	searchResult := service.PaperQueryByField("abstract", "abstract", abstract, 1, 10, is_precise)
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
