package service

import (
	"context"
	"encoding/json"
	"fmt"

	"gitee.com/online-publish/slime-scholar-go/utils"
	"github.com/olivere/elastic"

	"log"
	"os"
	"strconv"
	"time"
)

type EsClientType struct {
	EsCon *elastic.Client
}

var client *elastic.Client
var Timeout = "1s"        //超时时间
var EsClient EsClientType //连接类型

var host = utils.ELASTIC_SEARCH_HOST //这个是es服务地址,我的是配置到配置文件中了，测试的时候可以写死 比如 http://127.0.0.1:9200

//下面定义的是 聚合时候用的一些参数

type Aggregations struct {
	AVG_Metric AVG_Metric `json:"AVG_Metric"`
}

type AVG_Metric struct {
	Buckets []Metric `json:"buckets"`
}

type Metric struct {
	Avg_time Value `json:"avg_time"`
}

type Value struct {
	Value float64 `json:"value"`
}

func Init() {
	elastic.SetSniff(false) //必须 关闭 Sniffing
	//es 配置
	var err error
	//EsClient.EsCon, err = elastic.NewClient(elastic.SetURL(host))
	client, err = elastic.NewClient(
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
	info, code, err := client.Ping(host).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := client.ElasticsearchVersion(host)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
	fmt.Println("conn es succ", EsClient.EsCon)
}

//创建
func Create(Params map[string]string) string {
	//使用字符串
	var res *elastic.IndexResponse
	var err error

	res, err = client.Index().
		Index(Params["index"]).
		Type(Params["type"]).
		Id(Params["id"]).BodyJson(Params["bodyJson"]).
		Do(context.Background())

	if err != nil {
		panic(err)
	}
	return res.Result
}

//删除
func (client *EsClientType) Delete(Params map[string]string) string {
	var res *elastic.DeleteResponse
	var err error

	res, err = client.EsCon.Delete().Index(Params["index"]).
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
func (client *EsClientType) Update(Params map[string]string) string {
	var res *elastic.UpdateResponse
	var err error

	res, err = client.EsCon.Update().
		Index(Params["index"]).
		Type(Params["type"]).
		Id(Params["id"]).
		Doc(Params["doc"]).
		Do(context.Background())

	if err != nil {
		println(err.Error())
	}
	fmt.Printf("update age %s\n", res.Result)
	return res.Result

}

//查找
func Gets(Params map[string]string) *elastic.GetResult {
	//通过id查找
	var get1 *elastic.GetResult
	var err error
	if len(Params["id"]) < 0 {
		fmt.Printf("param error")
		return get1
	}

	get1, err = client.Get().Index(Params["index"]).Type(Params["type"]).Id(Params["id"]).Do(context.Background())

	if err != nil {
		panic(err)
	}

	return get1
}

//搜索
func (client EsClientType) Query(Params map[string]string) *elastic.SearchResult {
	var res *elastic.SearchResult
	var err error
	//取所有
	res, err = client.EsCon.Search(Params["index"]).Type(Params["type"]).Do(context.Background())
	if len(Params["queryString"]) > 0 {
		//字段相等
		q := elastic.NewQueryStringQuery(Params["queryString"])
		res, err = client.EsCon.Search(Params["index"]).Type(Params["type"]).Query(q).Do(context.Background())
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

func (client *EsClientType) List(Params map[string]string) *elastic.SearchResult {
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
		res, err = client.EsCon.Search(Params["index"]).
			Type(Params["type"]).
			Query(q).
			Size(size).
			From((page)*size).
			Sort(Params["sort"], sort_type).
			Timeout(Timeout).
			Do(context.Background())

	} else {
		res, err = client.EsCon.Search(Params["index"]).
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
func (client *EsClientType) Aggregation(Params map[string]string) *elastic.SearchResult {
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

	res, err = client.EsCon.Search(Params["index"]).
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

func main() {
	Init()
	fmt.Println("123")
	var map_param map[string]string = make(map[string]string)
	e1, _ := json.Marshal(Value{Value: 132})

	map_param["index"], map_param["type"], map_param["id"], map_param["bodyJson"] = "megacorp", "employee", "53", string(e1)
	// ret := Create(map_param)
	// fmt.Printf(ret)
	get_ret := Gets(map_param)
	fmt.Printf(get_ret.Id)

}
