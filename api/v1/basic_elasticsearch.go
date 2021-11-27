package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
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
		fmt.Println("this id %s not existed", this_id)
		return
	}
	ret, _ := service.Gets(map_param)
	body_byte, _ := json.Marshal(ret.Source)
	var paper = make(map[string]interface{})
	_ = json.Unmarshal(body_byte, &paper)
	// 查找信息
	map_param["index"] = "paper_author"
	// reference_msg := make(map[string]interface{})

	paper_authors, err := service.Gets(map_param)
	if err != nil {
		panic(err)
	}

	paper_reference_rel_map := make(map[string]interface{})
	err = json.Unmarshal(paper_authors.Source, &paper_reference_rel_map)
	if err != nil {
		panic(err)
	}
	paper["authors"] = service.ParseRelPaperAuthor(paper_reference_rel_map)["rel"]
	paper["abstract"] = service.SemanticScholarApiSingle(this_id, "abstract")

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

// TitleQueryPaper doc
// @description es 根据title查询论文
// @Tags elasticsearch
// @Param title formData string true "title"
// @Success 200 {string} string "{"success": true, "message": "获取成功"}"
// @Failure 404 {string} string "{"success": false, "message": "论文不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/paper/title [POST]
func TitleQueryPaper(c *gin.Context) {
	title := c.Request.FormValue("title")
	searchResult := service.QueryByField("paper", "paper_title", title, 1, 10)

	if searchResult.TotalHits() == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "论文不存在", "status": 404})
		fmt.Printf("this title query %s not existed", title)
		return
	}
	fmt.Println("search title", title, "hits :", searchResult.TotalHits())
	var paper_sequences []interface{} = make([]interface{}, 0, 1000)
	for _, paper := range searchResult.Hits.Hits {
		paper_sequences = append(paper_sequences, paper.Source)
		//paper_sequences[strconv.FormatInt(int64(i), 10)] = paper.Source
	}
	//body_byte,_ := json.Marshal(ret.Source)
	//var paper = make(map[string]interface{})
	//_ = json.Unmarshal(body_byte,&paper)
	//fmt.Println(paper)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查找成功", "status": 200, "total_hits": searchResult.TotalHits(),
		"details": paper_sequences})
	return
}

// NameQueryAuthor doc
// @description es 根据姓名查询作者：精确查询,isPrecise=0 为模糊匹配，为1为精准匹配
// @Tags elasticsearch
// @Param name formData string true "name"
// @Param isPrecise formData int flase "isPrecise"
// @Success 200 {string} string "{"success": true, "message": "获取作者成功"}"
// @Failure 404 {string} string "{"success": false, "message": "作者不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "错误500"}"
// @Router /es/query/author/name [POST]
func NameQueryAuthor(c *gin.Context) {
	name := c.Request.FormValue("name")
	isPrecise, err := strconv.Atoi(c.Request.FormValue("isPrecise"))
	if err != nil {
		panic(err)
	}
	boolQuery := elastic.NewBoolQuery()
	if isPrecise == 1 {
		query := elastic.NewMatchPhraseQuery("authors.name", name)
		boolQuery.Must(query)
	} else {
		query := elastic.NewMatchQuery("authors.name", name)
		boolQuery.Must(query)
	}
	searchResult, err := service.Client.Search().Index("paper").Query(boolQuery).From(0).Size(10).Do(context.Background())
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

	searchResult := service.MatchPhraseQuery("paper", "doi", doi, 1, 10)
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
		"details": paper_sequences})
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
