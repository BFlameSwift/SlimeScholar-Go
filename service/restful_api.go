package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 向服务器发起请求调用aoi， 最初用作调用semanticscholar的api
func GetUrl(url string) string {
	res, err := http.Get(url)
	if err != nil {
		return ""
	}
	fmt.Println("request api:", res)
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		//panic(err)
		return ""
	}
	return string(robots)
}
func SemanticScholarApiSingle(mag_id string, field string) string {
	// 时刻记住go 参数传递数组是复制在传递，直接用指针可以节省开销

	str := GetUrl("https://api.semanticscholar.org/graph/v1/paper/MAG:" + mag_id + "?fields=" + field)
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		panic(err)
	}
	if m[field] == nil {
		return "null"
	} else {
		return m[field].(string)
	}

}

// 不过也用不上（
func SemanticScholarApiMulti(mag_id string, fields_pointer *[]string) map[string]interface{} {
	// 时刻记住go 参数传递数组是复制在传递，直接用指针可以节省开销
	fields := *fields_pointer
	fields_request := ""
	for i, field := range fields {
		fields_request += field
		if i < len(fields)-1 {
			fields_request += ", "
		}
	}
	str := GetUrl("https://api.semanticscholar.org/graph/v1/paper/MAG:" + mag_id + "?fields=" + fields_request)
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		panic(err)
	}

	return m
}

//func main() {
//	//fmt.Println(SemanticScholarApiSingle("1582271227","abstract"))
//	//fmt.Println(GetUrl("127.0.0.1:9200/paper/_count"))
//
//	//requests := make([]string,0,30)
//	//requests = append(requests,"abstract")
//	//requests = append(requests,"title")
//	//fmt.Println(SemanticScholarApi("1582271227",&requests))
//}
