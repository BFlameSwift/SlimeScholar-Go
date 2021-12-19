package service

import (
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"sort"
	"strings"
	"time"

	"gitee.com/online-publish/slime-scholar-go/global"
	"gitee.com/online-publish/slime-scholar-go/model"
	"gorm.io/gorm"
)

// 将paper简化到最贱格式
func MostSimplifyPaper(m map[string]interface{}) (ret map[string]interface{}) {
	ret["paper_id"] = m["paper_id"]
	ret["paper_title"] = m["paper_title"]
	return ret
}

// SimplifyPapers 检查map形式的paper列表
func SimplifyPapers(inter []interface{}) []interface{} {
	ret_list := make([]interface{}, len(inter))
	for _, v := range inter {
		ret_list = append(ret_list, MostSimplifyPaper(v.(map[string]interface{})))
	}
	return ret_list
}

// BrowerPaper 浏览记录保存
func BrowerPaper(paper map[string]interface{}) {
	title := paper["paper_title"].(string)
	authors := paper["authors"].([]interface{})
	paper_id := paper["paper_id"].(string)
	authors_name := ""
	for i, author := range authors {
		authors_name += author.(map[string]interface{})["name"].(string)
		if i < len(authors)-1 {
			authors_name += ", "
		}
	}
	browsing_history := model.BrowsingHistory{BrowsingTime: time.Now(), Title: title, Authors: authors_name, PaperID: paper_id}
	err := CreateBrowseHistory(&browsing_history)
	if err != nil {
		panic(err)
	}
}

// FindExistingTransfer 精准查询已经存在的transfer
func FindExistingTransfer(author_id string, paper_id string, user_id uint64, kind int) (transfer *model.Transfer, notFound bool) {
	err := global.DB.Where("author_id = ? AND paper_id = ? AND user_id = ? AND kind = ?", author_id, paper_id, user_id, kind).First(&transfer).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return transfer, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return transfer, false
	}
}

// TransferPaper 创建transfer
func TransferPaper(user model.User, author_id string, paper_id string, kind int, obj_user_id uint64) {
	_, notFound := FindExistingTransfer(author_id, paper_id, user.UserID, kind)
	if notFound {
		transfer := model.Transfer{UserID: user.UserID, AuthorId: author_id, PaperId: paper_id, Kind: kind, Status: 1, ObjUserID: obj_user_id}
		if err := global.DB.Create(&transfer).Error; err != nil {
			panic(err)
		}
	}
}

// FindAllAuthorManagePapers 根据作者id找到作者的所有transfer
func FindAllAuthorManagePapers(author_id string) (transfer_list *[]model.Transfer, notFound bool) {
	err := global.DB.Where("author_id = ?", author_id).Find(&transfer_list).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return transfer_list, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return transfer_list, false
	}
}

// PaperMapToPaperList 将papermap转化为paper————list
func PaperMapToPaperList(m map[string]interface{}) (ret_list []interface{}) {
	for key := range m {
		item := m[key]
		author_map := make(map[string]interface{})
		author_map["rel"] = item.(map[string]interface{})["authors"]
		item.(map[string]interface{})["authors"] = ParseRelPaperAuthor(author_map)["rel"]
		ret_list = append(ret_list, item)
	}
	return ret_list
}
func GetAuthorAllPapersIds(author_id string) []string {
	paper_result := QueryByField("paper", "authors.aid.keyword", author_id, 1, 100)
	paper_ids_origin := make([]string, 0, 10000)
	//authors_map := make(map[string]interface{})
	for _, hit := range paper_result.Hits.Hits {
		hit_map := make(map[string]interface{})
		err := json.Unmarshal([]byte(hit.Source), &hit_map)
		if err != nil {
			panic(err)
		}
		paper_ids_origin = append(paper_ids_origin, hit_map["paper_id"].(string))
	}
	transfer_list, notFound := FindAllAuthorManagePapers(author_id)
	papre_ids_del := make([]string, 0, 1000)
	// 找到应该删除的paper和应该添加的paperids
	if !notFound || len(*transfer_list) != 0 {
		for _, transfer := range *transfer_list {
			if transfer.Status == 1 {
				if transfer.Kind == 2 || transfer.Kind == 1 {
					papre_ids_del = append(papre_ids_del, transfer.PaperId)
				} else if transfer.Kind == 0 {
					paper_ids_origin = append(paper_ids_origin, transfer.PaperId)
				}
			}
		}
	}
	// 去重与删除操作
	paper_ids_final := make([]string, 0, 1000)
	paper_ids_map := make(map[string]int)
	for _, id := range paper_ids_origin {
		paper_ids_map[id] = 1
	}
	for _, id := range papre_ids_del {
		paper_ids_map[id] = 0
	}
	for key := range paper_ids_map {
		if paper_ids_map[key] == 1 {
			paper_ids_final = append(paper_ids_final, key)
		}
	}
	return paper_ids_final
}

