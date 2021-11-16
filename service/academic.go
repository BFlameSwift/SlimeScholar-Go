package service

import (
	"encoding/json"
	"gitee.com/online-publish/slime-scholar-go/model"
)

func JsonToPaper(jsonStr string) model.Paper {
	var item map[string]interface{} = make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &item)
	ok := false
	if err != nil {
		panic(err)
	}
	var paper model.Paper
	paper.Id = item["id"].(string)
	paper.Title = item["title"].(string)
	paper.Abstract,ok = item["paperAbstract"].(string)
	if !ok {paper.Abstract = ""}
	paper.Url = item["s2Url"].(string)
	paper.S2PdfUrl = item["s2PdfUrl"].(string)
	year,ok := item["year"].(float64)
	if(!ok){year = 0}
	paper.Year = int(year)
	paper.JournalPages = item["journalPages"].(string)
	paper.JournalName = item["journalName"].(string)
	paper.JournalVolume = item["journalVolume"].(string)
	paper.Doi = item["doi"].(string)
	paper.DoiUrl = item["doiUrl"].(string)
	pdf_urls  := make([]string,len(item["pdfUrls"].([]interface{})))
	for i,url := range (item["pdfUrls"].([]interface{})){
		pdf_urls[i] = url.(string)
	};paper.PdfUrls = pdf_urls
	in_citations  := make([]string,len(item["inCitations"].([]interface{})))
	for i,str := range (item["inCitations"].([]interface{})){
		in_citations[i] = str.(string)
	};paper.InCitations = in_citations
	out_citations  := make([]string,len(item["outCitations"].([]interface{})))
	for i,str := range (item["outCitations"].([]interface{})){
		out_citations[i] = str.(string)
	};paper.OutCitations = out_citations
	_,ok = item["FieldsOfStudy"].([]interface{})
	if !ok{item["FieldsOfStudy"] = make([]interface{},0)   }
	fields  := make([]string,len(item["FieldsOfStudy"].([]interface{})))
	for i,str := range (item["FieldsOfStudy"].([]interface{})){
		fields[i] = str.(string)
	};paper.FieldsOfStudy = fields
	authors := make([]model.Author,len(item["authors"].([]interface{})))
	_,ok = item["authors"].([]map[string]interface{})
	if !ok{item["authors"] = make([]map[string]interface{},0)   }
	for i,item_author := range (item["authors"].([]map[string]interface{})){
		author_new := model.Author{Id: item_author["id"].(string),Name: item_author["name"].(string)}
		authors[i] = author_new
	};paper.Authors = authors

	//author.position, ok = item["position"].(string)
	//if !ok {
	//	author.position = ""
	//}

	if err != nil {
		panic(err)
	}
	return paper
}
