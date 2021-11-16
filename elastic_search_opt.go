package main

import (
	"context"
	//"encoding/json"
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/service"
	"reflect"

	"gitee.com/online-publish/slime-scholar-go/utils"

	"github.com/olivere/elastic/v7"
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
	e1 := Employee{"zht", "zhou", 18, "zht tql!!!!", []string{"coding"}}
	put1, err := client.Index().
		Index("megacorp").
		Type("employee").
		Id("5").
		BodyJson(e1).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)

	//使用字符串

}

//查找
//func gets() {
//	//通过id查找
//	get1, err := client.Get().Index("megacorp").Type("employee").Id("1").Do(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	if get1.Found {
//		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
//		var bb Employee
//		err := json.Unmarshal(*get1.Source, &bb) // 个人修改，原来模板存在问题
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(bb.FirstName)
//		fmt.Println(string(*get1.Source))
//	}
//}

//删除
func delete() {
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

func query_by_field(index string ,field string,content string)  *elastic.SearchResult{
	client = service.Client
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewMatchQuery(field, content))
	//boolQuery.Filter(elastic.NewRangeQuery("age").Gt("30"))
	searchResult, err := client.Search(index).Query(boolQuery).Pretty(true).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(searchResult.TotalHits())
	//var paper_ret map[string]interface{} = make(map[string]interface{})

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
	var res *elastic.SearchResult
	var err error
	//取所有
	res, err = client.Search("megacorp").Type("employee").Do(context.Background())
	printEmployee(res, err)

	//字段相等
	q := elastic.NewQueryStringQuery("last_name:Smith")
	res, err = client.Search("megacorp").Type("employee").Query(q).Do(context.Background())
	if err != nil {
		println(err.Error())
	}
	printEmployee(res, err)

	//条件查询
	//年龄大于30岁的
	boolQ := elastic.NewBoolQuery()
	boolQ.Must(elastic.NewMatchQuery("last_name", "smith"))
	boolQ.Filter(elastic.NewRangeQuery("age").Gt(30))
	res, err = client.Search("megacorp").Type("employee").Query(q).Do(context.Background())
	printEmployee(res, err)

	//短语搜索 搜索about字段中有 rock climbing
	matchPhraseQuery := elastic.NewMatchPhraseQuery("about", "rock climbing")
	res, err = client.Search("megacorp").Type("employee").Query(matchPhraseQuery).Do(context.Background())
	printEmployee(res, err)

	//分析 interests
	aggs := elastic.NewTermsAggregation().Field("interests")
	res, err = client.Search("megacorp").Type("employee").Aggregation("all_interests", aggs).Do(context.Background())
	printEmployee(res, err)

}

//
////简单分页
func list(size, page int) {
	if size < 0 || page < 1 {
		fmt.Printf("param error")
		return
	}
	res, err := client.Search("megacorp").
		Type("employee").
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

func main() {
	service.Init()
	query_by_field("paper","title","Business")
//	Create()
//	gets()
//	//delete()

//	// gets()
//	// query()
//	// list(2, 1)
}
