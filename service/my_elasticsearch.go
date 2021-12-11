package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ESClient *elastic.Client
var Client *elastic.Client
var Timeout = "1s" //超时时间

var host = utils.ELASTIC_SEARCH_HOST //这个是es服务地址,我的是配置到配置文件中了，测试的时候可以写死 比如 http://127.0.0.1:9200

//下面定义的是 聚合时候用的一些参数

func Init() {
	elastic.SetSniff(false) //必须 关闭 Sniffing
	//es 配置
	var err error
	//EsClient.EsCon, err = elastic.NewClient(elastic.SetURL(host))
	Client, err = elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		panic(err)
	}
	info, code, err := Client.Ping(host).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := Client.ElasticsearchVersion(host)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
	ESClient = Client
}

//创建
func Create(Params map[string]string) string {
	//使用字符串
	var res *elastic.IndexResponse
	var err error
	m := make(map[string]interface{})
	//fmt.Println("Creating bodyJson", Params["bodyJson"])
	//fmt.Println([]byte(Params["bodyJson"]))
	err = json.Unmarshal([]byte(Params["bodyJson"]), &m)
	//fmt.Println("m", m)
	res, err = Client.Index().
		Index(Params["index"]).
		Id(Params["id"]).
		BodyJson(m).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	return res.Result
}

//删除
func Delete(Params map[string]string) string {
	var res *elastic.DeleteResponse
	var err error

	res, err = Client.Delete().Index(Params["index"]).
		Type(Params["type"]).
		Id(Params["id"]).
		Do(context.Background())

	if err != nil {
		println(err.Error())
	}

	fmt.Printf("delete result %s\n", res.Result)
	return res.Result
}

//修改
func Update(Params map[string]string) string {
	var res *elastic.IndexResponse
	var err error

	res, err = Client.Index().
		Index(Params["index"]).
		Id(Params["id"]).BodyJson(Params["bodyJson"]).
		Do(context.Background())

	if err != nil {
		panic(err)
	}
	return res.Result

}

