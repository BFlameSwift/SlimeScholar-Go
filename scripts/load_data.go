package scripts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"io"
	"os"

	"gitee.com/online-publish/slime-scholar-go/service"
	"golang.org/x/net/context"
)

const AUTHOR_DIR = "H:\\Scholar"
const PAPER_DIR = "E:\\Paper"
const FILE_NUM = 3
const AUTHOR_FILE_PREFIX = "aminer_authors_"
const PAPER_FILE_PREFIX = "s2-corpus-"

// TODO 设置bulk的大小
const BULK_SIZE = 10000

var fieldsMap map[string]int = make(map[string]int)
var success_num, fail_num = 0, 0

var max_citation_num = 0   // 看一下所有论文的最大引用数目
var max_references_num = 0 //
type pub struct {
	id           string `json:"id"`
	author_order int    `json:"author_order"`
}

func proc_single_paper(m map[string]interface{}) map[string]interface{} {
	m["rank"], _ = strconv.ParseInt(m["rank"].(string), 10, 64)
	m["citation_count"], _ = strconv.ParseInt(m["citation_count"].(string), 10, 64)
	m["reference_count"], _ = strconv.ParseInt(m["reference_count"].(string), 10, 64)
	m["paper_id"], _ = strconv.ParseInt(m["paper_id"].(string), 10, 64)
	m["year"], _ = strconv.ParseInt(m["year"].(string), 10, 64)
	return m
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
	//simpleBulkRequest := client.Bulk()
	reader := bufio.NewReader(fin)
	for {
		line, error_read := reader.ReadString('\n')
		if len(line) == 0 {
			break
		}
		json_str := line
		i++

		var m map[string]interface{} = make(map[string]interface{})
		err = json.Unmarshal([]byte(json_str), &m)
		if err != nil {
			panic(err)
		}

		if m["date"] == "" {
			m["date"] = "2020-05-01"
		}
		//m = proc_single_paper(m)
		// 因为这些数据到es中已经超过了100G 由于io的限制会导致查询的特别慢。。于是杉树一些不必哟啊的属性。 将引用，被引用信息分开存储，减少paper 索引的数据量
		//m["comment_num"] ,m["download_num"],m["collect_num"],m["browser_num"]= 0,0,0,0
		//TODO 存到数据库中吧
		doc := elastic.NewBulkIndexRequest().Index(index).Id(m["id"].(string)).Doc(m)
		//simpleBulkRequest.Add(elastic.NewBulkIndexRequest().Index("simple_paper").Id(m["id"].(string)).Doc(service.SimplifyPaper(m)))

		bulkRequest.Add(doc)
		if i%BULK_SIZE == 0 {
			response, err := bulkRequest.Do(context.Background())
			if err != nil {
				panic(err)
			}
			success_num += len(response.Succeeded())
			fail_num += len(response.Failed())
			fmt.Println("success_num", success_num, "fail_num", fail_num, time.Now())
			if fail_num > 0 {
				for _, item := range response.Failed() {
					fmt.Println(item.Error)
				}
			}
		}

		if error_read != nil {
			if err == io.EOF {
				fmt.Printf("%#v\n", line)
				break
			}
			panic(err)
		}
		// i++
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
	//response,err = simpleBulkRequest.Do(context.Background())
	if len(response.Failed()) > 0 {
		panic(err)
	}
	if fail_num > 0 {
		fmt.Println("error:")
	}
	for _, item := range response.Failed() {
		fmt.Println(item.Error)
	}
	fmt.Println("line sum", i)
	fmt.Println("success_num", success_num, "fail_num", fail_num, "max_citation_num", max_citation_num, "max_references_num", max_references_num)
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
	for {
		line, error_read := reader.ReadString('\n')
		if len(line) == 0 {
			break
		}
		json_str := line

		//_ = JsonToPaper(json_str)
		//if(i<5){fmt.Println(paper)}
		var m map[string]interface{}
		_ = json.Unmarshal([]byte(json_str), &m)
		//if len(m["author_id"].(interface{})) == 0 {
		//	continue
		//} // 数据501行中存在"author_id": [],  过滤
		//m["author_id"] = m["author_id"].([]interface{})[0].(string)
		m["paper_count"] = int(m["paper_count"].(float64))
		m["citation_count"] = int(m["citation_count"].(float64))
		m["rank"] = int(m["rank"].(float64))
		doc := elastic.NewBulkIndexRequest().Index(index).Id(m["author_id"].(string)).Doc(m)
		bulkRequest.Add(doc)
		if i%BULK_SIZE == 0 {
			response, err := bulkRequest.Do(context.Background())
			if err != nil {
				panic(err)
			}
			success_num += len(response.Succeeded())
			fail_num += len(response.Failed())
			fmt.Println("success_num", success_num, "fail_num", fail_num, time.Now())

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
func proc_journal(file_path string, index string, main_id string) {
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
	for {
		line, error_read := reader.ReadString('\n')
		if len(line) == 0 {
			break
		}
		json_str := line
		i++

		var m map[string]interface{}
		_ = json.Unmarshal([]byte(json_str), &m)
		if m[main_id] == nil {
			fmt.Println("linenum!!!!", i)
			continue
		}
		if index == "conference" {
			if m["start"].(string) == "" {
				m["start"] = "2021-11-30"
			}
			if m["end"].(string) == "" {
				m["end"] = "2021-11-30"
			}
		}
		//if len(m["author_id"].([]interface{})) == 0{continue} // 数据501行中存在"author_id": [],  过滤
		//m["id"] = m["id"].([]interface{})[0].(string)
		doc := elastic.NewBulkIndexRequest().Index(index).Id(m[main_id].(string)).Doc(m)
		bulkRequest.Add(doc)
		if i%BULK_SIZE == 0 {
			response, err := bulkRequest.Do(context.Background())
			if err != nil {
				panic(err)
			}
			success_num += len(response.Succeeded())
			fail_num += len(response.Failed())
			if fail_num > 0 {
				fmt.Println((response.Failed()[0].Error))
			}
			fmt.Println("success_num", success_num, "fail_num", fail_num, i, time.Now())
		}
		if error_read != nil {
			if err == io.EOF {
				fmt.Printf("%#v\n", line)
				break
			}
			panic(err)
		}
		// i++
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
	if fail_num > 0 {
		for _, fail := range response.Failed() {
			fmt.Println(fail.Error)
		}
	}
	fmt.Println("line sum", i)
	fmt.Println("success_num", success_num, "fail_num", fail_num)
	fmt.Println(fieldsMap)
}
func proc_paper_rel(file_path string, index string, main_id string, other_type string) {
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
	for {
		line, error_read := reader.ReadString('\n')
		if len(line) == 0 {
			break
		}
		json_str := line
		i++

		var m map[string]interface{}
		_ = json.Unmarshal([]byte(json_str), &m)
		if m[main_id] == nil {
			fmt.Println("linenum!!!!", i)
			continue
		}
		// TODO 数据格式为{paperid,rel(abstract、authors)}，让map只含有abstract或者authors、，直接插入

		m[other_type] = m["rel"]
		// fmt.Println(m["rel"])
		id := m[main_id].(string)
		delete(m, "rel")
		delete(m, main_id)
		//TODO 让map只含有abstract或者authors、，直接插入
		// fmt.Println(m)
		doc := elastic.NewBulkUpdateRequest().Index(index).Id(id).Doc(m).DocAsUpsert(true)
		bulkRequest.Add(doc)
		if i%BULK_SIZE == 0 {
			// TODO 每一个size 输出结果
			response, err := bulkRequest.Do(context.Background())
			if err != nil {
				panic(err)
			}
			fmt.Println(id)
			success_num += len(response.Succeeded())
			fail_num += len(response.Failed())
			if fail_num > 0 {
				fmt.Println((response.Failed()[0].Error))
			}
			fmt.Println("success_num", success_num, "fail_num", fail_num)
		}
		if error_read != nil {
			if err == io.EOF {
				fmt.Printf("%#v\n", line)
				break
			}
			panic(err)
		}
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
	if fail_num > 0 {
		for _, fail := range response.Failed() {
			fmt.Println(fail.Error)
		}
	}
	fmt.Println("line sum", i)
	fmt.Println("success_num", success_num, "fail_num", fail_num)
	fmt.Println(fieldsMap)
}

func proc_abstract_title(file_path string, index string, main_id string, other_type string, location int) {
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
	for {
		line, error_read := reader.ReadString('\n')
		if len(line) == 0 {
			break
		}
		// json_str := line/
		i++

		line_list := strings.Split(line, "\t")

		theMap := make(map[string]interface{})
		theMap[other_type] = line_list[location]
		// m[other_type] = m["rel"]
		// fmt.Println(m["rel"])
		id := line_list[0]
		// delete(m, "rel")
		// delete(m, main_id)
		// fmt.Println(m)
		doc := elastic.NewBulkUpdateRequest().Index(index).Id(id).Doc(theMap).DocAsUpsert(false)
		bulkRequest.Add(doc)
		if i%BULK_SIZE == 0 {
			response, err := bulkRequest.Do(context.Background())
			if err != nil {
				panic(err)
			}
			fmt.Println(id)
			success_num += len(response.Succeeded())
			fail_num += len(response.Failed())
			if fail_num > 0 {
				fmt.Println((response.Failed()[0].Error))
			}
			fmt.Println("success_num", success_num, "fail_num", fail_num, time.Now())
		}
		if error_read != nil {
			if err == io.EOF {
				fmt.Printf("%#v\n", line)
				break
			}
			panic(err)
		}
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
	if fail_num > 0 {
		for _, fail := range response.Failed() {
			fmt.Println(fail.Error)
		}
	}
	fmt.Println("line sum", i)
	fmt.Println("success_num", success_num, "fail_num", fail_num, time.Now())
	fmt.Println(fieldsMap)
}
func load_paper() {
	//cluster_block_exception index [paper] blocked by: [TOO_MANY_REQUESTS/12/disk usage exceeded flood-stage watermark, index has read-only-allow-delete block
	//考虑磁盘空间问题 ：方法：curl -XPUT -H "Content-Type: application/json" http://10.2.7.70:9204/_all/_settings -d '{"index.blocks.read_only_allow_delete": null}'
	//max_citation_num 220497 max_references_num 26676
	service.Init()
	proc_file("F:\\微软文献库\\myPapers.txt", "paper")
	//for i := 0; i < 6000; i++ {
	//	var str string
	//	if i < 1000 {
	//		str = fmt.Sprintf("%03d", i)
	//	} else {
	//		str = strconv.Itoa(i)
	//	}
	//	fmt.Println(str)
	//	proc_file(PAPER_DIR+"\\"+PAPER_FILE_PREFIX+str, "paper")
	//
	//}
}
func load_simple_paper() {
	//cluster_block_exception index [paper] blocked by: [TOO_MANY_REQUESTS/12/disk usage exceeded flood-stage watermark, index has read-only-allow-delete block
	//考虑磁盘空间问题 ：方法：curl -XPUT -H "Content-Type: application/json" http://10.2.7.70:9204/_all/_settings -d '{"index.blocks.read_only_allow_delete": null}'
	//max_citation_num 220497 max_references_num 26676
	service.Init()
	proc_file("F:\\微软文献库\\mySimplePapers.txt", "simple_paper")

}
func load_authors() {
	service.Init()
	proc_author("F:\\微软文献库\\myAuthors.txt", "author")
}
func load_journal() {
	service.Init()
	proc_journal("F:\\微软文献库\\myJournals.txt", "journal", "journal_id")
}

func load_abstract_title() {
	service.Init()
	proc_abstract_title("F:\\微软文献库\\Papers.txt", "abstract", "paper_id", "title", 3)
}

//func load_incitations() {
//	service.Init()
//	proc_journal("H:\\ScholarinCitations.txt", "incitations")
//}
func load_paper_author() {
	service.Init()
	proc_paper_rel("F:\\微软文献库\\mydata\\myPaperAuthorAffiliations.txt", "paper", "paper_id", "authors")
}
func load_paper_rel() {
	service.Init()
	proc_journal("F:\\微软文献库\\mydata\\myPaperReferences.txt", "reference", "paper_id")
}
func load_conference() {
	service.Init()
	proc_journal("F:\\微软文献库\\myConferenceInstances.txt", "conference", "conference_id")
}
func load_paper_url() {
	service.Init()
	proc_journal("F:\\微软文献库\\myPaperUrls.txt", "url", "paper_id")
}
func load_fields() {
	service.Init()
	proc_paper_rel("F:\\微软文献库\\myPaperFields.txt", "paper", "paper_id", "fields")
}
func load_abstract() {
	service.Init()
	for i := 5; i <= 5; i++ {
		str := strconv.Itoa(i)
		proc_journal("F:\\微软文献库\\myPaperAbstractsInvertedIndex.txt."+str, "abstract", "paper_id")
	}
}
func load_affiliations() {
	service.Init()
	proc_journal("F:\\微软文献库\\myAffiliations.txt", "affiliation", "affiliation_id")
}
func load_citation() {
	service.Init()
	proc_journal("F:\\微软文献库\\mydata\\myPaperCitationContexts.txt", "citation", "paper_id")
}
func load_reference_year() {
	service.Init()
	proc_abstract_title("H:\\mag\\Papers.txt", "reference", "paper_id", "year", 7)

}
func print1() {
	for i := 0; i < 1; i++ {
		fmt.Printf("%s\n", fmt.Sprintf("%04d", i))
	}
}
func loadPaperAbstract() {
	service.Init()
	// TODO 分成五进程直接插入
	for i := 5; i < i+1; i++ {
		str := strconv.Itoa(i)
		proc_paper_rel("H:\\myPaperAbstractsInvertedIndex.txt."+str, "paper", "paper_id", "abstract")
	}
}
func CountFile(filename string) int {
	open, err := os.Open(filename)
	if err != nil {
		fmt.Println(filename + "打开失败")
		return 0
	}
	scanner := bufio.NewScanner(open)
	i := 0
	lenStr := 0
	for scanner.Scan() {
		line := scanner.Text()
		lenStr += len(line)
		i += 1

	}
	fmt.Println(scanner.Err().Error())
	fmt.Println(lenStr)
	return i
}
func loadDataMain() {
	load_reference_year()
	load_abstract_title()
	load_fields()
	load_paper()
	load_affiliations()
	load_authors()
	load_fields()
	print1()
	load_simple_paper()
	load_journal()
	load_paper_rel()
	load_authors()
	load_abstract()
	load_citation()
	load_paper_author()
	loadPaperAbstract()
	load_conference()
	load_journal()
	load_paper_url()
}
