package model

// 由于相关数据量过大，数据可能会放在elasti search上，此处暂时保存一下

type Paper struct {
	PaperID        string `gorm:"type:varchar(32);primary_key;" json:"paper_id"`
	Title          string `gorm:"type:varchar(256);not null" json:"title"`
	Abstract       string `gorm:"type:varchar(1024);" json:"abstract"`
	PublishYear    string `gorm:"type:varchar(5)" json:"paper_publish_year"`
	AuthorID       string `gorm:"type:varchar(30);primary_key;" json:"author_id"`
	AuthorName     string `gorm:"type:varchar(100)" json:"author_name"`
	AuthorOrg      string `gorm:"type:varchar(128)" json:"author_org"`
	Keywords       string `gorm:"type:varchar(128)" json:"keywords"` // 使用特殊字符间隔开
	CitationNumber int    `gorm:"default:0" ,json:"citation_number"`
	Publisher      string `goem:"type:varchar(32)" json:"publisher"`
	DOI            string `gorm:"type:varchar(64)" json:"doi"`
	ISBN           string `gorm:"type:varchar(32)" json:"isbn"`
	URL            string `gorm:"type:varchar(128)" json:"url"` // 是否需要限定大小
	PDF            string `gorm:"type:varchar(128" json:"pdf"`
	ConferenceID   string `gorm:"type:varchar(30)" json:"conference_id"`
	Lang           string `gorm:"type:varchar(32)" json:"language"`
	Venue          string `gorm:"type:varchar(64) "json:"venue"` // 领域 数据中看到的

	// TODO
}

type PaperReference struct {
	PaperID             string `gorm:"type:varchar(32);primary_key;" json:"paper_id"`
	ReferencePaperID    string `gorm:"type:varchar(32);primary_key;" json:"reference_paper_id"`
	PaperTitle          string `gorm:"type:varchar(256);not null" json:"paper_title"`
	ReferencePaperTitle string `gorm:"type:varchar(256);not null" json:"reference_paper_title"`
}

type Conference struct {
	ConferenceID   string `gorm:"type:varchar(30);primary_key;" json:"conference_id"`
	ConferenceName string `gorm:"type:varchar(30)" json:"conference_name"`
}