//修改
func RealButerrorUpdate(Params map[string]string) string {
	var res *elastic.UpdateResponse
	var err error
	script := elastic.NewScript("ctx._source.retweets += params.num").Param("num", 1)
	res, err = Client.Update().
		Index(Params["index"]).
		Type(Params["type"]).
		Id(Params["id"]).
		Script(script).
		Do(context.Background())

	if err != nil {
		println(err.Error())
		panic(err)
	}
	fmt.Printf("update age %s\n", res.Result)
	return res.Result

}
func GetsByIndexId(index string, id string) (*elastic.GetResult, error) {
	//通过id查找
	var get1 *elastic.GetResult
	var err error

	get1, err = Client.Get().Index(index).Id(id).Do(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	return get1, err
}
func GetsByIndexIdWithout(index string, id string) *elastic.GetResult {
	//通过id查找
	var get1 *elastic.GetResult
	//var err error

	get1, _ = Client.Get().Index(index).Id(id).Do(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	return get1
}

//查找
func Gets(Params map[string]string) (*elastic.GetResult, error) {
	//通过id查找
	var get1 *elastic.GetResult
	var err error
	if len(Params["id"]) < 0 {
		fmt.Printf("param error")
		return get1, errors.New("param error")
	}
	get1, err = Client.Get().Index(Params["index"]).Id(Params["id"]).Do(context.Background())

	return get1, err
}

//搜索
func Query(Params map[string]string) *elastic.SearchResult {
	var res *elastic.SearchResult
	var err error
	//取所有
	res, err = Client.Search(Params["index"]).Type(Params["type"]).Do(context.Background())
	if len(Params["queryString"]) > 0 {
		//字段相等
		q := elastic.NewQueryStringQuery(Params["queryString"])
		res, err = Client.Search(Params["index"]).Type(Params["type"]).Query(q).Do(context.Background())
	}
	if err != nil {
		println(err.Error())
	}

	//if res.Hits.TotalHits > 0 {
	//	fmt.Printf("Found a total of %d Employee \n", res.Hits.TotalHits)
	//}
	return res
}

//简单分页 可用

func List(Params map[string]string) *elastic.SearchResult {
	var res *elastic.SearchResult
	var err error
	size, _ := strconv.Atoi(Params["size"])
	page, _ := strconv.Atoi(Params["page"])
	q := elastic.NewQueryStringQuery(Params["queryString"])

	//排序类型 desc asc es 中只使用 bool 值  true or false
	sort_type := true
	if Params["sort_type"] == "desc" {
		sort_type = false
	}
	//fmt.Printf(" sort info  %s,%s\n", Params["sort"],Params["sort_type"])
	if size < 0 || page < 0 {
		fmt.Printf("param error")
		return res
	}
	if len(Params["queryString"]) > 0 {
		res, err = Client.Search(Params["index"]).
			Type(Params["type"]).
			Query(q).
			Size(size).
			From((page)*size).
			Sort(Params["sort"], sort_type).
			Timeout(Timeout).
			Do(context.Background())

	} else {
		res, err = Client.Search(Params["index"]).
			Type(Params["type"]).
			Size(size).
			From((page)*size).
			Sort(Params["sort"], sort_type).
			//SortBy(elastic.NewFieldSort("add_time").UnmappedType("long").Desc(), elastic.NewScoreSort()).
			Timeout(Timeout).
			Do(context.Background())
	}

	if err != nil {
		println("func list error:" + err.Error())
	}
	return res

}

//聚合 平均 可用
func Aggregation(Params map[string]string) *elastic.SearchResult {
	var res *elastic.SearchResult
	var err error
	//需要聚合的指标 求平均
	avg := elastic.NewAvgAggregation().Field(Params["avg"])
	//单位时间和指定字段
	aggs := elastic.NewDateHistogramAggregation().
		Interval("day").
		Field(Params["field"]).
		//TimeZone("Asia/Shanghai").
		SubAggregation(Params["agg_name"], avg)

	res, err = Client.Search(Params["index"]).
		Type(Params["type"]).
		Size(0).
		Aggregation(Params["aggregation_name"], aggs).
		//Sort(Params["sort"],sort_type).
		Timeout(Timeout).
		Do(context.Background())

	if err != nil {
		println("func Aggregation error:" + err.Error())
	}
	println("func Aggregation here 297")
	return res
}
func GetPaperById(id string) {
	// TODO
}

// 匹配搜索，非完全匹配按照index和字段搜索
func QueryByField(index string, field string, content string, page int, size int) *elastic.SearchResult {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewMatchQuery(field, content))
	//boolQuery.Filter(elastic.NewRangeQuery("age").Gt("30"))
	searchResult, err := Client.Search(index).Query(boolQuery).Size(size).
		From((page - 1) * size).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(searchResult.TotalHits())

	return searchResult
}
func PaperQueryByField(index string, field string, content string, page int, size int, is_precise bool) *elastic.SearchResult {
	doc_type_agg := elastic.NewTermsAggregation().Field("doctype.keyword") // 设置统计字段
	fields_agg := elastic.NewTermsAggregation().Field("fields.keyword")
	conference_agg := elastic.NewTermsAggregation().Field("conference_id.keyword") // 设置统计字段
	journal_id_agg := elastic.NewTermsAggregation().Field("journal_id.keyword")    // 设置统计字段
	publisher_agg := elastic.NewTermsAggregation().Field("publisher.keyword")

	boolQuery := elastic.NewBoolQuery()
	if is_precise == false {
		boolQuery.Must(elastic.NewMatchQuery(field, content))
	} else {
		boolQuery.Must(elastic.NewMatchPhraseQuery(field, content))
	}
	//boolQuery.Filter(elastic.NewRangeQuery("age").Gt("30"))
	searchResult, err := Client.Search(index).Query(boolQuery).Size(size).Aggregation("conference", conference_agg).
		Aggregation("journal", journal_id_agg).Aggregation("doctype", doc_type_agg).Aggregation("fields", fields_agg).Aggregation("publisher", publisher_agg).
		From((page - 1) * size).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(searchResult.TotalHits())

	return searchResult
}

func MatchPhraseQuery(index string, field string, content string, page int, size int) *elastic.SearchResult {
	query := elastic.NewMatchPhraseQuery(field, content)
	searchResult, err := Client.Search().Index("paper").Query(query).From((page - 1) * size).Size(size).Do(context.Background())
	if err != nil {
		panic(err)
	}
	return searchResult
}
func IdsGetList(id_list []string, index string) (retList []interface{}) {
	mul_item := Client.MultiGet()
	fmt.Println("mget : ", index)
	//fmt.Println("len!!!!",len(id_list))
	for _, id := range id_list {
		if len(id) == 0 {
			break
		}
		//res,err := Client.Get().Index(index).Id(id).Do(context.Background())
		q := elastic.NewMultiGetItem().Index(index).Id(id)
		mul_item.Add(q)
	}
	response, err := mul_item.Do(context.Background())
	if err != nil {
		fmt.Println(id_list)
		fmt.Println(index)
		return make([]interface{}, 0)
		//panic(err)
	}
	for _, hit := range response.Docs {
		var m map[string]interface{} = make(map[string]interface{})
		_ = json.Unmarshal([]byte(hit.Source), &m)
		retList = append(retList, m)
	}
	return retList
}

// 通过[]string id—list 来获取结果，其中未命中的结果返回为nil 表示此id文件中不存在
func IdsGetItems(id_list []string, index string) map[string]interface{} {
	mul_item := Client.MultiGet()
	fmt.Println("mget : ", index)
	//fmt.Println("len!!!!",len(id_list))
	for _, id := range id_list {
		if len(id) == 0 {
			break
		}
		//res,err := Client.Get().Index(index).Id(id).Do(context.Background())
		q := elastic.NewMultiGetItem().Index(index).Id(id)
		mul_item.Add(q)
	}
	//response, err := Client.Search().Index(index).Query(elastic.NewIdsQuery().Ids(id_list...)).Size(len(id_list)).Do(context.Background())
	response, err := mul_item.Do(context.Background())
	if err != nil {
		fmt.Println(id_list)
		fmt.Println(index)
		return make(map[string]interface{})
		//panic(err)
	}

	var result_map map[string]interface{} = make(map[string]interface{})
	for _, id := range id_list {
		result_map[id] = ""
	}
	for i, hit := range response.Docs {
		var m map[string]interface{} = make(map[string]interface{})
		_ = json.Unmarshal([]byte(hit.Source), &m)
		result_map[id_list[i]] = m
	}
	//fmt.Println(result_map)
	return result_map
}

func SimplifyPaper(m map[string]interface{}) map[string]interface{} {
	var ret map[string]interface{} = make(map[string]interface{})
	ret["id"], ret["authors"], ret["citation_count"], ret["journalName"], ret["paperAbstract"], ret["reference_count"], ret["year"], ret["title"] = m["id"], m["authors"], m["citation_num"], m["journalName"], m["paperAbstract"], m["reference_num"], m["year"], m["title"]
	return ret
}

func ParseRelPaperAuthor(m map[string]interface{}) map[string]interface{} {
	var inter []interface{} = m["rel"].([]interface{})
	// ret_arr := make([]interface{}, 0, len(inter))
	ret_map := make(map[string]interface{})
	for _, v := range inter {
		v_map := v.(map[string]interface{})
		v_map["author_id"] = v_map["aid"]
		v_map["author_name"] = v_map["aname"]
		v_map["affiliation_id"] = v_map["afid"]
		v_map["affiliation_name"] = v_map["afname"]
		delete(v_map, "aid")
		delete(v_map, "afid")
		delete(v_map, "aname")
		delete(v_map, "afname")
	}
	// 按照作者次序排序
	sort.Slice(inter, func(i, j int) bool {
		if inter[i].(map[string]interface{})["order"] == inter[j].(map[string]interface{})["order"] {
			return inter[i].(map[string]interface{})["author_id"].(string) < inter[j].(map[string]interface{})["author_id"].(string)
		}
		aid1, _ := strconv.Atoi(inter[i].(map[string]interface{})["order"].(string))
		aid2, _ := strconv.Atoi(inter[j].(map[string]interface{})["order"].(string))
		return aid1 < aid2
	})
	ret_map["rel"] = inter
	return ret_map
}
func InterfaceListToStringList(list []interface{}) []string {
	ret_list := make([]string, 0, 1000)
	for _, id := range list {
		ret_list = append(ret_list, id.(string))
	}
	return ret_list
}

func ParseFields(ids []string, index string) []interface{} {
	the_map := IdsGetItems(ids, index)
	ret_list := make([]interface{}, 0, 1000)
	for _, id := range ids {
		ret_list = append(ret_list, the_map[id])
	}
	return ret_list
}
func ComplePaper(paper map[string]interface{}) (paper_map map[string]interface{}) {
	// 补全paper中的作者与领域信息，主要是paper作者可能为空字段

	paper_map = paper
	if paper_map["fields"] != nil {
		paper_map["fields"] = ParseFields(InterfaceListToStringList(paper_map["fields"].([]interface{})), "fields")
	} else {
		paper_map["fields"] = make([]string, 0)
	}

	if paper_map["authors"] != nil {
		authors_map := make(map[string]interface{})
		authors_map["rel"] = paper_map["authors"]
		paper_map["authors"] = (ParseRelPaperAuthor(authors_map))["rel"]
	} else {
		paper_map["authors"] = make([]interface{}, 0, 0)
	}
	return paper_map
}
func PaperGetAuthors(paper_id string) map[string]interface{} {
	var map_param map[string]string = make(map[string]string)
	map_param["index"], map_param["id"] = "paper", paper_id
	map_param["index"] = "paper_author"
	paper_authors, err := Gets(map_param)
	if err != nil {
		panic(err)
	}

	paper_reference_rel_map := make(map[string]interface{})
	err = json.Unmarshal(paper_authors.Source, &paper_reference_rel_map)
	if err != nil {
		panic(err)
	}
	return paper_reference_rel_map
}
func PaperRelMakeMap(str string) []interface{} {
	ret_map := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &ret_map)
	if err != nil {
		panic(err)
	}
	return ret_map["rel"].([]interface{})

}

