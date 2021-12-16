package service

import (
	"errors"
	"gitee.com/online-publish/slime-scholar-go/global"
	"gitee.com/online-publish/slime-scholar-go/model"
	"gitee.com/online-publish/slime-scholar-go/utils"
	"gorm.io/gorm"
	"strconv"
)

//创建标签
func CreateATag(tag *model.Tag) (err error) {
	if err = global.DB.Create(&tag).Error; err != nil {
		return err
	}
	return nil
}

//收藏文章-打标签
func CreateATagPaper(tagPaper *model.TagPaper) (err error) {
	if err = global.DB.Create(&tagPaper).Error; err != nil {
		return err
	}
	return nil
}

//用户第一次收藏文章
func CreateACollect(collect *model.Collect) (err error) {
	if err = global.DB.Create(&collect).Error; err != nil {
		return err
	}
	return nil
}

//点赞
func UpdateACollect(collect *model.Collect) (err error) {
	err = global.DB.Save(collect).Error
	return err
}

//查询用户所有标签
func QueryTagList(userID uint64) (tags []model.Tag) {
	tags = make([]model.Tag, 0)
	global.DB.Where("user_id=?", userID).Find(&tags)
	return tags
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
func QueryTagPaper(tagID uint64) (papers []model.TagPaper) {
	papers = make([]model.TagPaper, 0)
	global.DB.Where("tag_id=?", tagID).Order("create_time desc").Find(&papers)
	return papers
}

type PaperCollect struct{
	Num 		uint64 		`json:"num"`
	PaperId 	string 		`json:"paper_id"`
}
//查询收藏文章Top10
func QueryCollectTop10() (collects []PaperCollect) {
	collects = make([]PaperCollect, 0)
	global.DB.Table("collect").Select("COUNT(user_id) as num,paper_id").Group("paper_id").Order("COUNT(user_id) desc").Limit(10).Find(&collects)
	return collects
}

//查询用户的所有收藏文章
func QueryUserCollect(userid uint64)(collects []model.Collect){
	collects = make([]model.Collect, 0)
	//TODO 按时间排序
	global.DB.Where("user_id=?", userid).Find(&collects)
	return collects
}

//精确查询标签文章
func QueryATagPaper(tagID uint64, paperID string) (tagPaper model.TagPaper, not bool) {
	err := global.DB.Where("tag_id = ? AND paper_id = ?", tagID, paperID).First(&tagPaper).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return tagPaper, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return tagPaper, false
	}
}

//判断用户是否收藏文章
func QueryACollect(userID uint64,paperID string)(collect model.Collect, notFound bool) {
	db := global.DB
	db = db.Where("user_id = ?", userID)
	db = db.Where("paper_id = ?", paperID)
	err := db.First(&collect).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return collect, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return collect, false
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

//删除收藏
func DeleteACollect(ID uint64) (err error) {
	if err = global.DB.Where("id = ?", ID).Delete(model.Collect{}).Error; err != nil {
		return err
	}
	return nil
}

//创建评论
func CreateAComment(comment *model.Comment) (notCreated bool) {
	if err := global.DB.Create(&comment).Error; err != nil {
		//更新回复数量
		com := comment
		for com.RelateID != 0 {
			relateCom, _ := QueryAComment(com.RelateID)
			relateCom.ReplyCount++
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

//点赞
func UpdateCommentLike(comment *model.Comment, user model.User) (err error) {
	comment.Like++
	err = global.DB.Save(comment).Error
	if err != nil {
		return err
	}

	like := model.Like{UserID: user.UserID, CommentID: comment.CommentID}
	err = global.DB.Create(&like).Error
	return err
}

//查询用户是否点赞评论
func UserLike(userID uint64, commentID uint64) (isLike bool) {
	like := model.Like{}
	err := global.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&like).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return true
	}
}

//取消点赞
func CancelLike(comment *model.Comment, user model.User) (notFound bool) {
	like := model.Like{}
	err := global.DB.Where("user_id = ? AND comment_id = ?", user.UserID, comment.CommentID).First(&like).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		global.DB.Delete(&like)
		comment.Like--
		global.DB.Save(&comment)
		return false
	}
}

//根据文献id获取文献所有评论
func QueryComsByPaperId(paperId string) (coms []model.Comment) {
	coms = make([]model.Comment, 0)
	global.DB.Where(map[string]interface{}{"paper_id": paperId, "relate_id": 0}).Order("comment_time desc").Find(&coms)
	return coms
}

// 根据文献Id获取文献的所有tag
func QueryTagByPaperId(paperId string) (tags []model.TagPaper) {
	global.DB.Where(map[string]interface{}{"paper_id": paperId}).Find(&tags)
	return tags
}

func QueryTagByTagId(tagId uint64) (tag model.Tag, notFound bool) {
	err := global.DB.Where("tag_id = ?", tagId).First(&tag).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return tag, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return tag, false
	}
}

