package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/olivere/elastic"
	"golang.org/x/net/context"

	"gitee.com/online-publish/slime-scholar-go/service"
)

var successNum, failedNum = 0, 0
var thisfolder, document = 0, 0

func PrintMagPaper(file_path string) {
	open, err := os.Open(file_path)
	if err != nil {
		fmt.Println("打开失败")
	}
	scanner := bufio.NewScanner(open)
	i := 0
	for scanner.Scan() {
		if i < 10 {
			fmt.Println(scanner.Text())
		}
		i++
	}
	fmt.Println("folder ", thisfolder, "document", document, "line sum", i)

}
func InputData(file_path string, index string, index_type string) {
	client := service.Client
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
		if i%10000 == 0 {
			response, err := bultRequest.Do(context.Background())
			if err != nil {
				panic(err)
			}
			fmt.Println("success", len(response.Succeeded()), "failed", len(response.Failed()))
			failedNum += len(response.Failed())
			successNum += len(response.Succeeded())
		}
		i++

		//fmt.Println("document",document,"line sum", i)
	}
	response, err := bultRequest.Do(context.Background())
	failedNum += len(response.Failed())
	successNum += len(response.Succeeded())
	fmt.Println("Over document", document)
	fmt.Println("success", len(response.Succeeded()), "failed", len(response.Failed()))
	fmt.Println("successnum", successNum, "failed", failedNum)
}

func main() {
	// printMagPaper("E:\\Paper\\mag_papers_0\\mag_papers_1.txt")
	service.Init()
	//var param_paper map[string]string = make(map[string]string)
	//param_paper["index"] = "test"
	//param_paper["type"] = "test_paper"
	//param_paper["id"] = "53e99784b7602d9701f3e131"
	//InputData("E:\\Paper\\mag_papers_0\\mag_papers_1.txt")

	//InputData("E:\\Paper\\mag_papers_0\\mag_papers_0.txt","mag","paper");document ++
	//InputData("E:\\Paper\\mag_papers_0\\mag_papers_1.txt","mag","paper");document ++
	//InputData("E:\\Paper\\mag_papers_0\\mag_papers_2.txt","mag","paper");document ++
	//InputData("E:\\Paper\\mag_papers_0\\mag_papers_3.txt","mag","paper");document ++
	//InputData("E:\\Paper\\mag_papers_1\\mag_papers_4.txt","mag","paper");document ++
	//InputData("E:\\Paper\\mag_papers_1\\mag_papers_5.txt","mag","paper");document ++
	//InputData("E:\\Paper\\mag_papers_1\\mag_papers_6.txt","mag","paper");document ++
	//InputData("E:\\Paper\\mag_papers_1\\mag_papers_7.txt","mag","paper");document ++
	//InputData("E:\\Paper\\mag_papers_2\\mag_papers_8.txt","mag","paper");document ++
	//InputData("E:\\Paper\\mag_papers_2\\mag_papers_9.txt","mag","paper");document ++
	//InputData("E:\\Paper\\mag_papers_2\\mag_papers_10.txt","mag","paper");document =0

	InputData("E:\\Paper\\mag_authors_0\\mag_authors_0.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_0\\mag_authors_1.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_0\\mag_authors_2.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_0\\mag_authors_3.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_0\\mag_authors_4.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_1\\mag_authors_5.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_1\\mag_authors_6.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_1\\mag_authors_7.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_1\\mag_authors_8.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_1\\mag_authors_9.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_2\\mag_authors_10.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_2\\mag_authors_11.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_2\\mag_authors_12.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_2\\mag_authors_13.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_2\\mag_authors_14.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_3\\mag_authors_15.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_3\\mag_authors_16.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_3\\mag_authors_17.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_3\\mag_authors_18.txt", "mag_author", "author")
	document++
	InputData("E:\\Paper\\mag_authors_3\\mag_authors_19.txt", "mag_author", "author")
	document++

	//	ret, _ := service.Gets(param_paper)
	//
	//	body_byte, _ := json.Marshal(ret.Source)
	//	fmt.Println(string(body_byte))
}
