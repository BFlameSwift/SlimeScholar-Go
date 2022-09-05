package model

// 由于相关数据量过大，数据全部会放在elasti search上，此处暂时保存一下
// 此处的数据库只是简单复述，实际的数据格式要看https://docs.microsoft.com/en-us/academic-services/graph/reference-data-schema
// 不过我的操作是将abstract paperauthor paperfield 都合并到paper这个大index中，方便聚合搜索。
type Paper struct {
	PaperId       string `gorm:"type:varchar(32); primary_key;" json:"paper_id"`
	Rank          int    `gorm:"type:integer;" json:"rank"`
	Doi           string `gorm:"type:varchar(64)" json:"doi"`
	DocType       string `gorm:"type: varchar(32)" json:"doc_type"`
	Title         string `gorm:"type: varchar(32)" json:"title"`
	BookTitle     string `gorm:"type: varchar(32)" json:"book_title"`
	Year          int    `gorm:"type: integer; " json:"year"`
	Date          string `gorm:"type: varchar(16)" json:"date"`
	JournalId     string `gorm:"type: varchar(32)" json:"journal_id"`
	ConferenceId  string `gorm:"type: varchar(32)" json:"conference_id"`
	Volume        string `gorm:"type: varchar(16)" json:" volume"`
	FirstPage     string `gorm:"type: varchar(8)" json:"first_page"`
	LastPage      string `gorm:"type:varchar(8)" json:"last_page"`
	PaperCount    int    `gorm:"type:integer" json:"paper_count"`
	CitationCount int    `gorm:"type:integer" json:"citation_count"`
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
	Location       string `gorm:"type:varchar(30)" json:"location"`
	OfficalPage    string `gorm:"type:varchar(64)" json:"offical_page"`
	PaperCount     int    `gorm:"type:int" json:"paper_count"`
	CitationCount  int    `gorm:"type:int" json:"citation_count"`
}

type Transfer struct {
	TransferID uint64 `gorm:"primary_key; not null;" json:"transfer_id"`
	PaperId    string `gorm:"type:varchar(32);" json:"paper_id"`
	AuthorId   string `gorm:"type:varchar(32);not null;" json:"author_id"`
	UserID     uint64 `gorm:"not null;" json:"user_id"`
	ObjUserID  uint64 `gorm:";" json:"obj_user_id"`
	Kind       int    `gorm:"type:int" json:"kind"`
	Status     int    `gorm:"type:int;default:0" json:"status"`
}