// GetAuthorAllPaper 根据作者id获取该作者所有的papers
func GetAuthorAllPaper(author_id string) (paper_list []interface{}) {

	return GetPapers(GetAuthorAllPapersIds(author_id))
}

// JudgeAuthorIsSettled 判断作者是否已经入驻
func JudgeAuthorIsSettled(author_id string) (bool, uint64) {
	submit, notFound := QueryASubmitByAuthor(author_id)
	return !notFound, submit.UserID
}

// GetAuthorMsg 未入驻作者在展示个人中心之前的格式转化
func GetAuthorMsg(author_id string) (author_map map[string]interface{}) {
	author_json := GetsByIndexIdWithout("author", author_id)
	if author_json == nil {
		return nil
	}
	err := json.Unmarshal(author_json.Source, &author_map)
	if err != nil {
		panic(err)
	}
	if author_map["affiliation_id"].(string) != "" {
		affiliation_byte := GetsByIndexIdWithout("affiliation", author_map["affiliation_id"].(string)).Source
		affiliation_map := make(map[string]interface{})
		err = json.Unmarshal([]byte(affiliation_byte), &affiliation_map)
		author_map["affiliation"] = affiliation_map["name"]
		author_map["affiliation_entry"] = affiliation_map

	} else {
		author_map["affiliation"] = ""
	}
	author_map["author_name"] = author_map["name"]
	delete(author_map, "name")
	return author_map
}

// ProcAuthorMsg 处理作者的基本信息：生成作者的领域等等
func ProcAuthorMsg(people map[string]interface{}, papers []interface{}) map[string]interface{} {
	fields_map := make(map[string]int)
	for _, paper := range papers {
		if paper.(map[string]interface{})["fields"] != nil {
			for _, field := range paper.(map[string]interface{})["fields"].([]interface{}) {
				//fmt.Println(field)
				fieldStr := field.(map[string]interface{})["field_id"].(string)
				if _, ok := fields_map[fieldStr]; ok {
					fields_map[fieldStr]++
				} else {
					fields_map[fieldStr] = 1
				}
			}
		}
	}
	fields_items := IdsGetItems(GetTopNKey(fields_map, 5), "fields")
	fields := make([]string, 0)
	for id := range fields_items {
		fields = append(fields, fields_items[id].(map[string]interface{})["name"].(string))
	}
	people["fields"] = fields
	people["fields_graph"] = getFieldsMap(fields_map)
	return people
}

// 根据上面的函数来获取领域列表
func getFieldsMap(m map[string]int) (ret []interface{}) {
	keys := GetAllSortedKey(m)
	fields_items := IdsGetItems(keys, "fields")

	for _, id := range keys {
		item := fields_items[id].(map[string]interface{})
		item["count"] = m[id]
		ret = append(ret, item)
	}
	return ret
}

// SortPapers 对paper按照年份引用次数进行排序
func SortPapers(papers []interface{}, sort_type int) []interface{} {
	inter := papers
	sort.Slice(inter, func(i, j int) bool {
		if sort_type == 1 {
			return inter[i].(map[string]interface{})["year"].(string) < inter[j].(map[string]interface{})["year"].(string)
		} else if sort_type == 2 {
			return inter[i].(map[string]interface{})["year"].(string) >= inter[j].(map[string]interface{})["year"].(string)
		} else if sort_type == 3 {
			return inter[i].(map[string]interface{})["citation_count"].(float64) < inter[j].(map[string]interface{})["citation_count"].(float64)
		} else if sort_type == 4 {
			return inter[i].(map[string]interface{})["citation_count"].(float64) >= inter[j].(map[string]interface{})["citation_count"].(float64)
		} else {
			return true
		}
	})
	return inter
}

// GetPaperAuthorsName 获取paper的所有作者组成的列表
func GetPaperAuthorsName(paper map[string]interface{}) (ret []string) {

	for _, author := range paper["authors"].([]interface{}) {

		ret = append(ret, author.(map[string]interface{})["author_name"].(string))
	}
	return ret
}

// GetPaperPages 根据引用文献，获取论文的起止页数，来格式化输出
func GetPaperPages(paper map[string]interface{}) string {
	ret := ""
	if paper["first_page"].(string) != "" {
		ret += ":" + paper["first_page"].(string)
		if paper["last_page"].(string) != "" {
			ret += "-" + paper["last_page"].(string)
		}
		ret += "."
	}
	return ret
}

// GetPaperCiteType 获取引用文献的最后一段
func GetPaperCiteType(paper map[string]interface{}) string {
	doctype := paper["doctype"].(string)
	if doctype == "Patent" {
		return "[P]."
	} else if doctype == "Conference" {
		return "[C]."
	} else if doctype == "Journal" {
		journal, err := GetsByIndexId("journal", paper["journal_id"].(string))
		journalName := ""
		if err == nil {
			journalMap := make(map[string]interface{})
			_ = json.Unmarshal(journal.Source, &journalMap)
			journalName += journalMap["name"].(string)
		}
		return "[J]." + journalName + "," + paper["year"].(string) + "," + paper["volume"].(string) + GetPaperPages(paper)
	} else {
		return "[M]" + "." + paper["publisher"].(string) + "," + paper["year"].(string) + GetPaperPages(paper)
	}
	//else if doctype == "BookChapter" || doctype == "Book" {
	//	return "M"
	//}
}
func FormatCite(rank int, name string, content string) map[string]interface{} {
	ret := make(map[string]interface{})
	ret["id"] = rank
	ret["name"] = name
	ret["content"] = content
	return ret
}