func Paper_Aggregattion(result *elastic.SearchResult, index string) (my_list []interface{}) {
	agg, found := result.Aggregations.Terms(index)
	if !found {
		log.Fatal("没有找到聚合数据")
	}
	fmt.Println(result.TotalHits())

	// 遍历桶数据
	bucket_len := len(agg.Buckets)
	result_ids := make([]string, 0, 10000)
	result_map := make(map[string]interface{})
	if index == "journal" || index == "conference" || index == "fields" || index == "author" || index == "affiliation" {
		for _, bucket := range agg.Buckets {
			if bucket.Key.(string) == "" {
				continue
			}
			result_ids = append(result_ids, bucket.Key.(string))
		}
		result_map = IdsGetItems(result_ids, index)
	}
	if len(result_map) == 0 && (index == "journal" || index == "conference" || index == "fields") || index == "author" || index == "affiliation" {
		fmt.Println("啥也没聚合到", len(result_ids))
		return make([]interface{}, 0, 0)
	}
	for _, bucket := range agg.Buckets {
		m := make(map[string]interface{})
		// 每一个桶都有一个key值，其实就是分组的值，可以理解为SQL的group by值
		if bucket.Key.(string) == "" && bucket_len != 1 {
			continue
		}
		if index == "journal" || index == "conference" || index == "fields" || index == "author" || index == "affiliation" {
			m = result_map[bucket.Key.(string)].(map[string]interface{})
			m["count"] = bucket.DocCount
			m["id"] = bucket.Key
		} else {
			m[bucket.Key.(string)] = bucket.DocCount
		}
		my_list = append(my_list, m)
	}
	return my_list
}

