package model

// 本文件下主要记录社交化带来的表
import "time"

type Message struct {
	MsgID      uint64    `gorm:"primary_key; not null" json:"msg_id"`
	Content    string    `gorm:"size :256;" json:"content"`
	Title      string    `gorm:"size:128" json:"title"`
	CreateTime time.Time `gorm:"type:datetime" json:"create_time"`
}

type Comment struct {
	CommentID   uint64    `gorm:"primary_key;not null" json:"comment_id"`
	Like        uint64    `gorm:"default:0" json:"like"`
	UnLike      uint64    `gorm:"default:0" json:"unlike"`
	UserID      uint64    `gorm:" not null;" json:"user_id"`
	PaperID     string    `gorm:"size:64" json:"paper_id"`
	CommentTime time.Time `gorm:"type:datetime" json:"comment_time"`
	Content     string    `gorm:"size:255" json:"content"`
	// OnTop       bool      `gorm:"default:false" json:"on_top"`
	ReplyCount	uint64	  `gorm:"default:0" json:"reply_count"`
	RelateID	uint64	  `gorm:"default:0" json:"relate_id"`
}

type Like struct { // 点赞
	IsLike    bool   `gorm:"default:false" json:"is_like"`
	CommentID uint64 `gorm:"primary_key;" json:"comment_id"`
	UserID    uint64 `gorm:"primary_key;" json:"user_id"`
}

type Follow struct {
	FollowID     uint64    `gorm:"primary_key; not null;" json:"id"`
	UserID       uint64    `gorm:"not null" json:"follow_user_id"`
	FollowUserID uint64    `gorm:"not null" json:"be_follow_user_id"`
	FollowTime   time.Time `gorm:"type:datetime" json:"follow_time"`
}

//标签
type Tag struct{
	TagID		uint64	`gorm:"primary_key;" json:"tag_id"`
	TagName		string	`gorm:"type:varchar(32);" json:"tag_name"`
	UserID		uint64	`gorm:" not null;" json:"user_id"`
	Username      string    `gorm:"type:varchar(32); unique" json:"username"`
	CreateTime	time.Time	`gorm:"type:datetime" json:"create_time"`
}

//标签-文章
type TagPaper struct{
	ID		uint64	`gorm:"primary_key;" json:"id"`
	TagID		uint64	`json:"tag_id"`
	TagName		string	`gorm:"type:varchar(32);" json:"tag_name"`
	PaperID		string	`gorm:"type:varchar(32);" json:"paper_id"`
	Title  string `gorm:"type:varchar(256);" json:"title"`
	Abstract string `gorm:"type:varchar(1000);" json:"paperAbstrac"`
	JournalName string `gorm:"type:varchar(256);" json:"journal_name"`
	CreateTime	time.Time	`gorm:"type:datetime" json:"create_time"`
}

// 浏览记录
type BrowsingHistory struct {
	BrowsingTime time.Time `gorm:"type:datetime" json:"browsing_time"`
	UserID       uint64    `gorm:" not null;" json:"user_id"`
	PaperID      string    `gorm:"type:varchar(32);" json:"paper_id"`
	Title  string `gorm:"type:varchar(256);not null" json:"title"`
}