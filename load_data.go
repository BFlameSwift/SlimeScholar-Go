package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/service"
	"github.com/olivere/elastic/v7"

	"golang.org/x/net/context"
	"io"
	"os"
	"strconv"
)

const AUTHOR_DIR = "H:\\Scholar"
const PAPER_DIR = "H:\\Scholar"
const FILE_NUM = 3
const AUTHOR_FILE_PREFIX = "aminer_authors_"
const PAPER_FILE_PREFIX = "s2-corpus-"
const BULK_SIZE = 100000
var fieldsMap  map[string]int = make(map[string]int)
var success_num, fail_num = 0, 0

type pub struct {
	id           string `json:"id"`
	author_order int    `json:"author_order"`
}

type Author struct {
	id             string `json:"id"`
	name           string `json:"name"`
	n_pubs         int    `json:"n_pubs"`
	n_citation     int    `json:"n_citation"`
	h_index        int    `json:"h_index"`
	papers []string `json:"papers"`
}


type Paper struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Abstract string `json:"abstract"`
	Url string `json:"url"`
	PdfUrls []string `json:"pdf_urls"`
	S2PdfUrl string `json:"s2pdf_urls"`
	InCitations []string `json:"in_citations"`
	OutCitations []string `json:"out_citations"`
	FieldsOfStudy []string `json:"study_fields"`
	Year int `json:"year"`
	Venue string   `json:"venue"`
	JournalName string `json:"journal_name"`
	JournalVolume string `json:"journal_volume"`
	JournalPages string `json:"journal_pages"`
	Doi string `json:"doi"`
	DoiUrl string `json:"doi_url"`
	MagId string `json:"mag_id"`
	Authors []Author `json:"authors"`
}