func SelectTypeQuery(doctypes []string, journals []string, conferences []string, publishers []string, min_year int, max_year int) *elastic.BoolQuery {
	boolQuery := elastic.NewBoolQuery()

	//fmt.Println(len(doctypes))
	if len(doctypes) > 0 {
		doctype_query := elastic.NewBoolQuery()
		for _, doctype := range doctypes {
			doctype_query.Should(elastic.NewMatchQuery("doctype", doctype))
		}
		boolQuery.Must(doctype_query)
	}

	if len(journals) > 0 {
		journal_query := elastic.NewBoolQuery()
		for _, journal := range journals {
			journal_query.Should(elastic.NewTermQuery("journal_id", journal))
		}
		boolQuery.Must(journal_query)
	}
	if len(conferences) > 0 {
		conference_query := elastic.NewBoolQuery()
		for _, conference := range conferences {
			conference_query.Should(elastic.NewTermQuery("conference_id", conference))
		}
		boolQuery.Must(conference_query)
	}
	if len(publishers) > 0 {
		publisher_query := elastic.NewBoolQuery()
		for _, publisher := range publishers {
			publisher_query.Should(elastic.NewMatchPhraseQuery("publisher", publisher))
		}
		boolQuery.Must(publisher_query)
	}
	if min_year > 10 {
		boolQuery.Must(elastic.NewRangeQuery("year").Gte(min_year))
	}
	if max_year < 2022 {
		boolQuery.Must(elastic.NewRangeQuery("year").Lte(max_year))
	} // 尽量优化速度

	return boolQuery
}

