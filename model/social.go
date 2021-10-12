package model

// 本文件下主要记录社交化带来的表
import "time"

type Message struct {
  MsgID   uint64  `gorm:"primary_key; not null" json:"msg_id"`
  Content  string  `gorm:"size :256;" json:"content"`
  Title   string  `gorm:"size:128" json:"title"`
  CreateTime time.Time `gorm:"type:datetime" json:"create_time"`
}

type Comment struct {
  CommentID  uint64  `gorm:"primary_key;not null" json:"comment_id"`
  Like    uint64  `gorm:"default:0" json:"like"`
  UnLike   uint64  `gorm:"default:0" json:"unlike"`
  UserID   uint64  `gorm:" not null;" json:"user_id"`
  PaperID   string  `gorm:"size:30" json:"paper_id"`
  CommentTime time.Time `gorm:"type:datetime" json:"comment_time"`
  Content   string  `gorm:"size:255" json:"content"`
  OnTop    bool   `gorm:"default:false" json:"on_top"`
}

type Like struct { // 点赞
  IsLike  bool  `gorm:"default:false" json:"is_like"`
  CommentID uint64 `gorm:"primary_key;" json:"comment_id"`
  UserID  uint64 `gorm:"primary_key;" json:"user_id"`
}
