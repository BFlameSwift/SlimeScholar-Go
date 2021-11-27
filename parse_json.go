package main

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

var cspaper_map map[int64]byte = make(map[int64]byte)
var closure_map map[int64]byte = make(map[int64]byte)

func getKeys1(m map[int64]byte) int {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	j := 0
	for _ = range m {
		j++
	}
	return j
}

func load_map() {
	file_cspaper, err := os.Open("cspaper_ids.txt")
	if err != nil {
		fmt.Println(err)
	}
	file_closure, err := os.Open("cspaperclosure_ids.txt")
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file_closure)
	for scanner.Scan() {
		line_int, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			panic(err)
		}
		closure_map[line_int] = 1
	}
	scanner = bufio.NewScanner(file_cspaper)
	for scanner.Scan() {
		line_int, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			panic(err)
		}
		closure_map[line_int] = 1
	}
	// fmt.Println("load end ", time.Now(), "load sum", getKeys1(closure_map))
}
func make_paper(li []string) map[string]interface{} {
	paper_id, rank, doi, doctype, _, title, book_title, year, date, _, publisher, journal_id, conference_id, _, volume, _, first_page, last_page, reference_count, citation_count, _, _, _, _, _, _ :=

		li[0], li[1], li[2], li[3], li[4], li[5], li[6], li[7], li[8], li[9], li[10], li[11], li[12], li[13], li[14], li[15], li[16], li[17], li[18], li[19], li[20], li[21], li[22], li[23], li[24], li[25]
	paper_map := make(map[string]interface{})
	paper_map["paper_id"] = paper_id
	paper_map["rank"] = rank
	paper_map["doi"] = doi
	paper_map["doctype"] = doctype
	paper_map["title"] = title
	paper_map["book_title"] = book_title
	paper_map["year"] = year
	paper_map["date"] = date
	paper_map["publisher"] = publisher
	paper_map["journal_id"] = journal_id
	paper_map["conference_id"] = conference_id
	paper_map["volume"] = volume
	paper_map["first_page"] = first_page
	paper_map["last_page"] = last_page
	paper_map["reference_count"] = reference_count
	paper_map["citation_count"] = citation_count

	return paper_map
}
func make_author(li []string) map[string]interface{} {
	author_id, rank, _, name, affiliation_id, paper_count, _, citation_count, _ :=
		li[0], li[1], li[2], li[3], li[4], li[5], li[6], li[7], li[8]
	author_map := make(map[string]interface{})
	author_map["author_id"] = author_id
	author_map["rank"] = rank
	author_map["affiliation_id"] = affiliation_id
	author_map["name"] = name
	author_map["paper_count"] = paper_count
	author_map["citation_count"] = citation_count
	return author_map
}
func make_conference(li []string) map[string]interface{} {
	conference_id, _, name, _, location, offical_url, start, end, _, _, _, _, paper_count, _, citation_count, _, _, _ :=
		li[0], li[1], li[2], li[3], li[4], li[5], li[6], li[7], li[8], li[9], li[10], li[11], li[12], li[13], li[14], li[15], li[16], li[17]
	conference_map := make(map[string]interface{})
	conference_map["conference_id"] = conference_id
	conference_map["name"] = name
	conference_map["location"] = location
	conference_map["offical_url"] = offical_url
	conference_map["start"] = start
	conference_map["end"] = end
	conference_map["paper_count"] = paper_count
	conference_map["citation_count"] = citation_count
	return conference_map
}
func make_field(li []string) map[string]interface{} {
	field_id, rank, _, name, main_type, level, paper_count, _, citation_count, _ :=
		li[0], li[1], li[2], li[3], li[4], li[5], li[6], li[7], li[8], li[9]
	field_map := make(map[string]interface{})
	field_map["field_id"] = field_id
	field_map["rank"] = rank
	field_map["name"] = name
	field_map["main_type"] = main_type
	field_map["level"] = level
	field_map["paper_count"] = paper_count
	field_map["citation_count"] = citation_count
	return field_map
}
func make_journal(li []string) map[string]interface{} {
	journal_id, rank, _, name, issn, publisher, webpage, paper_count, _, citation_count, _ :=
		li[0], li[1], li[2], li[3], li[4], li[5], li[6], li[7], li[8], li[9], li[10]
	journal_map := make(map[string]interface{})
	journal_map["journalid"] = journal_id
	journal_map["rank"] = rank
	journal_map["name"] = name
	journal_map["issn"] = issn
	journal_map["publisher"] = publisher
	journal_map["webpage"] = webpage
	journal_map["paper_count"] = paper_count
	journal_map["citation_count"] = citation_count
	return journal_map
}
func make_rel_paper_author(li []string) map[string]interface{} {
	paper_id, author_id, affiliation_id, sequence, author_name, affiliation_name :=
		li[0], li[1], li[2], li[3], li[4], li[5]
	rel_map := make(map[string]interface{})
	rel_map["pid"] = paper_id
	rel_map["aid"] = author_id
	rel_map["afid"] = affiliation_id
	rel_map["order"] = sequence
	rel_map["aname"] = author_name
	rel_map["afname"] = affiliation_name
	return rel_map
}
func make_rel_paper_paper(li []string, other_type string) map[string]interface{} {
	paper_id, reference_id := li[0], li[1]
	rel_map := make(map[string]interface{})
	rel_map["pid"] = paper_id
	rel_map[other_type] = reference_id
	return rel_map
}