func SimplifyAdvanceSearch(must []string, should []string, not []string, field string, boolQuery *elastic.BoolQuery) *elastic.BoolQuery {
	if len(must) > 0 {
		query := elastic.NewBoolQuery()
		for _, item := range must {
			item_list := strings.Split(strings.TrimSpace(item), " ")
			for _, str := range item_list {
				query.Must(elastic.NewMatchQuery(field, str))
			}
		}
		boolQuery.Must(query)
	}
	if len(should) > 0 {
		query := elastic.NewBoolQuery()
		for _, item := range should {
			item_list := strings.Split(strings.TrimSpace(item), " ")
			for _, str := range item_list {
				query.Should(elastic.NewMatchQuery(field, str))
			}
		}
		boolQuery.Should(query)
	}
	if len(not) > 0 {
		query := elastic.NewBoolQuery()
		for _, item := range not {
			item_list := strings.Split(strings.TrimSpace(item), " ")
			for _, str := range item_list {
				query.Should(elastic.NewMatchQuery(field, str))
			}
		}
		boolQuery.MustNot(query)
	}
	return boolQuery
}

func LogicSearch(musts map[string]([]string), shoulds map[string][]string, nots map[string][]string, boolQuery *elastic.BoolQuery) *elastic.BoolQuery {
	var nilStringArray = make([]string, 0)
	boolQuery = SimplifyAdvanceSearch(musts["title"], nilStringArray, nilStringArray, "paper_title", boolQuery)
	boolQuery = SimplifyAdvanceSearch(musts["author"], nilStringArray, nilStringArray, "authors.aname", boolQuery)
	boolQuery = SimplifyAdvanceSearch(musts["field"], nilStringArray, nilStringArray, "field", boolQuery)
	boolQuery = SimplifyAdvanceSearch(musts["doi"], nilStringArray, nilStringArray, "doi.keyword", boolQuery)
	boolQuery = SimplifyAdvanceSearch(musts["author_affiliation"], nilStringArray, nilStringArray, "authors.afname", boolQuery)
	boolQuery = SimplifyAdvanceSearch(musts["source"], nilStringArray, nilStringArray, "publisher", boolQuery)

	mustQuery := boolQuery
	boolQuery = boolQuery.Must(mustQuery)

	boolQuery = SimplifyAdvanceSearch(nilStringArray, shoulds["title"], nilStringArray, "paper_title", boolQuery)
	boolQuery = SimplifyAdvanceSearch(nilStringArray, shoulds["author"], nilStringArray, "authors.aname", boolQuery)
	boolQuery = SimplifyAdvanceSearch(nilStringArray, shoulds["field"], nilStringArray, "field", boolQuery)
	boolQuery = SimplifyAdvanceSearch(nilStringArray, shoulds["doi"], nilStringArray, "doi.keyword", boolQuery)
	boolQuery = SimplifyAdvanceSearch(nilStringArray, shoulds["author_affiliation"], nilStringArray, "authors.afname", boolQuery)
	boolQuery = SimplifyAdvanceSearch(nilStringArray, shoulds["source"], nilStringArray, "publisher", boolQuery)

	orQuery := boolQuery
	boolQuery = boolQuery.Should(orQuery)

	boolQuery = SimplifyAdvanceSearch(nilStringArray, nilStringArray, nots["title"], "paper_title", boolQuery)
	boolQuery = SimplifyAdvanceSearch(nilStringArray, nilStringArray, nots["author"], "authors.aname", boolQuery)
	boolQuery = SimplifyAdvanceSearch(nilStringArray, nilStringArray, nots["field"], "field", boolQuery)
	boolQuery = SimplifyAdvanceSearch(nilStringArray, nilStringArray, nots["doi"], "doi.keyword", boolQuery)
	boolQuery = SimplifyAdvanceSearch(nilStringArray, nilStringArray, nots["author_affiliation"], "authors.afname", boolQuery)
	boolQuery = SimplifyAdvanceSearch(nilStringArray, nilStringArray, nots["source"], "publisher", boolQuery)

	notQuery := boolQuery
	boolQuery = boolQuery.MustNot(notQuery)
	return boolQuery
}

