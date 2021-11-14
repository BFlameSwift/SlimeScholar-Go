package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"gitee.com/online-publish/slime-scholar-go/service"
)

type Paper struct {
	id         string
	title      string
	authors    map[string]string
	year       string
	keywords   []string
	fos        []string
	n_citation int
	reference  []string
	doc_type   string
	lang       string
	publisher  string
	isbn       string
	doi        string
	pdf        string
	url        []string
	abstract   string
	page_start int
	page_end   int
	volume     int
}

func printMagPaper(file_path string) {
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
	fmt.Println("line sum", i)

}
func printAminerPaper(file_path string) {
	open, err := os.Open(file_path)
	if err != nil {
		fmt.Println("打开失败")
	}
	scanner := bufio.NewScanner(open)
	i := 0
	for scanner.Scan() {

			bodyJsonbyte := scanner.Bytes()
			var param_paper map[string]string = make(map[string]string)
			param_paper["index"] = "test"
			param_paper["type"] = "test_paper"
			var tempMap map[string]interface{}
			err := json.Unmarshal(bodyJsonbyte, &tempMap)
			if err != nil {
				panic(err)
			}
			json_map, _ := json.Marshal(tempMap)
			param_paper["id"] = fmt.Sprintf("%s", tempMap["id"])
			param_paper["bodyJson"] = fmt.Sprintf("%s", json_map)

			// fmt.Println(tempMap)
			//fmt.Println(tempMap["id"])
			ret := service.Create(param_paper)
			fmt.Println("return ", ret)


		i++
	}
	fmt.Println("line sum", i)

}

func main() {

	service.Init()
	printAminerPaper("E:\\Paper\\aminer_papers_0\\aminer_papers_0.txt")
	//var param_paper map[string]string = make(map[string]string)
	//param_paper["index"] = "test"
	//param_paper["type"] = "test_paper"
	//param_paper["id"] = "53e99784b7602d9701f3e131"
	//// printAminerPaper("D:\\Desktop\\aminer_papers_0.txt")
	//ret, _ := service.Gets(param_paper)
	//
	//body_byte, _ := json.Marshal(ret.Source)
	//fmt.Println(string(body_byte))
}
