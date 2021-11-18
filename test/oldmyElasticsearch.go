package main

//reference 英文官网网址 https://pkg.go.dev/github.com/olivere/elastic?utm_source=godoc#section-documentation
//reference 中文文档，仅有基本操作  https://www.topgoer.com/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%93%8D%E4%BD%9C/go%E6%93%8D%E4%BD%9Celasticsearch/%E6%93%8D%E4%BD%9Celasticsearch.html

// Kibana 可视化操作 https://www.topgoer.com/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%93%8D%E4%BD%9C/go%E6%93%8D%E4%BD%9Celasticsearch/kibana%E5%AE%89%E8%A3%85.html
import (
	"context"
	"encoding/json"
	"fmt"

	"gitee.com/online-publish/slime-scholar-go/utils"
	"github.com/olivere/elastic/v7"
)

type employee struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

var client *elastic.Client

var host = utils.ELASTIC_SEARCH_HOST // 自己定义常量

type MyType struct {
	Id string `json:"id"`
}

//初始化
func OldInit() {
	var err error
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(host))
	if err != nil {
		panic(err)
	}
	_, _, err = client.Ping(host).Do(context.Background())
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	_, err = client.ElasticsearchVersion(host)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Elasticsearch version %s\n", esversion)
}

//创建
// index 通过index 查找 id 主键，直接查找
func OldCreate(index_type string, index string, id string, body string) {

	//使用结构体
	// e1 := employee{"zht", "zhou", 18, "zht tql!!!!", []string{"coding"}}
	mytype := MyType{Id: id}
	fmt.Println("index:", index_type, index, body, mytype)
	e1 := employee{"zht", "zhou", 18, "zht tql!!!!", []string{"coding"}}
	put1, err := client.Index().
		Index("megacorp").
		Type("employee").
		Id("5").
		BodyJson(e1).
		Do(context.Background())
	//put1, err := client.Index().
	//	Index(index).
	//	Type(index_type).
	//	Id(id).
	//	BodyJson(mytype).
	//	Do(context.Background())

	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)

}

//查找
// TODO 根据实际类别选择返回值 定制
func OldGets(index string, index_type string, id string) {
	//通过id查找
	get1, err := client.Get().Index(index).Type(index_type).Id(id).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
		var bb MyType
		err := json.Unmarshal(get1.Source, &bb) // 个人修改，原来模板存在问题
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(bb.FirstName)
		fmt.Println(string(get1.Source))
	}
}

//删除
func OldDelete(index_type string, index string, id string) {
	res, err := client.Delete().Index(index).
		Type(index_type).
		Id(id).
		Do(context.Background())
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("delete result %s\n", res.Result)
}