func AdvancedSearch(min_year int, max_year int, musts map[string]([]string), shoulds map[string][]string, nots map[string][]string) *elastic.BoolQuery {
	boolQuery := elastic.NewBoolQuery()
	// 很臭。。有办法但是懒得弄了。。should must的逻辑麻烦
	//fmt.Println(len(doctypes))
	boolQuery = LogicSearch(musts, shoulds, nots, boolQuery)
	if min_year > 10 {
		boolQuery.Must(elastic.NewRangeQuery("year").Gte(min_year))
	}
	if max_year < 2022 {
		boolQuery.Must(elastic.NewRangeQuery("year").Lte(max_year))
	} // 尽量减少筛选优化速度

	return boolQuery
}
func SearchAggregates(searchResult *elastic.SearchResult) map[string]interface{} {
	aggregation := make(map[string]interface{})

	aggregation["doctype"] = Paper_Aggregattion(searchResult, "doctype")
	fmt.Println(aggregation["doctype"])
	aggregation["journal"] = Paper_Aggregattion(searchResult, "journal")
	aggregation["conference"] = Paper_Aggregattion(searchResult, "conference")
	aggregation["fields"] = Paper_Aggregattion(searchResult, "fields")
	aggregation["publisher"] = Paper_Aggregattion(searchResult, "publisher")
	return aggregation
}

// 根据paperids 获取一组完整的paperlist。 加速版，减少多次获取。简化代码
// 从现在开始修正码风！！！go的变量命名用驼峰
// 其中，abstract，field，都不一定有，所以要尽可能保证安全性
func GetPapers(paperIds []string) []interface{} {
	papers := IdsGetList(paperIds, "paper")
	needFieldList := make([]string, 0)
	abstractMap := IdsGetItems(paperIds, "abstract")
	for _, paper := range papers {
		paper := paper.(map[string]interface{}) // 省点事
		if paper["fields"] != nil {
			for _, field := range paper["fields"].([]interface{}) {
				needFieldList = append(needFieldList, field.(string))
				// 可能会冗余几个，但是也不太碍事
			}
		}

	}
	fieldsItems := IdsGetItems(needFieldList, "fields")
	thisFieldList := make([]interface{}, 0)

	for i, paper := range papers {
		paper := paper.(map[string]interface{}) // 省点事
		if paper["fields"] != nil {
			for _, field := range paper["fields"].([]interface{}) {
				thisFieldList = append(thisFieldList, fieldsItems[field.(string)])
			}
		}
		// 格式化authors
		if paper["authors"] != nil {
			authors_map := make(map[string]interface{})
			authors_map["rel"] = paper["authors"]
			paper["authors"] = (ParseRelPaperAuthor(authors_map))["rel"]
		} else {
			paper["authors"] = make([]interface{}, 0)
		}
		abstract := abstractMap[paperIds[i]].(map[string]interface{})["abstract"]
		if abstract != nil {
			paper["abstract"] = abstract
		} else {
			paper["abstract"] = ""
		}
		paper["is_collected"] = false
		paper["fields"] = thisFieldList
		// 格式化paper的fields
		thisFieldList = make([]interface{}, 0)
		papers[i] = paper
	}
	return papers
}

