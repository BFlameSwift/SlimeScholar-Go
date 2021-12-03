package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"gitee.com/online-publish/slime-scholar-go/utils"
	"github.com/olivere/elastic/v7"

	"log"
	"os"
	"strconv"
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
func GetsByIndexId(index string, id string) *elastic.GetResult {
	//通过id查找
	var get1 *elastic.GetResult
	var err error

	get1, err = Client.Get().Index(index).Id(id).Do(context.Background())
	if err != nil {
		panic(err)
	}
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
func PaperQueryByField(index string, field string, content string, page int, size int) *elastic.SearchResult {
	doc_type_agg := elastic.NewTermsAggregation().Field("doctype.keyword") // 设置统计字段
	//TODO 领域
	conference_agg := elastic.NewTermsAggregation().Field("conference_id.keyword") // 设置统计字段
	journal_id_agg := elastic.NewTermsAggregation().Field("journal_id.keyword")    // 设置统计字段
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewMatchQuery(field, content))
	//boolQuery.Filter(elastic.NewRangeQuery("age").Gt("30"))
	searchResult, err := Client.Search(index).Query(boolQuery).Size(size).Aggregation("conference", conference_agg).
		Aggregation("journal", journal_id_agg).Aggregation("doctype", doc_type_agg).
		From((page - 1) * size).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(searchResult.TotalHits())

	return searchResult
}

func MatchPhraseQuery(index string, field string, content string, page int, size int) *elastic.SearchResult {
	query := elastic.NewMatchPhraseQuery(field, content)
	searchResult, err := Client.Search().Index("paper").Query(query).From(0).Size(10).Do(context.Background())
	if err != nil {
		panic(err)
	}
	return searchResult
}

// 通过[]string id—list 来获取结果，其中未命中的结果返回为nil 表示此id文件中不存在
func IdsGetItems(id_list []string, index string) map[string]interface{} {
	mul_item := Client.MultiGet()
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
		panic(err)
	}
	//如果有字段未命中怎么办，可能出现:返回空
	// TODO 调用接口
	var result_map map[string]interface{} = make(map[string]interface{})
	for _, id := range id_list {
		result_map[id] = "null"
	}
	for i, hit := range response.Docs {
		//fmt.Println(hit.Source)
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
		return inter[i].(map[string]interface{})["order"].(string) < inter[j].(map[string]interface{})["order"].(string)
	})
	ret_map["rel"] = inter
	return ret_map
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
	if index == "journal" || index == "conference" || index == "field" {
		for _, bucket := range agg.Buckets {
			if bucket.Key.(string) == "" {
				continue
			}
			result_ids = append(result_ids, bucket.Key.(string))
		}
		result_map = IdsGetItems(result_ids, index)
	}
	for _, bucket := range agg.Buckets {
		m := make(map[string]interface{})
		// 每一个桶都有一个key值，其实就是分组的值，可以理解为SQL的group by值
		if bucket.Key.(string) == "" && bucket_len != 1 {
			continue
		}
		if index == "journal" || index == "conference" || index == "field" {
			m = result_map[bucket.Key.(string)].(map[string]interface{})
			m["count"] = bucket.DocCount
		} else {
			m[bucket.Key.(string)] = bucket.DocCount
		}
		my_list = append(my_list, m)
	}
	return my_list
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
