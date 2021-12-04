package service

import (
	"gitee.com/online-publish/slime-scholar-go/model"
	"time"
)

func MostSimplifyPaper(m map[string]interface{}) (ret map[string]interface{}) {
	ret["paper_id"] = m["paper_id"]
	ret["paper_title"] = m["paper_title"]
	return ret
}
func SimplifyPapers(inter []interface{}) []interface{} {
	ret_list := make([]interface{}, len(inter))
	for _, v := range inter {
		ret_list = append(ret_list, MostSimplifyPaper(v.(map[string]interface{})))
	}
	return ret_list
}

func BrowerPaper(paper map[string]interface{}) {
	title := paper["paper_title"].(string)
	authors := paper["authors"].([]interface{})
	paper_id := paper["paper_id"].(string)
	authors_name := ""
	for i, author := range authors {
		authors_name += author.(map[string]interface{})["name"].(string)
		if i < len(authors)-1 {
			authors_name += ", "
		}
	}
	browsing_history := model.BrowsingHistory{BrowsingTime: time.Now(), Title: title, Authors: authors_name, PaperID: paper_id}
	err := CreateBrowseHistory(&browsing_history)
	if err != nil {
		panic(err)
	}

}