//查询回复对应的最初的评论
func QueryABaseCom(comment *model.Comment) (base *model.Comment) {
	for comment.RelateID != 0 {
		id := comment.RelateID
		comment, _ = QueryAComment(id)
	}
	return comment
}

//查询某条评论的所有回复
func QueryComReply(relateID uint64) (coms []model.Comment) {
	coms = make([]model.Comment, 0)
	global.DB.Where("relate_id = ?", relateID).Order("comment_time").Find(&coms)
	tmp := coms
	for _, com := range tmp {
		comcom := QueryComReply(com.CommentID)
		for _, tmptmp := range comcom {
			coms = append(coms, tmptmp)
		}
	}
	return coms
}

// 根据Paperid找到认领该paper的username
func PaperGetCollectedUsers(paperId string) []string {
	tags := QueryTagByPaperId(paperId)

	tagUserMap := make(map[string]interface{})
	for _, tag := range tags {
		tag_user, notFound := QueryTagByTagId(tag.TagID)
		if !notFound {
			tagUserMap[tag_user.Username] = 1
		}
	}
	return GetMapAllKey(tagUserMap)
}

// 根据paperids列表找到与用户找到是否被用户所收藏
func PapersGetIsCollectedByUser(paperIds []string, user model.User) (ret []interface{}) {
	for _, paperId := range paperIds {
		item := make(map[string]interface{})
		item["paper_id"] = paperId
		item["is_collected"] = false
		for _, the_username := range PaperGetCollectedUsers(paperId) {
			if the_username == user.Username {
				item["is_collected"] = true
				break
			}
		}
		ret = append(ret, item)
	}
	return ret
}

// 根据userid与be_follow_user_id 实现关注操作，user关注列表增加，befollow被关注列表增加
func FollowUser(userId uint64, followedUserId uint64) {
	userIdStr := strconv.Itoa(int(userId))
	followedUserIdStr := strconv.Itoa(int(followedUserId))
	RedisSaveValue(utils.FOLLOW_USER_PREFIX+userIdStr, followedUserIdStr)
	RedisSaveValue(utils.BE_FOLLOWED_USER_PREFIX+followedUserIdStr, userIdStr)
}

func CanCelFollowUser(userId uint64, followedUserId uint64) bool {
	userIdStr := strconv.Itoa(int(userId))
	followedUserIdStr := strconv.Itoa(int(followedUserId))
	/// 用户未关注，却仍然使用
	if !RedisKeyIsExistValue(utils.FOLLOW_USER_PREFIX+userIdStr, followedUserIdStr) {
		return false
	}
	RedisRemoveValue(utils.FOLLOW_USER_PREFIX+userIdStr, followedUserIdStr)
	RedisRemoveValue(utils.BE_FOLLOWED_USER_PREFIX+followedUserIdStr, userIdStr)
	return true
}

// 根据userid与be_follow_user_id 实现关注操作，user关注列表增加，befollow被关注列表增加
func GetUserFollowingList(userId uint64) (userIdList []uint64) {
	userIdStr := strconv.Itoa(int(userId))
	list := RedisGetValueList(utils.FOLLOW_USER_PREFIX + userIdStr)
	for _, item := range list {
		itemU64, _ := strconv.ParseUint(item, 10, 64)
		userIdList = append(userIdList, itemU64)
	}
	return userIdList
}

// 格努befollowuserid获取用户的粉丝列表
func GetUserFollowedList(userId uint64) (userIdList []uint64) {
	userIdStr := strconv.Itoa(int(userId))
	list := RedisGetValueList(utils.BE_FOLLOWED_USER_PREFIX + userIdStr)
	for _, item := range list {
		itemU64, _ := strconv.ParseUint(item, 10, 64)
		userIdList = append(userIdList, itemU64)
	}
	return userIdList
}