// 获取基本的paper信息
func GetSimplePaper(paper_id string) map[string]interface{} {
	return (GetPapers(append(make([]string, 0), paper_id))[0]).(map[string]interface{})
}

// 获取基本的paper信息
func GetFullPaper(paper_id string) map[string]interface{} {
	paper := GetSimplePaper(paper_id)
	paper["doi_url"] = ""
	if paper["doi"].(string) != "" {
		paper["doi_url"] = "https://dx.doi.org/" + paper["doi"].(string)
	} // 原文链接 100%
	reference_result, err := GetsByIndexId("reference", paper_id)
	if err != nil {
		paper["reference_msg"] = make([]string, 0)
	} else {
		reference_ids_interfaces := PaperRelMakeMap(string(reference_result.Source))
		reference_ids := make([]string, 0, 1000)
		for _, str := range reference_ids_interfaces {
			reference_ids = append(reference_ids, str.(string))
		}
		paper["reference_msg"] = GetPapers(reference_ids)
	}

	paper["citation_msg"] = make([]string, 0)
	paper["journal"] = make(map[string]interface{})
	if paper["journal_id"].(string) != "" {
		paper["journal"] = GetsByIndexIdWithout("journal", paper["journal_id"].(string)).Source
	}
	paper["conference"] = make(map[string]interface{})
	if paper["conference_id"].(string) != "" {
		paper["conference"] = GetsByIndexIdWithout("conference", paper["conference_id"].(string)).Source
	}
	return paper
}

// 补充Paper的社交属性
func FullPaperSocial(paper map[string]interface{}) map[string]interface{} {
	paperId := paper["paper_id"].(string)
	// 收集数目
	paper["collect_count"] = len(PaperGetCollectedUsers(paperId))

	return paper
}

func CheckSelectPaperParams(c *gin.Context, page_str string, size_str string, minYear string, maxYear string, doctypesJson string, journalsJson string, conferenceJson string, publisherJson string, sort_ascending_str string) error {
	_, err := strconv.Atoi(page_str)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "page 不为整数", "status": 401})
		return err
	}
	_, err = strconv.Atoi(size_str)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "size 不为整数", "status": 401})
		return err
	}
	_, err = strconv.Atoi(minYear)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "min_year 不为整数", "status": 401})
		return err
	}
	_, err = strconv.Atoi(maxYear)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "max_year 不为整数", "status": 401})
		return err
	}

	doctypes, conferences, journals, publishers := make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100), make([]string, 0, 100)

	//sort_type, _ := strconv.Atoi(c.Request.FormValue("sort_type"))

	if sort_ascending_str == "true" {

	} else if sort_ascending_str == "false" {

	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "sort_ascending 不是truefalse", "status": 401})
		return err
	}
	err = json.Unmarshal([]byte(doctypesJson), &doctypes)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "doctypes格式错误", "status": 401})
		return err
	}
	err = json.Unmarshal([]byte(journalsJson), &journals)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "journals格式错误", "status": 401})
		return err
	}
	err = json.Unmarshal([]byte(conferenceJson), &conferences)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "conferneces格式错误", "status": 401})
		return err
	}
	err = json.Unmarshal([]byte(publisherJson), &publishers)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "publisher格式错误", "status": 401})
		return err
	}

	return nil
}

