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
	this_id := c.Request.FormValue("id")
	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["id"] = "paper", this_id
	_, error_get := service.Gets(map_param)
	if error_get != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在", "status": 404})
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在", "status": 404})
		fmt.Printf("this id %s not existed", this_id)
		return
	}
	ret, _ := service.Gets(map_param)
	body_byte, _ := json.Marshal(ret.Source)
	var paper = make(map[string]interface{})
	_ = json.Unmarshal(body_byte, &paper)
	// 查找信息
	paper["journal"] = make(map[string]interface{})
	if paper["journal_id"].(string) != "" {
		paper["journal"] = service.GetsByIndexIdWithout("journal", paper["journal_id"].(string)).Source
	}
	paper["conference"] = make(map[string]interface{})
	if paper["conference_id"].(string) != "" {
		paper["conference"] = service.GetsByIndexIdWithout("conference", paper["conference_id"].(string)).Source
	}
	//paper["authors"] = service.ParseRelPaperAuthor(service.PaperGetAuthors(this_id))["rel"]
	paper["abstract"] = service.SemanticScholarApiSingle(this_id, "abstract")
	paper["doi_url"] = ""
	if paper["doi"].(string) != "" {
		paper["doi_url"] = "https://dx.doi.org/" + paper["doi"].(string)
	} // 原文链接 100%
	reference_result, err := service.GetsByIndexId("reference", this_id)
	if err != nil {
		paper["reference_msg"] = make([]string, 0)
	} else {
		reference_ids_interfaces := service.PaperRelMakeMap(string(reference_result.Source))
		reference_ids := make([]string, 0, 1000)
		for _, str := range reference_ids_interfaces {
			reference_ids = append(reference_ids, str.(string))
		}
		paper["reference_msg"] = (service.GetMapAllContent(service.IdsGetItems(reference_ids, "paper")))
	}

	paper["citation_msg"] = make([]string, 0)
	if paper["fields"] != nil {
		paper["fields"] = service.ParseFields(service.InterfaceListToStringList(paper["fields"].([]interface{})), "fields")
	} else {
		paper["fields"] = make([]string, 0)
	}
	//paper["fields"] = make([]string, 0)
	//service.BrowerPaper(paper)
	// id_inter_list := paper["outCitations"].([]interface{})
	// var id_list []string = make([]string, 0, 3000)
	// for _, id := range id_inter_list {
	// 	id_list = append(id_list, id.(string))
	// }
	// fmt.Println(id_list)
	// reference_map := service.IdsGetPapers(id_list, "paper")
	// reference_list := make([]interface{}, 0, 3000)
	// TOOD 最大饮用量可能不一样
	// for _, id := range id_list {
	// 	var item interface{} = reference_map[id]
	// 	//fmt.Println(item)
	// 	reference_list = append(reference_list, service.SimplifyPaper(item.(map[string]interface{})))
	// }
	// paper["reference_msg"] = reference_list
	//fmt.Println(paper)
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
	this_id := c.Request.FormValue("id")
	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["id"], map_param["bodyJson"] = "author", this_id, ""

	_, error_get := service.Gets(map_param)
	if error_get != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在", "status": 404})
		fmt.Println("this id %s not existed", this_id)
		return
	}
	ret, _ := service.Gets(map_param)
	var author_map map[string]interface{} = make(map[string]interface{})
	body_byte, _ := json.Marshal(ret.Source)
	err := json.Unmarshal(body_byte, &author_map)
	if err != nil {
		panic(err)
	}
	fmt.Println(author_map)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "details": author_map, "body": string(body_byte)})
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
	this_id := c.Request.FormValue("id")
	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["id"], map_param["bodyJson"] = "conference", this_id, ""

	_, error_get := service.Gets(map_param)
	if error_get != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在", "status": 404})
		fmt.Println("this id %s not existed", this_id)
		return
	}
	ret, _ := service.Gets(map_param)
	var conference_map map[string]interface{} = make(map[string]interface{})
	body_byte, _ := json.Marshal(ret.Source)
	err := json.Unmarshal(body_byte, &conference_map)
	if err != nil {
		panic(err)
	}
	fmt.Println(conference_map)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "details": conference_map})
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
	this_id := c.Request.FormValue("id")
	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["id"], map_param["bodyJson"] = "journal", this_id, ""

	_, error_get := service.Gets(map_param)
	if error_get != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "索引不存在", "status": 404})
		fmt.Println("this id %s not existed", this_id)
		return
	}
	ret, _ := service.Gets(map_param)
	var journal_map map[string]interface{} = make(map[string]interface{})
	body_byte, _ := json.Marshal(ret.Source)
	err := json.Unmarshal(body_byte, &journal_map)
	if err != nil {
		panic(err)
	}
	fmt.Println(journal_map)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "details": journal_map})
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

	var paper_sequences []interface{} = make([]interface{}, 0, 1000)
	paper_ids := make([]string, 0, 1000)
	for _, paper := range searchResult.Hits.Hits {
		body_byte, _ := json.Marshal(paper.Source)
		var paper_map = make(map[string]interface{})
		_ = json.Unmarshal(body_byte, &paper_map)
		paper_ids = append(paper_ids, paper_map["paper_id"].(string))
		paper_map = service.ComplePaper(paper_map)

		paper_sequences = append(paper_sequences, paper_map)
	}

	//其他的aggregation都集成起来了，毕竟每一个都查询都十行代码挺臭的
	aggregation := make(map[string]interface{})

	aggregation["doctype"] = service.Paper_Aggregattion(searchResult, "doctype")
	fmt.Println(aggregation["doctype"])
	aggregation["journal"] = service.Paper_Aggregattion(searchResult, "journal")
	aggregation["conference"] = service.Paper_Aggregattion(searchResult, "conference")
	aggregation["fields"] = service.Paper_Aggregattion(searchResult, "fields")
	aggregation["publisher"] = service.Paper_Aggregattion(searchResult, "publisher")

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paper_sequences, "aggregation": aggregation})
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
	title := c.Request.FormValue("title")
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

	doctypes_json, journals_json, conference_json, publisher_json := c.Request.FormValue("doctypes"), c.Request.FormValue("journals"), c.Request.FormValue("conferences"), c.Request.FormValue("publishers")
	doctypes, conferences, journals, publishers := make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100)
	var searchResult *elastic.SearchResult
	var sort_ascending bool

	sort_type, _ := strconv.Atoi(c.Request.FormValue("sort_type"))
	sort_ascending_str := c.Request.FormValue("sort_ascending")
	if sort_ascending_str == "true" {
		sort_ascending = true
	} else if sort_ascending_str == "false" {
		sort_ascending = false
	} else {
		sort_ascending = true
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "sort_ascending 不是truefalse", "status": 401})
		return
	}

	err = json.Unmarshal([]byte(doctypes_json), &doctypes)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "doctypes格式错误", "status": 401})
		return
	}
	err = json.Unmarshal([]byte(journals_json), &journals)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "journals格式错误", "status": 401})
		return
	}
	err = json.Unmarshal([]byte(conference_json), &conferences)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "conferneces格式错误", "status": 401})
		return
	}
	err = json.Unmarshal([]byte(publisher_json), &publishers)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "publisher格式错误", "status": 401})
		return
	}
	fmt.Println(doctypes, journals, conferences, publishers)
	boolQuery := service.SelectTypeQuery(doctypes, journals, conferences, publishers, min_year, max_year)
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
	var paper_sequences []interface{} = make([]interface{}, 0, 1000)
	paper_ids := make([]string, 0, 1000)
	for _, paper := range searchResult.Hits.Hits {
		body_byte, _ := json.Marshal(paper.Source)
		var paper_map = make(map[string]interface{})
		_ = json.Unmarshal(body_byte, &paper_map)
		paper_map = service.ComplePaper(paper_map)

		paper_ids = append(paper_ids, paper_map["paper_id"].(string))
		paper_sequences = append(paper_sequences, paper_map)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paper_sequences})
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
	var paper_sequences []interface{} = make([]interface{}, 0, 1000)
	for _, paper := range searchResult.Hits.Hits {
		paper_sequences = append(paper_sequences, paper.Source)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paper_sequences})
	return
}

