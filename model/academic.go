package model

type Paper struct {
	PaperID          string `gorm:"type:varchar(30);primary_key;" json:"paper_id"`
	Title       string `gorm:"type:string" json:"title"`
	PublishYear string `gorm:"type:varchar(5)" json:"paper_publish_year"`
	AuthorID   string `gorm:"type:varchar(30);primary_key;" json:"author_id"`
	AuthorName string `gorm:"type:varchar(100)" json:"author_name"`
	Keywords string `gorm:"type:varchar(100)",json:"keywords"` // 使用特殊字符间隔开
	CitationNumber int `gorm:"default:0" ,json:"citation_number"`
	DOI string `gorm:"type:varchar(50)" json:"doi"`
	URL string `gorm:"" json:"url"`

	// TODO
}