func SearchSort(boolQuery *elastic.BoolQuery, sort_type int, sort_ascending bool, page int, size int) *elastic.SearchResult {
	var searchResult *elastic.SearchResult

	if sort_type == 1 {
		searchResult, _ = Client.Search("paper").Query(boolQuery).Size(size).
			From((page - 1) * size).Do(context.Background())
	} else if sort_type == 2 {
		searchResult, _ = Client.Search("paper").Query(boolQuery).Size(size).Sort("citation_count", sort_ascending).
			From((page - 1) * size).Do(context.Background())
	} else if sort_type == 3 {
		searchResult, _ = Client.Search("paper").Query(boolQuery).Size(size).Sort("date", sort_ascending).
			From((page - 1) * size).Do(context.Background())
	}
	return searchResult
}

func parseCondition(condition map[string]interface{}) elastic.Query {
	theMap := condition
	delete(theMap, "type")
	key := GetMapAllKey(theMap)[0]
	switch key {
	case "source":
		return elastic.NewMatchQuery("publisher", condition[key])
	case "title":
		return elastic.NewMatchPhraseQuery("paper_title", condition[key])
	case "author":
		return elastic.NewMatchPhraseQuery("authors.aname", condition[key])
	case "doi":
		return elastic.NewTermQuery("doi.keyword", condition[key])
	case "author_affiliation":
		return elastic.NewMatchPhraseQuery("authors.afname", condition[key])

	}
	return nil
}

// 高级检索条件设置
func AdvancedCondition(conditions []interface{}) *elastic.BoolQuery {
	boolQuery := elastic.NewBoolQuery()
	var condition int
	orQuery := elastic.NewBoolQuery().Must(parseCondition(conditions[0].(map[string]interface{})))
	for i := 1; i < len(conditions); i++ {
		condition = int((conditions[i]).(map[string]interface{})["type"].(float64))
		if condition == 3 {
			boolQuery.MustNot(parseCondition(conditions[i].(map[string]interface{})))
		} else if condition == 2 {
			boolQuery.Should(orQuery)
			orQuery = elastic.NewBoolQuery()
			orQuery.Must(parseCondition(conditions[i].(map[string]interface{})))
		} else if condition == 1 {
			orQuery.Must(parseCondition(conditions[i].(map[string]interface{})))
		}
	}
	boolQuery.Should(orQuery)
	return boolQuery
}

// 搜索作者返回结果
func AuthorQuery(page int, size int, sort_type int, sort_ascending bool, index string, boolQuery *elastic.BoolQuery) (searchResult *elastic.SearchResult) {
	//authorNameAgg := elastic.NewTermsAggregation().Field("name.keyword") // 设置统计字段
	affiliationNameAgg := elastic.NewTermsAggregation().Field("affiliation_id.keyword")
	if sort_type == 0 {
		searchResult, err := Client.Search().Index(index).Query(boolQuery).Aggregation("affiliation", affiliationNameAgg).From((page - 1) * size).Size(size).Do(context.Background())
		if err != nil {
			panic(err)
		}
		return searchResult
	} else if sort_type == 1 {
		searchResult, err := Client.Search().Index(index).Query(boolQuery).Aggregation("affiliation", affiliationNameAgg).From((page-1)*size).Size(size).Sort("paper_count", sort_ascending).Do(context.Background())
		if err != nil {
			panic(err)
		}
		return searchResult
	} else if sort_type == 2 {
		searchResult, err := Client.Search().Index(index).Query(boolQuery).Aggregation("affiliation", affiliationNameAgg).From((page-1)*size).Size(size).Sort("citation_count", sort_ascending).Do(context.Background())
		if err != nil {
			panic(err)
		}
		return searchResult
	}
	return nil
}

// func main() {
// 	Init()
// 	fmt.Println("123")
// 	var map_param map[string]string = make(map[string]string)
// 	e1, _ := json.Marshal(model.ValueString{Value: "132"})

// 	map_param["index"], map_param["type"], map_param["id"], map_param["bodyJson"] = "megacorp", "employee", "53", string(e1)
// 	// ret := Create(map_param)
// 	// fmt.Printf(ret)
// 	get_ret, _ := Gets(map_param)
// 	fmt.Printf(get_ret.Id)

// }