func make_multi_map(lines []string, mytype string, main_type, main_id string, other_type string) map[string]interface{} {
	multi_map := make(map[string]interface{})
	multi_map[main_type] = main_id
	other := make([]interface{}, 0, 10010)
	for _, line := range lines {
		line_list := strings.Split(line, "\t")
		signal_map := make(map[string]interface{})
		if mytype == "rel_paper_author" {
			signal_map = make_rel_paper_author(line_list)
			delete(signal_map, "pid")
		} else if mytype == "rel_paper_reference" || mytype == "rel_paper_field" {
			signal_map = make_rel_paper_paper(line_list, other_type)
			delete(signal_map, "pid")
		}
		if len(signal_map) == 1 {
			other = append(other, signal_map[other_type])
			continue
		}
		other = append(other, signal_map)
	}
	multi_map[other_type] = other
	return multi_map
}

func make_map(line string, mytype string) map[string]interface{} {
	line_list := strings.Split(line, "\t")

	if mytype == "paper" {
		return make_paper(line_list)
	} else if mytype == "author" {
		return make_author(line_list)
	} else {
		return make(map[string]interface{})
	}
}
func write_to_file(json_list_pointer *[]string, write_obj *bufio.Writer) {
	json_list := *json_list_pointer
	for _, str := range json_list {
		_, err := write_obj.Write([]byte(str + "\n"))
		if err != nil {
			panic(err)
		}
	}
}

func make_json(filename string, mytype string, is_multi bool, main_type string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)
	i := 0

	write_file, err := os.OpenFile("my"+filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	writeObj := bufio.NewWriterSize(write_file, 4096)
	if err != nil {
		panic(err)
	}
	json_list := make([]string, 0, 3000000)
	the_map := make(map[string]interface{})
	line_num := 0
	this_id := "bababababbabab"
	lines := make([]string, 0, 5000)
	for scanner.Scan() {
		line := ""
		line_num += 1
		if is_multi {

			for scanner.Scan() {
				line = scanner.Text()
				line_list := strings.Split(line, "\t")
				if this_id == "bababababbabab" || this_id == line_list[0] {
					lines = append(lines, line)
					this_id = line_list[0]
				} else {
					the_map = make_multi_map(lines, mytype, main_type, this_id, "rel")
					lines = make([]string, 0, 5000)
					this_id = line_list[0]
					lines = append(lines, line)
					break
				}
			}
		} else {
			line = scanner.Text()
			the_map = make_map(line, mytype)
		}

		if the_map != nil {

			json_str, err := json.Marshal(the_map)
			if err != nil {
				panic(err)
			}
			i += 1
			json_list = append(json_list, string(json_str))
			if i%1000000 == 0 {
				if err != nil {
					fmt.Println("cacheFileList.yml file create Failed. err: " + err.Error())
				}
				write_to_file(&json_list, writeObj) // 传递数组为拷贝传递，穿指针省时间
				json_list = make([]string, 0, 3000000)
				fmt.Println(time.Now(), i, line_num)
			}

		}

	}
	write_to_file(&json_list, writeObj)
	json_list = make([]string, 0, 1000000)
	fmt.Println(time.Now(), i, line_num)
}

func make_multi_rel(filename string, mytype string, is_multi bool, main_type string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)
	i := 0

	write_file, err := os.OpenFile("my"+filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	writeObj := bufio.NewWriterSize(write_file, 4096)
	if err != nil {
		panic(err)
	}
	json_list := make([]string, 0, 3000000)
	the_map := make(map[string]interface{})
	line_num := 0
	this_id := "bababababbabab"
	lines := make([]string, 0, 5000)
	for scanner.Scan() {
		line := ""
		line_num += 1

		line = scanner.Text()
		line_list := strings.Split(line, "\t")
		if this_id == "bababababbabab" || this_id == line_list[0] {
			lines = append(lines, line)
			this_id = line_list[0]
		} else {
			the_map = make_multi_map(lines, mytype, main_type, this_id, "rel")
			lines = make([]string, 0, 5000)
			this_id = line_list[0]
			lines = append(lines, line)
			if the_map != nil {

				json_str, err := json.Marshal(the_map)
				if err != nil {
					panic(err)
				}
				i += 1
				json_list = append(json_list, string(json_str))
				if i%1000000 == 0 {
					if err != nil {
						fmt.Println("cacheFileList.yml file create Failed. err: " + err.Error())
					}
					write_to_file(&json_list, writeObj) // 传递数组为拷贝传递，穿指针省时间
					json_list = make([]string, 0, 3000000)
					fmt.Println(time.Now(), i, line_num)
				}

			}

		}

	}
	write_to_file(&json_list, writeObj)
	json_list = make([]string, 0, 1000000)
	fmt.Println(time.Now(), i, line_num)
}

func main() {
	s := time.Now()
	fmt.Println(s)
	//load_map()

	// make_json("Papers.txt", "paper",false,"paper_id")
	fmt.Println(time.Now(), "load paper end")
	// make_json("Authors.txt", "author",false,"author_id")	make_multi_rel("PaperAuthorAffiliations.txt", "rel_paper_author", true, "paper_id")//2021-11-26 22:51:14.0211883 +0800 CST m=+11773.545295501 269412162 731587019
	make_multi_rel("PaperReferences.txt", "rel_paper_reference", true, "paper_id")
	fmt.Println(time.Now())
	make_multi_rel("PaperFieldsOfStudy.txt", "rel_paper_field", true, "paper_id")
	make_multi_rel("PaperCitationContexts.txt", "rel_paper_citation", true, "paper_id")
}
