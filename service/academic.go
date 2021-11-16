package service

import (
	"gitee.com/online-publish/slime-scholar-go/global"
	"gitee.com/online-publish/slime-scholar-go/model"
)

//创建标签
func CreateATag(tag *model.Tag) (err error) {
	if err = global.DB.Create(&tag).Error; err != nil {
		return err
	}
	return nil
}

//查询用户所有标签
func QueryTagList(userID uint64)(tags []model.Tag, not bool){
	tags = make([]model.Tag, 0)
	db := global.DB
	db = db.Where("user_id=?",userID)
    if err := db.Find(&tags).Error; err != nil {
		return nil, true
	}
	return tags, false	
}