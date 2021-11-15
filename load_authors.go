package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/olivere/elastic"
	"golang.org/x/net/context"
	"os"
	"strconv"
)

const AUTHOR_DIR = "H:\\Author"
const FILE_NUM = 3
const FILE_PREFIX = "mag_authors_"
const BULK_SIZE = 10000

var success_num,fail_num = 0,0

type pub struct {
	id   string  `json:"id"`
	author_order int `json:"author_order"`
}

type Author struct {
	id             string `json:"id"`
	name           string `json:"name"`
	org            string `json:"org"`
	affiliation_id string `json:"affiliation_id"`
	position       string `json:"position"`
	n_pubs         int    `json:"n_pubs"`
	n_citation     int    `json:"n_citation"`
	h_index        int    `json:"h_index"`
	//pubs [10000]pub `json:"pubs"`
}

func jsonToAuthor(jsonStr string) (Author) {
	var item map[string]interface{} = make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &item)
	ok := false
	if err != nil {
		panic(err)
	}
	var author Author
	author.id = strconv.Itoa((int(item["id"].(float64))))
	author.name = item["name"].(string)
	author.org,ok = item["org"].(string)
	if !ok{author.org = ""}
	author.affiliation_id ,ok= item["last_known_aff_id"].(string)
	if !ok{author.affiliation_id = ""}
	author.position ,ok= item["position"].(string)
	if !ok{author.position = ""}
	author.n_pubs = int(item["n_pubs"].(float64))
	author.h_index, ok = item["h_index"].(int)
	if !ok{author.h_index = 0}
	n_citation ,ok := item["n_citation"].(float64)
	if !ok{author.n_citation = 0} else {author.n_citation = int(n_citation)}
	//pub_list := item["pubs"].([]interface {})
	//var pub_set [10000]pub
	//
	//for i,publish := range pub_list{
	//	var pub_item pub
	//	var map_item map[string]interface{} = make(map[string]interface{})
	//	map_item = publish.(map[string]interface{})
	//	publish_id := fmt.Sprintf("%s",item["i"])
	//	publish_item := pub{publish_id,int(map_item["r"].(float64))}
	//	pub_set[i] = publish_item
	//	if err != nil {panic(err)}
	//	pub_set[i] = pub_item
	//}
	//author.pubs = pub_set
	if err != nil {panic(err)}
	return author
}
func proc_file(file_path string) {
	open, err := os.Open(file_path)
	if err != nil {
		fmt.Println(file_path+"打开失败")
		return
	}
	scanner := bufio.NewScanner(open)
	i := 0

	client := service.ESClient
	bulkRequest := client.Bulk()
	for scanner.Scan() {
		author := jsonToAuthor(scanner.Text())
		doc := elastic.NewBulkIndexRequest().Index("author").Type("author").Id(author.id).Doc(author)
		bulkRequest = bulkRequest.Add(doc)
		if i % BULK_SIZE == 0 {
			response , err := bulkRequest.Do(context.Background())
			if(err!=nil) { panic(err)}
			success_num += len(response.Succeeded())
			fail_num += len(response.Failed())
			fmt.Println("success_num", success_num,"fail_num",fail_num)
		}
		i++
	}
	response , err := bulkRequest.Do(context.Background())
	if(err!=nil) { panic(err)}
	success_num += len(response.Succeeded())
	fail_num += len(response.Failed())
	fmt.Println("line sum", i)
	fmt.Println("success_num", success_num,"fail_num",fail_num)
}
func main() {
	service.Init()
	for i:=0 ; i<2;i++{
		for j:=0 ;j<FILE_NUM;j++{
			proc_file(AUTHOR_DIR+"\\"+FILE_PREFIX+string(i+'0')+"\\"+FILE_PREFIX+string(i*FILE_NUM+j+'0')+".txt")
		}
	}



}