// // AbstractQueryPaper doc
// // @description es 根据abstract查询论文
// // @Tags elasticsearch
// // @Param paperAbstract formData string true "paperAbstract"
// // @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// // @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// // @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// // @Router /es/query/paper/abstract [POST]
// func AbstractQueryPaper(c *gin.Context) {
// 	abstract := c.Request.FormValue("paperAbstract")
// 	searchResult := service.QueryByField("paper", "paperAbstract", abstract, 1, 10)

// 	if searchResult.TotalHits() == 0 {
// 		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
// 		fmt.Printf("this abstract query %s not existed", abstract)
// 		return
// 	}
// 	fmt.Println("search abstract", abstract, "hits :", searchResult.TotalHits())
// 	var paper_sequences []interface{} = make([]interface{}, 0, 1000)
// 	for _, paper := range searchResult.Hits.Hits {
// 		paper_sequences = append(paper_sequences, paper.Source)
// 	}
// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
// 		"details": paper_sequences})
// 	return
// }

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

	searchResult := service.MatchPhraseQuery("paper", "doi.keyword", doi, 1, 1)
	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this abstract query %s not existed", doi)
		return
	}
	fmt.Println("search doi", doi, "hits :", searchResult.TotalHits())
	var paper_sequences []interface{} = make([]interface{}, 0, 1000)
	for _, paper := range searchResult.Hits.Hits {
		paper_sequences = append(paper_sequences, paper.Source)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paper_sequences, "aggregation": service.SearchAggregates(searchResult)})
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
//	var paper_sequences []interface{} = make([]interface{}, 0, 1000)
//	for _, paper := range searchResult.Hits.Hits {
//		paper_sequences = append(paper_sequences, paper.Source)
//	}
//	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
//		"details": paper_sequences})
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

	musts_json, shoulds_json, nots_json := c.Request.FormValue("musts"), c.Request.FormValue("ors"), c.Request.FormValue("nots")
	musts, nots, shoulds := make(map[string]([]string)), make(map[string]([]string)), make(map[string]([]string))
	err = json.Unmarshal([]byte(musts_json), &musts)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "musts格式错误", "status": 401})
		return
	}
	err = json.Unmarshal([]byte(shoulds_json), &shoulds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "shoulds 格式错误", "status": 401})
		return
	}
	err = json.Unmarshal([]byte(nots_json), &nots)
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

	var paper_sequences []interface{} = make([]interface{}, 0, 1000)
	paper_ids := make([]string, 0, 1000)
	for _, paper := range searchResult.Hits.Hits {
		body_byte, _ := json.Marshal(paper.Source)
		var paper_map = make(map[string]interface{})
		_ = json.Unmarshal(body_byte, &paper_map)
		paper_map = service.ComplePaper(paper_map)

		paper_ids = append(paper_ids, paper_map["paper_id"].(string))
		paper_sequences = append(paper_sequences, paper_map)
	}
	paper_author_map := service.IdsGetItems(paper_ids, "paper_author")
	for i, paper_map_item := range paper_sequences {
		paper_map_item.(map[string]interface{})["authors"] = service.ParseRelPaperAuthor(paper_author_map[paper_ids[i]].(map[string]interface{}))["rel"]
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paper_sequences, "aggregation": service.SearchAggregates(searchResult)})
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

	var paper_sequences []interface{} = make([]interface{}, 0, 1000)
	paper_ids := make([]string, 0, 1000)
	for _, paper := range searchResult.Hits.Hits {
		body_byte, _ := json.Marshal(paper.Source)
		var paper_map = make(map[string]interface{})
		_ = json.Unmarshal(body_byte, &paper_map)
		paper_ids = append(paper_ids, paper_map["paper_id"].(string))
		paper_map = service.ComplePaper(paper_map)

		paper_sequences = append(paper_sequences, paper_map)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paper_sequences, "aggregation": service.SearchAggregates(searchResult)})
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

	var paper_sequences []interface{} = make([]interface{}, 0, 1000)
	paper_ids := make([]string, 0, 1000)
	for _, paper := range searchResult.Hits.Hits {
		body_byte, _ := json.Marshal(paper.Source)
		var paper_map = make(map[string]interface{})
		_ = json.Unmarshal(body_byte, &paper_map)
		paper_ids = append(paper_ids, paper_map["paper_id"].(string))
		paper_map = service.ComplePaper(paper_map)

		paper_sequences = append(paper_sequences, paper_map)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paper_sequences, "aggregation": service.SearchAggregates(searchResult)})
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

	var paper_sequences []interface{} = make([]interface{}, 0, 1000)
	paper_ids := make([]string, 0, 1000)
	for _, paper := range searchResult.Hits.Hits {
		body_byte, _ := json.Marshal(paper.Source)
		var paper_map = make(map[string]interface{})
		_ = json.Unmarshal(body_byte, &paper_map)
		paper_ids = append(paper_ids, paper_map["paper_id"].(string))
		paper_map = service.ComplePaper(paper_map)

		paper_sequences = append(paper_sequences, paper_map)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paper_sequences, "aggregation": service.SearchAggregates(searchResult)})
	return
}
