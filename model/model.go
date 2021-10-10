package model

import "time"

// 用户
type User struct {
	UserID        uint64    `gorm:"primary_key; not null;" json:"user_id"`
	Username      string    `gorm:"size:25; not null; unique" json:"username"`
	Password      string    `gorm:"size:25; not null" json:"password"`
	UserInfo      string    `gorm:"size:255;" json:"user_info"`
	UserType      uint64    `gorm:"default:0" json:"user_type"` // 0: 普通用户，1: 认证机构用户,2 管理员
	Affiliation   string    `gorm:"size:25;" json:"affiliation"`
	Email         string    `gorm:"size:50;" json:"email"`
	HasComfirmed  bool      `gorm:"default:false" json:"has_comfirmed"`
	ConfirmNumber int       `gorm:"default:0" json:"confirm_number"`
	RegTime       time.Time `gorm:"column:reg_time;type:datetime" json:"reg_time"`
}
type Message struct {
	MsgID      uint64    `gorm:"primary_key; not null" json:"msg_id"`
	Content    string    `gorm:"size :256;" json:"content"`
	Title      string    `gorm:"size:128" json:"title"`
	CreateTime time.Time `gorm:"type:datetime" json:"create_time"`
}
type Followers struct {
	ID             uint64    `gorm:"primary_key; not null;" json:"id"`
	FollowUserID   uint64    `gorm:"not null" json:"follow_user_id"`
	BeFollowUserID uint64    `gorm:"not null" json:"be_follow_user_id"`
	FollowTime     time.Time `gorm:"type:datetime" json:"follow_time"`
}
type Paper struct {
	PaperID     uint64 `gorm:"primary_key;not null" json:"paper_id"`
	Title       string `gorm:"type:string" json:"title"`
	PublishYear string `gorm:"type:varchar(5)" json:"paper_publish_year"`
	// TODO
}
type Comment struct {
	CommentID uint64 `gorm:"primary_key;not null" json:"comment_id"`
	Like      uint64 `gorm:"default:0" json:"like"`
	UnLike    uint64 `gorm:"default:0" json:"unlike"`
	// TODO
}
