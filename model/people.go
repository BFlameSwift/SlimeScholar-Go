package model

// 本文件记录 主要涉及到人的表：包括人与人的关系
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
type Author struct {
	AuthorID       string `gorm:"type:varchar(30);primary_key;" json:"author_id"`
	AuthorName     string `gorm:"type:varchar(100)" json:"author_name"`
	Affiliation    string `gorm:"type:varchar(100)",json:"affiliation"`
	PublishNumber  int    `gorm:"default:0" ,json:"publish_number"`
	CitationNumber int    `gorm:"default:0" ,json:"citation_number"`
}

type AuthorConnection struct {
	ConnectionID uint64 `gorm:"primary_key; not null" json:"connection_id"`
	AuthorID1    string `gorm:"type:varchar(30);" json:"author_id"`
}

type Followers struct {
	ID             uint64    `gorm:"primary_key; not null;" json:"id"`
	FollowUserID   uint64    `gorm:"not null" json:"follow_user_id"`
	BeFollowUserID uint64    `gorm:"not null" json:"be_follow_user_id"`
	FollowTime     time.Time `gorm:"type:datetime" json:"follow_time"`
}
type CollectPapers struct {
	UserID  uint64 `gorm:"primary_key; not null;" json:"user_id"`
	PaperID string `gorm:"type:varchar(30);primary_key;" json:"paper_id"`
}
