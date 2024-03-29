package main

// elasticsearch 的基本操作，摘自go中文文档
//reference 英文官网网址 https://pkg.go.dev/github.com/olivere/elastic?utm_source=godoc#section-documentation
//reference 中文文档，仅有基本操作  https://www.topgoer.com/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%93%8D%E4%BD%9C/go%E6%93%8D%E4%BD%9Celasticsearch/%E6%93%8D%E4%BD%9Celasticsearch.html

// Kibana 可视化操作 https://www.topgoer.com/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%93%8D%E4%BD%9C/go%E6%93%8D%E4%BD%9Celasticsearch/kibana%E5%AE%89%E8%A3%85.html
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BFlameSwift/SlimeScholar-Go/initialize"
	"github.com/BFlameSwift/SlimeScholar-Go/service"
	"github.com/BFlameSwift/SlimeScholar-Go/utils"
	"github.com/olivere/elastic/v7"
	"log"
	"reflect"
)

var client *elastic.Client

var host = utils.ELASTIC_SEARCH_HOST

type Employee struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

//初始化
func Init() {
	//errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	var err error
	//这个地方有个小坑 不加上elastic.SetSniff(false) 会连接不上
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(host))
	if err != nil {

		panic(err)
	}
	_, _, err = client.Ping(host).Do(context.Background())
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	var esversion string
	esversion, err = client.ElasticsearchVersion(host)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
}

//创建
func Create() {

	//使用结构体
	e1 := Employee{"Jane", "Smith", 32, "I like to collect rock albums", []string{"music"}}
	put1, err := client.Index().
		Index("megacorp").
		Id("1").
		BodyJson(e1).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)

	//使用字符串
	e2 := `{"first_name":"John","last_name":"Smith","age":25,"about":"I love to go rock climbing","interests":["sports","music"]}`
	put2, err := client.Index().
		Index("megacorp").
		Id("2").
		BodyJson(e2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put2.Id, put2.Index, put2.Type)

	e3 := `{"first_name":"Douglas","last_name":"Fir","age":35,"about":"I like to build cabinets","interests":["forestry"]}`
	put3, err := client.Index().
		Index("megacorp").
		Id("3").
		BodyJson(e3).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put3.Id, put3.Index, put3.Type)

}

//查找
func gets(id string) {
	//通过id查找
	get1, err := client.Get().Index("megacorp").Id(id).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
		var bb Employee
		err := json.Unmarshal(get1.Source, &bb) // 个人修改，原来模板存在问题
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(bb.FirstName)
		fmt.Println(string(get1.Source))
	}
}

//删除
func delete(map[string]interface{}, string) {
	res, err := client.Delete().Index("megacorp").
		Type("employee").
		Id("1").
		Do(context.Background())
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("delete result %s\n", res.Result)
}
func update() {
	res, err := client.Update().
		Index("megacorp").
		Type("employee").
		Id("2").
		Doc(map[string]interface{}{"age": 88}).
		Do(context.Background())
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("update age %s\n", res.Result)
}

func query_by_field(index string, field string, content string) *elastic.SearchResult {
	client = service.Client
	//boolQuery := elastic.NewBoolQuery()
	q := elastic.NewMatchPhraseQuery(field+".name", content) //精确匹配
	//boolQuery.Must(elastic.NewMatchQuery(field+".name", content))
	//childQuery := elastic.NewHasChildQuery("name",boolQuery)
	//q := elastic.NewQueryStringQuery(field+".name:"+content)
	//boolQuery := elastic.NewTermQuery(field,content)
	//boolQuery.Filter(elastic.NewRangeQuery("age").Gt("30"))
	searchResult, err := client.Search(index).Query(q).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(searchResult.TotalHits())
	//var paper_ret map[string]interface{} = make(map[string]interface{})
	//var paper_item map[string]interface{} = make(map[string]interface{})
	//for _, item := range searchResult.Each(reflect.TypeOf(paper_item)) { //从搜索结果中取数据的方法
	//	t := item.(map[string]interface{})
	//	json_str, err := json.Marshal(item.(map[string]interface{}))
	//	if err != nil {
	//		panic(err)
	//	}
	//	paper := service.JsonToPaper(string(json_str))
	//	fmt.Printf("%#v\n", t)
	//	fmt.Printf("%#v\n", paper)
	//	fmt.Println(reflect.ValueOf(&paper).Elem())
	//}
	//for i,result := range(searchResult.Hits.Hits){
	//	json_str,err := json.Marshal(result.Source)
	//	if err != nil {panic(err)}
	//	var m map[string]interface{} = make(map[string]interface{})
	//	_ = json.Unmarshal([]byte(json_str),&m)
	//	if i<10{
	//		fmt.Println(i,m)
	//	}
	//}
	return searchResult
}

