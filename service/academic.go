package service

import (
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
	trans, notFound := FindExistingTransfer(author_id, paper_id, user.UserID, kind)
	if notFound {
		transfer := model.Transfer{UserID: user.UserID, AuthorId: author_id, PaperId: paper_id, Kind: kind, Status: false, ObjUserID: obj_user_id}
		if err := global.DB.Create(&transfer).Error; err != nil {
			panic(err)
		}
	}
}
