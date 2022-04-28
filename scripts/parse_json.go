package scripts

import (
	"bufio"
	"encoding/json"
	"fmt"

	// "io"
	"os"
	"strconv"
	"strings"
	"time"
)

var cspaperMap map[int64]byte = make(map[int64]byte)
var closureMap map[int64]byte = make(map[int64]byte)

func loadMap() {
	fileCspaper, err := os.Open("cspaper_ids.txt")
	if err != nil {
		fmt.Println(err)
	}
	fileClosure, err := os.Open("cspaperclosure_ids.txt")
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(fileClosure)
	for scanner.Scan() {
		lineInt, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			panic(err)
		}
		closureMap[lineInt] = 1
	}
	scanner = bufio.NewScanner(fileCspaper)
	for scanner.Scan() {
		lineInt, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			panic(err)
		}
		closureMap[lineInt] = 1
	}
	// fmt.Println("load end ", time.Now(), "load sum", getKeys1(closure_map))
}
func makePaper(li []string) map[string]interface{} {
	paperId, rank, doi, doctype, _, title, bookTitle, year, date, _, publisher, journalId, conferenceId, _, volume, _, firstPage, lastPage, referenceCount, citationCount, _, _, _, _, _, _ :=

		li[0], li[1], li[2], li[3], li[4], li[5], li[6], li[7], li[8], li[9], li[10], li[11], li[12], li[13], li[14], li[15], li[16], li[17], li[18], li[19], li[20], li[21], li[22], li[23], li[24], li[25]
	paperMap := make(map[string]interface{})
	paperMap["paper_id"] = paperId
	paperMap["rank"] = rank
	paperMap["doi"] = doi
	paperMap["doctype"] = doctype
	paperMap["title"] = title
	paperMap["book_title"] = bookTitle
	paperMap["year"] = year
	paperMap["date"] = date
	paperMap["publisher"] = publisher
	paperMap["journal_id"] = journalId
	paperMap["conference_id"] = conferenceId
	paperMap["volume"] = volume
	paperMap["first_page"] = firstPage
	paperMap["last_page"] = lastPage
	paperMap["reference_count"] = referenceCount
	paperMap["citation_count"] = citationCount

	return paperMap
}
func makeAuthor(li []string) map[string]interface{} {
	authorId, rank, _, name, affiliationId, paperCount, _, citationCount, _ :=
		li[0], li[1], li[2], li[3], li[4], li[5], li[6], li[7], li[8]
	authorMap := make(map[string]interface{})
	authorMap["author_id"] = authorId
	authorMap["rank"] = rank
	authorMap["affiliation_id"] = affiliationId
	authorMap["name"] = name
	authorMap["paper_count"] = paperCount
	authorMap["citation_count"] = citationCount
	return authorMap
}
func makeConference(li []string) map[string]interface{} {
	conferenceId, _, name, _, location, officalUrl, start, end, _, _, _, _, paperCount, _, citationCount, _, _, _ :=
		li[0], li[1], li[2], li[3], li[4], li[5], li[6], li[7], li[8], li[9], li[10], li[11], li[12], li[13], li[14], li[15], li[16], li[17]
	conferenceMap := make(map[string]interface{})
	conferenceMap["conference_id"] = conferenceId
	conferenceMap["name"] = name
	conferenceMap["location"] = location
	conferenceMap["offical_url"] = officalUrl
	conferenceMap["start"] = start
	conferenceMap["end"] = end
	conferenceMap["paper_count"] = paperCount
	conferenceMap["citation_count"] = citationCount
	return conferenceMap
}
func makeField(li []string) map[string]interface{} {
	fieldId, rank, _, name, mainType, level, paperCount, _, citationCount, _ :=
		li[0], li[1], li[2], li[3], li[4], li[5], li[6], li[7], li[8], li[9]
	fieldMap := make(map[string]interface{})
	fieldMap["field_id"] = fieldId
	fieldMap["rank"] = rank
	fieldMap["name"] = name
	fieldMap["main_type"] = mainType
	fieldMap["level"] = level
	fieldMap["paper_count"] = paperCount
	fieldMap["citation_count"] = citationCount
	return fieldMap
}
func makeJournal(li []string) map[string]interface{} {
	journalId, rank, _, name, issn, publisher, webpage, paperCount, _, citationCount, _ :=
		li[0], li[1], li[2], li[3], li[4], li[5], li[6], li[7], li[8], li[9], li[10]
	journalMap := make(map[string]interface{})
	journalMap["journalid"] = journalId
	journalMap["rank"] = rank
	journalMap["name"] = name
	journalMap["issn"] = issn
	journalMap["publisher"] = publisher
	journalMap["webpage"] = webpage
	journalMap["paper_count"] = paperCount
	journalMap["citation_count"] = citationCount
	return journalMap
}
func makeRelPaperAuthor(li []string) map[string]interface{} {
	paperId, authorId, affiliationId, sequence, authorName, affiliationName :=
		li[0], li[1], li[2], li[3], li[4], li[5]
	relMap := make(map[string]interface{})
	relMap["pid"] = paperId
	relMap["aid"] = authorId
	relMap["afid"] = affiliationId
	relMap["order"] = sequence
	relMap["aname"] = authorName
	relMap["afname"] = affiliationName
	return relMap
}
func makeRelPaperPaper(li []string, otherType string) map[string]interface{} {
	paperId, referenceId := li[0], li[1]
	relMap := make(map[string]interface{})
	relMap["pid"] = paperId
	relMap[otherType] = referenceId
	return relMap
}

