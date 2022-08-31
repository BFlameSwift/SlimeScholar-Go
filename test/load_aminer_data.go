package main

//
//import (
//	"bufio"
//	"encoding/json"
//	"fmt"
//	"os"
//	"strconv"
//
//	"github.com/BFlameSwift/SlimeScholar-Go/service"
//	"github.com/olivere/elastic/v7"
//	"golang.org/x/net/context"
//)
//
//const AUTHOR_DIR = "E:\\Paper"
//const PAPER_DIR = "E:\\Paper"
//const FILE_NUM = 3
//const AUTHOR_FILE_PREFIX = "aminer_authors_"
//const PAPER_FILE_PREFIX = "aminer_papers_"
//const BULK_SIZE = 10000
//
//var success_num, fail_num = 0, 0
//
//type pub struct {
//	id           string `json:"id"`
//	author_order int    `json:"author_order"`
//}
//
//type Author struct {
//	id             string `json:"id"`
//	name           string `json:"name"`
//	org            string `json:"org"`
//	affiliation_id string `json:"affiliation_id"`
//	position       string `json:"position"`
//	n_pubs         int    `json:"n_pubs"`
//	n_citation     int    `json:"n_citation"`
//	h_index        int    `json:"h_index"`
//	//pubs [10000]pub `json:"pubs"`
//}
//type author_paper struct{
//	id string `json:"id"`
//	org string `json:"org"`
//	org_id string `json:"org_id"`
//	name string `json:"name"`
//}
//type venue_paper struct{
//	id string `json:"id"`
//	name string `json:"name"`
//}
//type fos_paper struct{
//	name string `json:"name"`
//	w float64 `json:"weight"`
//}
//
//type Paper struct {
//	id string `json:"id"`
//	title string `json:"title"`
//	authors []author_paper `json:authors`
//	venue venue_paper `json:venue`
//	year int `json:"year"`
//	keywords []string `json:"keywords"`
//	fos []fos_paper `json:"fos"`
//	reference []string `json:"reference"`
//	n_citation int `json:"n_citation"`
//	page_start string `json:"page_start"`
//	page_end string `json:"page_end"`
//	publisher string `json:"publisher"`
//	volume string `json:"volume"`
//	issn string `json:"issn"`
//	isbn string `json:"isbn"`
//	doi string `json:"doi"`
//	pdf string `json:"pdf"`
//	url string `json:"url"`
//	abstract string `json:"abstract"`
//}
//
//func jsonToAuthor(jsonStr string) *Author {
//	var item map[string]interface{} = make(map[string]interface{})
//	err := json.Unmarshal([]byte(jsonStr), &item)
//	ok := false
//	if err != nil {
//		panic(err)
//	}
//	var author Author
//	author.id = strconv.Itoa((int(item["id"].(float64))))
//	author.name = item["name"].(string)
//	author.org, ok = item["org"].(string)
//	if !ok {
//		author.org = ""
//	}
//	author.affiliation_id, ok = item["last_known_aff_id"].(string)
//	if !ok {
//		author.affiliation_id = ""
//	}
//	author.position, ok = item["position"].(string)
//	if !ok {
//		author.position = ""
//	}
//	author.n_pubs = int(item["n_pubs"].(float64))
//	author.h_index, ok = item["h_index"].(int)
//	if !ok {
//		author.h_index = 0
//	}
//	n_citation, ok := item["n_citation"].(float64)
//	if !ok {
//		author.n_citation = 0
//	} else {
//		author.n_citation = int(n_citation)
//	}
//	//pub_list := item["pubs"].([]interface {})
//	//var pub_set [10000]pub
//	//
//	//for i,publish := range pub_list{
//	//	var pub_item pub
//	//	var map_item map[string]interface{} = make(map[string]interface{})
//	//	map_item = publish.(map[string]interface{})
//	//	publish_id := fmt.Sprintf("%s",item["i"])
//	//	publish_item := pub{publish_id,int(map_item["r"].(float64))}
//	//	pub_set[i] = publish_item
//	//	if err != nil {panic(err)}
//	//	pub_set[i] = pub_item
//	//}
//	//author.pubs = pub_set
//	if err != nil {
//		panic(err)
//	}
//	return &author
//}
//func proc_file(file_path string, index string) {
//	open, err := os.Open(file_path)
//	if err != nil {
//		fmt.Println(file_path + "打开失败")
//		return
//	}
//	scanner := bufio.NewScanner(open)
//	i := 0
//
//	client := service.ESClient
//	bulkRequest := client.Bulk()
//	for scanner.Scan() {
//		//author := jsonToAuthor(scanner.Text())
//		json_str := scanner.Text()
//		var m map[string]interface{}
//		_ = json.Unmarshal([]byte(json_str), &m)
//		doc := elastic.NewBulkIndexRequest().Index(index).Id(m["id"].(string)).Doc(m)
//
//		bulkRequest = bulkRequest.Add(doc)
//		if i%BULK_SIZE == 0 {
//			response, err := bulkRequest.Do(context.Background())
//			if err != nil {
//				panic(err)
//			}
//			success_num += len(response.Succeeded())
//			fail_num += len(response.Failed())
//			fmt.Println("success_num", success_num, "fail_num", fail_num)
//		}
//		i++
//	}
//	response, err := bulkRequest.Do(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	success_num += len(response.Succeeded())
//	fail_num += len(response.Failed())
//	fmt.Println("line sum", i)
//	fmt.Println("success_num", success_num, "fail_num", fail_num)
//}
//func load_authors() {
//	service.Init()
//	for i := 0; i < 2; i++ {
//		for j := 0; j < FILE_NUM; j++ {
//			proc_file(AUTHOR_DIR+"\\"+AUTHOR_FILE_PREFIX+string(i+'0')+"\\"+AUTHOR_FILE_PREFIX+strconv.Itoa(i*FILE_NUM+j)+".txt", "aminer_author")
//		}
//	}
//}
//func load_paper() {
//	service.Init()
//	for i := 0; i < 6; i++ {
//		for j := 0; j < FILE_NUM; j++ {
//			proc_file(PAPER_DIR+"\\"+PAPER_FILE_PREFIX+strconv.Itoa(i)+"\\"+PAPER_FILE_PREFIX+strconv.Itoa(i*FILE_NUM+j)+".txt", "aminer_paper")
//		}
//	}
//}
//func main() {
//	load_authors()
//	load_paper()
//
//}
