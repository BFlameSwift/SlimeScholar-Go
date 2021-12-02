package service

import (
	"errors"

	"gitee.com/online-publish/slime-scholar-go/global"
	"gitee.com/online-publish/slime-scholar-go/model"
	"gorm.io/gorm"
)

//创建标签
func CreateATag(tag *model.Tag) (err error) {
	if err = global.DB.Create(&tag).Error; err != nil {
		return err
	}
	return nil
}

//收藏文章
func CreateATagPaper(tagPaper *model.TagPaper) (err error) {
	if err = global.DB.Create(&tagPaper).Error; err != nil {
		return err
	}
	return nil
}

//查询用户所有标签
func QueryTagList(userID uint64) (tags []model.Tag, not bool) {
	tags = make([]model.Tag, 0)
	db := global.DB
	db = db.Where("user_id=?", userID)
	err := db.Find(&tags).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return tags, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return tags, false
	}
}

//查询用户某一个标签
func QueryATag(userID uint64, tagName string) (tag model.Tag, notFound bool) {
	db := global.DB
	db = db.Where("user_id = ?", userID)
	db = db.Where("tag_name = ?", tagName)
	err := db.First(&tag).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return tag, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return tag, false
	}
}

//查询用户标签下的所有文章
func QueryTagPaper(tagID uint64) (papers []model.TagPaper, not bool) {
	papers = make([]model.TagPaper, 0)
	db := global.DB
	db = db.Where("tag_id=?", tagID)
	err := db.Find(&papers).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return papers, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return papers, false
	}
}

//精确查询标签文章
func QueryATagPaper(tagID uint64, paperID string) (tagPaper model.TagPaper, not bool) {
	err := global.DB.Where("tag_id = ? AND paper_id = ?", tagID,paperID).First(&tagPaper).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return tagPaper, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return tagPaper, false
	}
}

//删除标签
func DeleteATag(tagID uint64) (err error) {
	if err = global.DB.Where("tag_id = ?", tagID).Delete(model.Tag{}).Error; err != nil {
		return err
	}
	if err = global.DB.Where("tag_id = ?", tagID).Delete(model.TagPaper{}).Error; err != nil {
		return err
	}
	return nil
}

//删除标签文章
func DeleteATagPaper(ID uint64) (err error) {
	if err = global.DB.Where("id = ?", ID).Delete(model.TagPaper{}).Error; err != nil {
		return err
	}
	return nil
}

//创建评论
func CreateAComment(comment *model.Comment) (notCreated bool) {
	if err := global.DB.Create(&comment).Error; err != nil {
		//更新回复数量
		com := comment
		for com.RelateID != 0{
			relateCom,_ := QueryAComment(com.RelateID)
			relateCom.ReplyCount++;
			global.DB.Save(relateCom)
			com = relateCom
		}
		return true
	}
	return false
}

// 根据评论 ID 查询某个评论
func QueryAComment(commentID uint64) (comment *model.Comment, notFound bool) {
	err := global.DB.Where("comment_id = ?", commentID).First(&comment).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return comment, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return comment, false
	}
}

//点赞或拉踩
func UpdateCommentLike(comment *model.Comment, option uint64) (err error) {
	if option == 0 {
		comment.Like++
	} else if option == 1 {
		comment.UnLike++
	}
	err = global.DB.Save(comment).Error
	return err
}

//根据文献id获取文献所有评论
func QueryComsByPaperId(paperId string)(coms []model.Comment){
	coms = make([]model.Comment, 0)
	global.DB.Where(map[string]interface{}{"paper_id":paperId,"relate_id":0}).Order("comment_time desc").Find(&coms)
	return coms
}

//查询某条评论的所有回复
func QueryComReply(relateID uint64)(coms []model.Comment){
	coms = make([]model.Comment, 0)
	global.DB.Where("relate_id = ?",relateID).Order("comment_time").Find(&coms)
	tmp := coms
	for _, com := range tmp{
		comcom := QueryComReply(com.CommentID)
		for _, tmptmp := range comcom{
			coms = append(coms,tmptmp)
		}
	}
	return coms
}

// func JsonToPaper(jsonStr string) model.Paper {
// 	var item map[string]interface{} = make(map[string]interface{})
// 	err := json.Unmarshal([]byte(jsonStr), &item)
// 	ok := false
// 	if err != nil {
// 		panic(err)
// 	}
// 	var paper model.Paper
// 	paper.Id = item["id"].(string)
// 	paper.Title = item["title"].(string)
// 	paper.Abstract, ok = item["paperAbstract"].(string)
// 	if !ok {
// 		paper.Abstract = ""
// 	}
// 	paper.Url = item["s2Url"].(string)
// 	paper.S2PdfUrl = item["s2PdfUrl"].(string)
// 	year, ok := item["year"].(float64)
// 	if !ok {
// 		year = 0
// 	}
// 	paper.Year = int(year)
// 	paper.JournalPages = item["journalPages"].(string)
// 	paper.JournalName = item["journalName"].(string)
// 	paper.JournalVolume = item["journalVolume"].(string)
// 	paper.Doi = item["doi"].(string)
// 	paper.DoiUrl = item["doiUrl"].(string)
// 	pdf_urls := make([]string, len(item["pdfUrls"].([]interface{})))
// 	for i, url := range item["pdfUrls"].([]interface{}) {
// 		pdf_urls[i] = url.(string)
// 	}
// 	paper.PdfUrls = pdf_urls
// 	in_citations := make([]string, len(item["inCitations"].([]interface{})))
// 	for i, str := range item["inCitations"].([]interface{}) {
// 		in_citations[i] = str.(string)
// 	}
// 	paper.InCitations = in_citations
// 	out_citations := make([]string, len(item["outCitations"].([]interface{})))
// 	for i, str := range item["outCitations"].([]interface{}) {
// 		out_citations[i] = str.(string)
// 	}
// 	paper.OutCitations = out_citations
// 	_, ok = item["FieldsOfStudy"].([]interface{})
// 	if !ok {
// 		item["FieldsOfStudy"] = make([]interface{}, 0)
// 	}
// 	fields := make([]string, len(item["FieldsOfStudy"].([]interface{})))
// 	for i, str := range item["FieldsOfStudy"].([]interface{}) {
// 		fields[i] = str.(string)
// 	}
// 	paper.FieldsOfStudy = fields
// 	authors := make([]model.Author, len(item["authors"].([]interface{})))
// 	_, ok = item["authors"].([]map[string]interface{})
// 	if !ok {
// 		item["authors"] = make([]map[string]interface{}, 0)
// 	}
// 	for i, item_author := range item["authors"].([]map[string]interface{}) {
// 		author_new := model.Author{AuthorId: item_author["id"].(string), AuthorName: item_author["name"].(string)}
// 		authors[i] = author_new
// 	}
// 	paper.Authors = authors

// 	//author.position, ok = item["position"].(string)
// 	//if !ok {
// 	//	author.position = ""
// 	//}

// 	if err != nil {
// 		panic(err)
// 	}
// 	return paper
// }