func makeMultiMap(lines []string, mytype string, mainType, mainId string, otherType string) map[string]interface{} {
	multiMap := make(map[string]interface{})
	multiMap[mainType] = mainId
	other := make([]interface{}, 0, 10010)
	for _, line := range lines {
		lineList := strings.Split(line, "\t")
		signalMap := make(map[string]interface{})
		if mytype == "rel_paper_author" {
			signalMap = makeRelPaperAuthor(lineList)
			delete(signalMap, "pid")
		} else if mytype == "rel_paper_reference" || mytype == "rel_paper_field" {
			signalMap = makeRelPaperPaper(lineList, otherType)
			delete(signalMap, "pid")
		}
		if len(signalMap) == 1 {
			other = append(other, signalMap[otherType])
			continue
		}
		other = append(other, signalMap)
	}
	multiMap[otherType] = other
	return multiMap
}

func makeMap(line string, mytype string) map[string]interface{} {
	lineList := strings.Split(line, "\t")

	if mytype == "paper" {
		return makePaper(lineList)
	} else if mytype == "author" {
		return makeAuthor(lineList)
	} else {
		return make(map[string]interface{})
	}
}
func writeToFile(jsonListPointer *[]string, writeObj *bufio.Writer) {
	jsonList := *jsonListPointer
	for _, str := range jsonList {
		_, err := writeObj.Write([]byte(str + "\n"))
		if err != nil {
			panic(err)
		}
	}
}

func makeJson(filename string, mytype string, isMulti bool, mainType string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)
	i := 0

	writeFile, err := os.OpenFile("my"+filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	writeObj := bufio.NewWriterSize(writeFile, 4096)
	if err != nil {
		panic(err)
	}
	jsonList := make([]string, 0, 3000000)
	theMap := make(map[string]interface{})
	lineNum := 0
	thisId := "bababababbabab"
	lines := make([]string, 0, 5000)
	for scanner.Scan() {
		line := ""
		lineNum += 1
		if isMulti {

			for scanner.Scan() {
				line = scanner.Text()
				lineList := strings.Split(line, "\t")
				if thisId == "bababababbabab" || thisId == lineList[0] {
					lines = append(lines, line)
					thisId = lineList[0]
				} else {
					theMap = makeMultiMap(lines, mytype, mainType, thisId, "rel")
					lines = make([]string, 0, 5000)
					thisId = lineList[0]
					lines = append(lines, line)
					break
				}
			}
		} else {
			line = scanner.Text()
			theMap = makeMap(line, mytype)
		}

		if theMap != nil {

			jsonStr, err := json.Marshal(theMap)
			if err != nil {
				panic(err)
			}
			i += 1
			jsonList = append(jsonList, string(jsonStr))
			if i%1000000 == 0 {
				if err != nil {
					fmt.Println("cacheFileList.yml file create Failed. err: " + err.Error())
				}
				writeToFile(&jsonList, writeObj) // 传递数组为拷贝传递，穿指针省时间
				jsonList = make([]string, 0, 3000000)
				fmt.Println(time.Now(), i, lineNum)
			}

		}

	}
	writeToFile(&jsonList, writeObj)
	jsonList = make([]string, 0, 1000000)
	fmt.Println(time.Now(), i, lineNum)
}

func makeMultiRel(filename string, mytype string, isMulti bool, mainType string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)
	i := 0

	writeFile, err := os.OpenFile("my"+filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	writeObj := bufio.NewWriterSize(writeFile, 4096)
	if err != nil {
		panic(err)
	}
	jsonList := make([]string, 0, 3000000)
	theMap := make(map[string]interface{})
	lineNum := 0
	thisId := "bababababbabab"
	lines := make([]string, 0, 5000)
	for scanner.Scan() {
		line := ""
		lineNum += 1

		line = scanner.Text()
		lineList := strings.Split(line, "\t")
		if thisId == "bababababbabab" || thisId == lineList[0] {
			lines = append(lines, line)
			thisId = lineList[0]
		} else {
			theMap = makeMultiMap(lines, mytype, mainType, thisId, "rel")
			lines = make([]string, 0, 5000)
			thisId = lineList[0]
			lines = append(lines, line)
			if theMap != nil {

				jsonStr, err := json.Marshal(theMap)
				if err != nil {
					panic(err)
				}
				i += 1
				jsonList = append(jsonList, string(jsonStr))
				if i%1000000 == 0 {
					if err != nil {
						fmt.Println("cacheFileList.yml file create Failed. err: " + err.Error())
					}
					writeToFile(&jsonList, writeObj) // 传递数组为拷贝传递，穿指针省时间
					jsonList = make([]string, 0, 3000000)
					fmt.Println(time.Now(), i, lineNum)
				}

			}

		}

	}
	writeToFile(&jsonList, writeObj)
	jsonList = make([]string, 0, 1000000)
	fmt.Println(time.Now(), i, lineNum)
}

func parseJsonMain() {
	s := time.Now()
	fmt.Println(s)
	//load_map()

	// make_json("Papers.txt", "paper",false,"paper_id")
	fmt.Println(time.Now(), "load paper end")
	// make_json("Authors.txt", "author",false,"author_id")	make_multi_rel("PaperAuthorAffiliations.txt", "rel_paper_author", true, "paper_id")//2021-11-26 22:51:14.0211883 +0800 CST m=+11773.545295501 269412162 731587019
	makeMultiRel("PaperReferences.txt", "rel_paper_reference", true, "paper_id")
	fmt.Println(time.Now())
	makeMultiRel("PaperFieldsOfStudy.txt", "rel_paper_field", true, "paper_id")
	makeMultiRel("PaperCitationContexts.txt", "rel_paper_citation", true, "paper_id")
}