////搜索
func query() {
	//var res *elastic.SearchResult
	//var err error
	//取所有
	//res, err = client.Search("paper").Do(context.Background())
	//printEmployee(res, err)

	//字段相等
	client = service.ESClient
	q := elastic.NewQueryStringQuery("authors:name:James")
	searchResult, err := client.Search().Index("paper").Query(q).Do(context.Background())
	if err != nil {
		println(err.Error())
	}
	fmt.Println(searchResult.TotalHits())
	var paper_item map[string]interface{} = make(map[string]interface{})
	for _, item := range searchResult.Each(reflect.TypeOf(paper_item)) { //从搜索结果中取数据的方法
		t := item.(map[string]interface{})
		//json_str,err := json.Marshal(item.(map[string]interface{}))
		if err != nil {
			panic(err)
		}
		//paper := service.JsonToPaper(string(json_str))
		fmt.Printf("%#v\n", t)
		fmt.Println(len(t["inCitations"].([]interface{})))

		//fmt.Printf("%#v\n",paper)
		//fmt.Println(reflect.ValueOf(&paper).Elem())
	}
	//printEmployee(res, err)

	//条件查询
	//年龄大于30岁的
	//boolQ := elastic.NewBoolQuery()
	//boolQ.Must(elastic.NewMatchQuery("last_name", "smith"))
	//boolQ.Filter(elastic.NewRangeQuery("age").Gt(30))
	//res, err = client.Search("megacorp").Type("employee").Query(q).Do(context.Background())
	//printEmployee(res, err)

	//短语搜索 搜索about字段中有 rock climbing
	//matchPhraseQuery := elastic.NewMatchPhraseQuery("about", "rock climbing")
	//res, err = client.Search("megacorp").Type("employee").Query(matchPhraseQuery).Do(context.Background())
	//printEmployee(res, err)
	//
	////分析 interests
	//aggs := elastic.NewTermsAggregation().Field("interests")
	//res, err = client.Search("megacorp").Type("employee").Aggregation("all_interests", aggs).Do(context.Background())
	//printEmployee(res, err)

}

//
////简单分页
func list(size, page int) {
	if size < 0 || page < 1 {
		fmt.Printf("param error")
		return
	}
	res, err := client.Search("megacorp").
		Size(size).
		From((page - 1) * size).
		Do(context.Background())
	printEmployee(res, err)

}

//
//打印查询到的Employee
func printEmployee(res *elastic.SearchResult, err error) {
	if err != nil {
		print(err.Error())
		return
	}
	var typ Employee
	for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
		t := item.(Employee)
		fmt.Printf("%#v\n", t)
	}
}

func test_aggregation() {
	service.Init()
	aggs := elastic.NewTermsAggregation().
		Field("doctype.keyword") // 设置统计字段

	searchResult, err := service.Client.Search().
		Index("paper"). // 设置索引名
		Query(elastic.NewMatchAllQuery()). // 设置查询条件
		Aggregation("doctype", aggs). // 设置聚合条件，并为聚合条件设置一个名字, 支持添加多个聚合条件，命名不一样即可。
		Size(0). // 设置分页参数 - 每页大小,设置为0代表不返回搜索结果，仅返回聚合分析结果
		Do(context.Background()) // 执行请求

	if err != nil {
		// Handle error
		panic(err)
	}

	// 使用ValueCount函数和前面定义的聚合条件名称，查询结果
	// 使用Terms函数和前面定义的聚合条件名称，查询结果
	agg, found := searchResult.Aggregations.Terms("doctypee")
	if !found {
		log.Fatal("没有找到聚合数据")
	}
	fmt.Println(searchResult.TotalHits())
	// 遍历桶数据
	for _, bucket := range agg.Buckets {
		// 每一个桶都有一个key值，其实就是分组的值，可以理解为SQL的group by值
		bucketValue := bucket.Key
		// 打印结果， 默认桶聚合查询，都是统计文档总数
		fmt.Printf("bucket = %q 文档总数 = %d\n", bucketValue, bucket.DocCount)
	}

	// query_by_field("paper", "authors", "Christopher  Quince")
	//var str_list = []string{"02d5380d2fb7a81019b124b079306cc5cd3794d3", "d80d307a463df3526bf12ef1974afa7352f7b863"}
	// IdsGetPapers(str_list, "paper")
	//	Create
	//	gets()
	//	//delete()
	//	// gets()
	//	query()
	//	// list(2, 1)
}

func optMain() {
	service.Init()
	initialize.Init()
	//result := service.GetsByIndexIdWithout("test_paper_fields", "22367272")
	//fmt.Println(string(result.Source))
	//initialize.Init()

	//mostCitationPapers := GetMostCitationPapers(1000)
	//for _, id := range mostCitationPapers {
	//	service.RedisSaveValueSorted("most1000sort", id)
	//	//fmt.Println(id)
	//}
	////initialize.Init()
	//initialize.InitRedis()
	ids := service.GetMost1000CitationPaperIds()
	fmt.Println(ids, len(ids))

}
