package model

import (
	"database/sql"
	"time"
)

// 用户
type User struct {
	UserID        uint64    `gorm:"primary_key; not null;" json:"user_id"`
	Username      string    `gorm:"size:32; not null; unique" json:"username"`
	Password      string    `gorm:"size:32; not null" json:"password"`
	UserInfo      string    `gorm:"size:255;" json:"user_info"`
	UserType      uint64    `gorm:"default:0" json:"user_type"` // 0: 普通用户，1: 认证机构用户,2 管理员
	Affiliation   string    `gorm:"size:64;" json:"affiliation"`
	AuthorName    string    `gorm:"size:64;" json:"author_name"`
	HomePage      string    `gorm:"size:64;" json:"home_page"`
	Email         string    `gorm:"size:32;" json:"email"`
	WorkEmail     string    `gorm:"size:32;" json:"work_email"`
	Fields        string    `gorm:"size:256;" json:"fields"`
	HasConfirmed  bool      `gorm:"default:false" json:"has_confirmed"`
	ConfirmNumber int       `gorm:"default:0" json:"confirm_number"`
	RegTime       time.Time `gorm:"column:reg_time;type:datetime" json:"reg_time"`
	PaperCount    int       `gorm:"default:0" json:"paper_count"`
}
type Author struct {
	AuthorId            string `gorm:"primary_key; not null;" json:"author_id"`
	Rank                string `gorm:"size :32;" json:"rank"`
	AuthorName          string `gorm:"size 64" json:"author_name"`
	AuthorAffiliationId string `gorm:"size :32;" json:"author_affiliation_id"`
	PaperCount          int    `gorm:"type:int" json:"paper_count"`
	CitationCount       int    `gorm:"type:int" json:"citation_count"`
}

type Affiliation struct {
	AffiliationName string `gorm:"type:varchar(150)" json:"affiliation_name"`
	AffiliationID   string `gorm:"type:varchar(32);primary_key" json:"affiliation_id"`
	OffcialPage     string `gorm:"type:varchar(86)" json:"offcial_page"`
	PaperCount      int    `gorm:"type:integer" json:"paper_count"`
	CitationCount   int    `gorm:"type:integer" json:"citation_count"`
}

type AuthorConnection struct {
	ConnectionID uint64 `gorm:"primary_key; not null" json:"connection_id"`
	AuthorID1    string `gorm:"type:varchar(32);" json:"author_id1"`
	AuthorID2    string `gorm:"type:varchar(32)" json:"author_id2"`
}

// TOOD 申请成为认证学者的申请表
type SubmitScholar struct {
	SubmitID        uint64       `gorm:"primary_key; not null" json:"submit_id"`
	UserID          uint64       `gorm:"not null;" json:"user_id"`
	AuthorName      string       `gorm:"not null;type:varchar(32)" json:"real_name"`
	Status          int          `gorm:"default:0" json:"status"`                    // 0:未处理，1，同意申请，2拒绝申请
	Content         string       `gorm:"type:varchar(256)" json:"content"`           // 填写内容
	AuthorID        string       `gorm:"type:varchar(32);not null" json:"author_id"` // 被申请的作者ID
	Fields          string       `gorm:"type:varchar(256);" json:"fields"`
	HomePage        string       `gorm:"type:varchar(64);" json:"home_page"`
	WorkEmail       string       `gorm:"type:varchar(64)" json:"work_email"`
	AffiliationName string       `gorm:"type:varchar(64)" json:"affiliation_name"`
	CreatedTime     time.Time    `gorm:"column:reg_time;type:datetime" json:"created_time"`
	AcceptTime      sql.NullTime `gorm:"column:accept_time;type:datetime" json:"accept_time"`
	PaperCount      int          `gorm:"type:integer;null" json:"paper_count"`
}