func JsonToPaper(jsonStr string) Paper {
	var item map[string]interface{} = make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &item)
	ok := false
	if err != nil {
		panic(err)
	}
	var paper Paper
	paper.Id = item["id"].(string)
	paper.Title = item["title"].(string)
	paper.Abstract,ok = item["abstract"].(string)
	if !ok {paper.Abstract = ""}
	paper.Url = item["s2Url"].(string)
	paper.S2PdfUrl = item["s2PdfUrl"].(string)
	year,ok := item["year"].(float64)
	if(!ok){year = 0}
	paper.Year = int(year)
	paper.JournalPages = item["journalPages"].(string)
	paper.JournalName = item["journalName"].(string)
	paper.JournalVolume = item["journalVolume"].(string)
	paper.Doi = item["doi"].(string)
	paper.DoiUrl = item["doiUrl"].(string)
	pdf_urls  := make([]string,10000)
	for i,url := range (item["pdfUrls"].([]interface{})){
		pdf_urls[i] = url.(string)
	};paper.PdfUrls = pdf_urls
	in_citations  := make([]string,10000)
	for i,str := range (item["inCitations"].([]interface{})){
		in_citations[i] = str.(string)
	};paper.InCitations = in_citations
	out_citations  := make([]string,10000)
	for i,str := range (item["outCitations"].([]interface{})){
		out_citations[i] = str.(string)
	};paper.OutCitations = out_citations
	fields  := make([]string,10000)
	_,ok = item["FieldsOfStudy"].([]interface{})
	if !ok{item["FieldsOfStudy"] = make([]interface{},0)   }
	for i,str := range (item["FieldsOfStudy"].([]interface{})){
		fields[i] = str.(string)
	};paper.FieldsOfStudy = fields
	authors := make([]Author,10000)
	_,ok = item["authors"].([]map[string]interface{})
	if !ok{item["authors"] = make([]map[string]interface{},0)   }
	for i,item_author := range (item["authors"].([]map[string]interface{})){
		author_new := Author{id: item_author["id"].(string),name: item_author["name"].(string)}
		authors[i] = author_new
	};paper.Authors = authors

	//author.position, ok = item["position"].(string)
	//if !ok {
	//	author.position = ""
	//}

	if err != nil {
		panic(err)
	}
	return paper
}
func proc_file(file_path string, index string) {
	open, err := os.Open(file_path)
	if err != nil {
		fmt.Println(file_path + "打开失败")
		return
	}
	scanner := bufio.NewScanner(open)
	i := 0
	fin, error := os.OpenFile(file_path, os.O_RDONLY, 0)
	if error != nil {
		panic(error)
	}
	defer fin.Close()
	client := service.ESClient
	bulkRequest := client.Bulk()
	reader := bufio.NewReader(fin)
	for  {
		line,error_read := reader.ReadString('\n')
		if(len(line) == 0){break;}
		json_str:= line

		//_ = JsonToPaper(json_str)
		//if(i<5){fmt.Println(paper)}
		var m map[string]interface{}
		_ = json.Unmarshal([]byte(json_str), &m)
		var fields []interface{} = m["fieldsOfStudy"].([]interface{})
		for _,field := range(fields){
			fieldsMap[field.(string)] += 1
			if field.(string) == "Computer Science" || field.(string) == "Mathematics"{
				doc := elastic.NewBulkIndexRequest().Index(index).Id(m["id"].(string)).Doc(m)

				bulkRequest.Add(doc)
				if i%BULK_SIZE == 0 {
					response, err := bulkRequest.Do(context.Background())
					if err != nil {
						panic(err)
					}
					success_num += len(response.Succeeded())
					fail_num += len(response.Failed())
					fmt.Println("success_num", success_num, "fail_num", fail_num)

				}
				break
			}
		}


		if error_read != nil {
			if err == io.EOF {
				fmt.Printf("%#v\n", line)
				break
			}
			panic(err)
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	response, err := bulkRequest.Do(context.Background())
	if err != nil {
		panic(err)
	}
	success_num += len(response.Succeeded())
	fail_num += len(response.Failed())
	fmt.Println("line sum", i)
	fmt.Println("success_num", success_num, "fail_num", fail_num)
	fmt.Println(fieldsMap)
}

func proc_author(file_path string, index string) {
	open, err := os.Open(file_path)
	if err != nil {
		fmt.Println(file_path + "打开失败")
		return
	}
	scanner := bufio.NewScanner(open)
	i := 0
	fin, error := os.OpenFile(file_path, os.O_RDONLY, 0)
	if error != nil {
		panic(error)
	}
	defer fin.Close()
	client := service.ESClient
	bulkRequest := client.Bulk()
	reader := bufio.NewReader(fin)
	for  {
		line,error_read := reader.ReadString('\n')
		if(len(line) == 0){break;}
		json_str := line

		//_ = JsonToPaper(json_str)
		//if(i<5){fmt.Println(paper)}
		var m map[string]interface{}
		_ = json.Unmarshal([]byte(json_str), &m)
		if len(m["author_id"].([]interface{})) == 0{continue} // 数据501行中存在"author_id": [],  过滤
		m["author_id"] = m["author_id"].([]interface{})[0].(string)
		doc := elastic.NewBulkIndexRequest().Index(index).Id(m["author_id"].(string)).Doc(m)
		bulkRequest.Add(doc)
		if i%BULK_SIZE == 0 {
			response, err := bulkRequest.Do(context.Background())
			if err != nil {
				panic(err)
			}
			success_num += len(response.Succeeded())
			fail_num += len(response.Failed())
			fmt.Println("success_num", success_num, "fail_num", fail_num)

		}

		if error_read != nil {
			if err == io.EOF {
				fmt.Printf("%#v\n", line)
				break
			}
			panic(err)
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	response, err := bulkRequest.Do(context.Background())
	if err != nil {
		panic(err)
	}
	success_num += len(response.Succeeded())
	fail_num += len(response.Failed())
	fmt.Println("line sum", i)
	fmt.Println("success_num", success_num, "fail_num", fail_num)
	fmt.Println(fieldsMap)
}
func load_paper() {
	service.Init()
	for i := 0; i < 6000; i++ {
		var str string
		if(i<1000){str = fmt.Sprintf("%03d",i);}else{str = strconv.Itoa(i)}
		fmt.Println(str)
		proc_file(PAPER_DIR+"\\"+PAPER_FILE_PREFIX+str, "paper")

	}
}
func load_authors(){
	service.Init()
	proc_author("H:\\Scholarauthors.txt","author")
}
func load_journal(){
	service.Init()
	proc_author("H:\\journal.txt","journal")
}
func print1(){
	for i := 0 ;i<1 ;i++{
		fmt.Printf("%s\n", fmt.Sprintf("%04d",i))
	}
}
func main() {
	//load_paper()
	//print1()
	//load_authors()
	load_journal()
	//fmt.Println(fieldsMap)
}
