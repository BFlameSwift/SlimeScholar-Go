package model

// 由于相关数据量过大，数据可能会放在elasti search上，此处暂时保存一下

type Paper struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Abstract string `json:"paperAbstrac"`
	Url string `json:"url"`
	PdfUrls []string `json:"pdf_urls"`
	S2PdfUrl string `json:"s2pdf_urls"`
	InCitations []string `json:"in_citations"`
	OutCitations []string `json:"out_citations"`
	FieldsOfStudy []string `json:"fieldsOfStudy"`
	Year int `json:"year"`
	Venue string   `json:"venue"`
	JournalName string `json:"journal_name"`
	JournalVolume string `json:"journal_volume"`
	JournalPages string `json:"journal_pages"`
	Doi string `json:"doi"`
	DoiUrl string `json:"doi_url"`
	MagId string `json:"mag_id"`
	Authors []Author `json:"authors"`
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
