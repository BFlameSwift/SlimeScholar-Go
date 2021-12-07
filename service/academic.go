package service

import (
	"encoding/json"
	"errors"
	"gitee.com/online-publish/slime-scholar-go/global"
	"gitee.com/online-publish/slime-scholar-go/model"
	"gorm.io/gorm"
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
func FindExistingTransfer(author_id string, paper_id string, user_id uint64, kind int) (transfer *model.Transfer, notFound bool) {
	err := global.DB.Where("author_id = ? AND paper_id = ? AND user_id = ? AND kind = ?", author_id, paper_id, user_id, kind).First(&transfer).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return transfer, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return transfer, false
	}
}

func TransferPaper(user model.User, author_id string, paper_id string, kind int, obj_user_id uint64) {
	_, notFound := FindExistingTransfer(author_id, paper_id, user.UserID, kind)
	if notFound {
		transfer := model.Transfer{UserID: user.UserID, AuthorId: author_id, PaperId: paper_id, Kind: kind, Status: 0, ObjUserID: obj_user_id}
		if err := global.DB.Create(&transfer).Error; err != nil {
			panic(err)
		}
	}
}
func FindAllAuthorManagePapers(author_id string) (transfer_list *[]model.Transfer, notFound bool) {
	err := global.DB.Where("author_id = ?", author_id).Find(&transfer_list).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return transfer_list, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return transfer_list, false
	}
}

func PaperMapToPaperList(m map[string]interface{}) (ret_list []interface{}) {
	for key := range m {
		item := m[key]
		author_map := make(map[string]interface{})
		author_map["rel"] = item.(map[string]interface{})["authors"]
		item.(map[string]interface{})["authors"] = ParseRelPaperAuthor(author_map)["rel"]
		ret_list = append(ret_list, item)
	}
	return ret_list
}
func GetAuthorAllPaper(author_id string) (paper_list []interface{}) {
	paper_result := QueryByField("paper", "authors.aid.keyword", author_id, 1, 100)
	paper_ids_origin := make([]string, 0, 10000)
	//authors_map := make(map[string]interface{})
	for _, hit := range paper_result.Hits.Hits {
		hit_map := make(map[string]interface{})
		err := json.Unmarshal([]byte(hit.Source), &hit_map)
		if err != nil {
			panic(err)
		}
		paper_ids_origin = append(paper_ids_origin, hit_map["paper_id"].(string))
	}
	transfer_list, notFound := FindAllAuthorManagePapers(author_id)
	papre_ids_del := make([]string, 0, 1000)
	// 找到应该删除的paper和应该添加的paperids
	if !notFound || len(*transfer_list) != 0 {
		for _, transfer := range *transfer_list {
			if transfer.Status == 1 {
				if transfer.Kind == 2 || transfer.Kind == 1 {
					papre_ids_del = append(papre_ids_del, transfer.PaperId)
				} else if transfer.Kind == 0 {
					paper_ids_origin = append(paper_ids_origin, transfer.PaperId)
				}
			}
		}
	}
	// 去重与删除操作
	paper_ids_final := make([]string, 0, 1000)
	paper_ids_map := make(map[string]int)
	for _, id := range paper_ids_origin {
		paper_ids_map[id] = 1
	}
	for _, id := range papre_ids_del {
		paper_ids_map[id] = 0
	}
	for key := range paper_ids_map {
		if paper_ids_map[key] == 1 {
			paper_ids_final = append(paper_ids_final, key)
		}
	}
	paper_map := IdsGetItems(paper_ids_final, "paper")
	//fmt.Println(paper_map)

	return PaperMapToPaperList(paper_map)
}

// 判断作者是否已经入驻
func JudgeAuthorIsSettled(author_id string) (bool, uint64) {
	submit, notFound := QueryASubmitByAuthor(author_id)
	return !notFound, submit.UserID
}

// 未入驻作者在展示个人中心之前的格式转化
func GetAuthorMsg(author_id string) (author_map map[string]interface{}) {
	author_json := GetsByIndexIdWithout("author", author_id)
	if author_json == nil {
		return nil
	}
	err := json.Unmarshal(author_json.Source, &author_map)
	if err != nil {
		panic(err)
	}
	if author_map["affiliation_id"].(string) != "" {
		author_map["affiliation"] = GetsByIndexIdWithout("affiliation", author_map["affiliation_id"].(string))
	} else {
		author_map["affiliation"] = ""
	}

	return author_map
}
