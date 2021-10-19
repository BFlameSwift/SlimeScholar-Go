package main

import (
	"context"

	"gitee.com/online-publish/slime-scholar-go/utils"
	"github.com/olivere/elastic"
)

var client *elastic.Client

var host = utils.ELASTIC_SEARCH_HOST // 自己定义常量

//初始化
func Init() {
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
func create(type string,index string,id string,body string) {

	//使用结构体
	e1 := Employee{"zht", "zhou", 18, "zht tql!!!!", []string{"coding"}}
	put1, err := client.Index().
		Index(index).
		Type(type).
		Id(id).
		BodyJson(body).
		Do(context.Background())

	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)

}