// 根据paper生成MLA格式的参考文献
func MLACitePaper(paper map[string]interface{}) (ret string) {
	//paper := GetSimplePaper(paperId)
	authors := GetPaperAuthorsName(paper)
	if len(authors) > 3 {
		ret += authors[0] + " et al"
	} else {
		ret += strings.Join(GetPaperAuthorsName(paper), ",")
	}
	ret += "." + "\"" + paper["paper_title"].(string) + "\""
	if paper["journal_id"].(string) != "" {
		journal := GetsByIndexIdRetMap("journal", paper["journal_id"].(string))
		ret += " " + journal["name"].(string)
	} else if paper["conference_id"].(string) != "" {
		conference := GetsByIndexIdRetMap("conference", paper["conference_id"].(string))
		ret += " " + conference["name"].(string)
	}
	ret += " (" + paper["year"].(string) + ")"
	return ret
}
func APACitePaper(paper map[string]interface{}) (ret string) {
	ret += strings.Join(GetPaperAuthorsName(paper), ".& ")
	ret += " (" + paper["year"].(string) + ")."
	ret += paper["paper_title"].(string) + "."
	if paper["journal_id"].(string) != "" {
		journal := GetsByIndexIdRetMap("journal", paper["journal_id"].(string))
		ret += " " + journal["name"].(string) + "Journal," + GetPaperPages(paper)
	} else if paper["conference_id"].(string) != "" {
		conference := GetsByIndexIdRetMap("conference", paper["conference_id"].(string))
		ret += " " + conference["name"].(string)
	}
	ret += "."
	return ret
}

// CitePaper 根据paperid引用文献
func CitePaper(paperId string) (ret []interface{}) {
	paper := GetSimplePaper(paperId)
	authors := strings.Join(GetPaperAuthorsName(paper), ",")
	//fmt.Println(authors)
	title := paper["paper_title"].(string)
	//fmt.Println(title)
	citedType := GetPaperCiteType(paper)
	//fmt.Println(citedType)
	ret = append(ret, FormatCite(1, "GB/T 7714", authors+title+citedType))
	ret = append(ret, FormatCite(2, "MLA", MLACitePaper(paper)))
	ret = append(ret, FormatCite(3, "APA", APACitePaper(paper)))
	//fmt.Println(gbt["GB/T 7714"])
	return ret
}

func GetPaperCitationIds(paperIds []string, size int, page int) ([]string, int) {
	boolQuery, idsQuery := elastic.NewBoolQuery(), elastic.NewBoolQuery()
	for _, id := range paperIds {
		idsQuery.Should(elastic.NewMatchPhraseQuery("rel.keyword", id))
	}
	boolQuery.Must(idsQuery)

	searchResult, err := Client.Search().Index("reference").Query(boolQuery).From((page - 1) * size).Size(size).Do(context.Background())
	if err != nil {
		panic(err)
	}
	citationsIds := make([]string, 0)

	for _, hit := range searchResult.Hits.Hits {
		citationsIds = append(citationsIds, hit.Id)
	}
	return citationsIds, int(searchResult.TotalHits())
}

// 根据文献id获取引用此文献的文献的引用图
func GetCitationPapersGraph(paperIds []string, size int) ([]string, []int) {
	citationsIds, _ := GetPaperCitationIds(paperIds, size, 1)
	mulIdsQuery := elastic.NewIdsQuery().Ids(citationsIds...)
	//fmt.Println("citation_count:!!!!!", len(citationsIds))
	yearAggregation := elastic.NewTermsAggregation().Field("year.keyword")
	searchResult, err := Client.Search().Index("paper").Query(mulIdsQuery).Size(0).Aggregation("year", yearAggregation).Do(context.Background())
	if err != nil {
		panic(err)
	}
	agg, found := searchResult.Aggregations.Terms("year")

	if !found {
		return make([]string, 0), make([]int, 0)
	}
	yearList, citationCountList := make([]string, 0), make([]int, 0)
	citationMap := make(map[string]int)
	for _, bucket := range agg.Buckets {
		if bucket.Key.(string) == "" {
			continue
		}
		citationMap[bucket.Key.(string)] = int(bucket.DocCount)
		yearList = append(yearList, bucket.Key.(string))

	}
	sort.Strings(yearList)
	for _, year := range yearList {
		citationCountList = append(citationCountList, citationMap[year])
	}
	return yearList, citationCountList
}
