package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/utils"
	"github.com/olivere/elastic"
	"golang.org/x/net/context"
	"log"
	"os"
	"time"
)
var successNum,failedNum = 0,0
var document = 0
var host = utils.ELASTIC_SEARCH_HOST
var client *elastic.Client
type author struct {
	name string
	org string
	org_id string
	id string
}
type venues struct {
	id string
	name string
}
type fos struct {
	name string
	w string
}
type paper strucct{
	id string
	title string   
	authors []author
	venues []venue
	year string
	keywords []string
	references []string
	n_citation int
	page_start int
	page_end int
	doc_type string
	lang string 
	publisher string
	volume string 
	issue string 
	issn string 
	isbn string 
	doi string 
	pdf string 
	url string 
	abstract string 
}

type Author struct {
	id string   
	name string 
	normalized_name string 
	orgs []org 
	org org
	last_known_aff_id string 
	position string 
	n_pubs int 
	n_citation int 
	h_index int 
	tags_t string 
	tags_w int
	pubs_i string 
	pubs_r int
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
}

func InputData(file_path string,index string,index_type string) {

	open, err := os.Open(file_path)
	if err != nil {
		fmt.Println("打开失败")
	}
	scanner := bufio.NewScanner(open)
	i := 0
	bultRequest := client.Bulk()
	for scanner.Scan() {
		bodyJsonbyte := scanner.Bytes()
		var param_paper map[string]string = make(map[string]string)
		param_paper["index"] = index
		param_paper["type"] = index_type
		var tempMap map[string]interface{}
		err := json.Unmarshal(bodyJsonbyte, &tempMap)
		if err != nil {
			panic(err)
		}
		param_paper["id"] = fmt.Sprintf("%s", tempMap["id"])
		res := elastic.NewBulkIndexRequest().
			Index(param_paper["index"]).
			Type(param_paper["type"]).
			Id(param_paper["id"]).
			Doc(tempMap)
		bultRequest.Add(res)
		// fmt.Println(tempMap)
		//fmt.Println(tempMap["id"])
		//_ = service.Create(param_paper)
		//fmt.Println("return ", ret)
		if(i %10000 == 0){
			response , err := bultRequest.Do(context.Background())
			if err != nil {
				panic(err)
			}
			fmt.Println("success",len(response.Succeeded()),"failed",len(response.Failed()))
			failedNum += len(response.Failed())
			successNum += len(response.Succeeded())
		}
		i++

		//fmt.Println("document",document,"line sum", i)
	}
	response , err := bultRequest.Do(context.Background())
	failedNum += len(response.Failed())
	successNum += len(response.Succeeded())
	fmt.Println("Over document",document)
	fmt.Println("success",len(response.Succeeded()),"failed",len(response.Failed()))
	fmt.Println("successnum",successNum,"failed",failedNum)
}
func main() {
	//service.Init()
	Init()
	//InputData("E:\\Paper\\aminer_papers_0\\aminer_papers_0.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_0\\aminer_papers_1.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_0\\aminer_papers_2.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_0\\aminer_papers_3.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_1\\aminer_papers_4.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_1\\aminer_papers_5.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_1\\aminer_papers_6.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_1\\aminer_papers_7.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_2\\aminer_papers_8.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_2\\aminer_papers_9.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_2\\aminer_papers_10.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_2\\aminer_papers_11.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_3\\aminer_papers_12.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_3\\aminer_papers_13.txt","aminer","paper");document ++
	//InputData("E:\\Paper\\aminer_papers_3\\aminer_papers_14.txt","aminer","paper");document ++
	fmt.Println("paper save end") ; document = 0
	InputData("E:\\Paper\\aminer_authors_0\\aminer_authors_0.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_0\\aminer_authors_1.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_0\\aminer_authors_2.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_0\\aminer_authors_3.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_0\\aminer_authors_4.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_1\\aminer_authors_5.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_1\\aminer_authors_6.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_1\\aminer_authors_7.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_1\\aminer_authors_8.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_1\\aminer_authors_9.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_2\\aminer_authors_10.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_2\\aminer_authors_11.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_2\\aminer_authors_12.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_2\\aminer_authors_13.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_2\\aminer_authors_14.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_3\\aminer_authors_15.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_3\\aminer_authors_16.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_3\\aminer_authors_17.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_3\\aminer_authors_18.txt","aminer_author","author");document ++
	InputData("E:\\Paper\\aminer_authors_3\\aminer_authors_19.txt","aminer_author","author");document ++
	//	ret, _ := service.Gets(param_paper)
	//
	//	body_byte, _ := json.Marshal(ret.Source)
	//	fmt.Println(string(body_byte))
}