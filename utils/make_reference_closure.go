package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// 根据cspaperid 根据reference 获取闭包。经计算10次迭代即可收敛
func getKeys1(m map[int64]int) int {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	j := 0
	for _ = range m {
		j++
	}
	return j
}

func MakeReferenceClosure() {
	s := time.Now()
	fmt.Println(s)

	file_cspaper, err := os.Open("cspaper_ids.txt")
	if err != nil {
		panic(err)
	}
	cs_paper_idmap := make(map[int64]int)
	scanner := bufio.NewScanner(file_cspaper)
	//last_paper_map := make(map[int64]int)
	this_paper_map := make(map[int64]int)
	cspaper_num := 0
	for scanner.Scan() {
		cspaper_num++
		line, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			panic(err)
		}
		cs_paper_idmap[line] = 1
		this_paper_map[line] = 1
	}
	all_paper_but_notcs := make(map[int64]int)

	len_closure := 0
	for len_closure == 0 || len_closure < getKeys1(all_paper_but_notcs) {
		fmt.Println("before", len_closure)
		len_closure = getKeys1(all_paper_but_notcs)
		fmt.Println("after", len_closure)
		last_paper_map := this_paper_map
		this_paper_map = make(map[int64]int)
		file2, err := os.Open("PaperReferences.txt")
		if err != nil {
			fmt.Println(err)
		}
		scanner = bufio.NewScanner(file2)
		line_num := 0
		for scanner.Scan() {
			line_num++
			if line_num%10000000 == 0 {
				fmt.Println(time.Now(), line_num, len_closure)
			}
			line := scanner.Text()
			lines := strings.Split(line, "\t")
			line0, err := strconv.ParseInt(lines[0], 10, 64)
			line1, err := strconv.ParseInt(lines[1], 10, 64)
			if err != nil {
				panic(err)
			}
			if _, ok := last_paper_map[line0]; ok {
				if _, ok2 := all_paper_but_notcs[line1]; !ok2 {
					if _, ok3 := cs_paper_idmap[line1]; !ok3 {
						this_paper_map[line1] = 1
						all_paper_but_notcs[line1] = 1
					}

				}
			}
		}
		fmt.Println(time.Now(), len_closure)
	}

	f, err1 := os.OpenFile("cspaperclosure_ids.txt", os.O_APPEND, 0666) //打开文件
	if err1 != nil {
		panic(err1)
	}
	for str := range all_paper_but_notcs {
		_, err12 := io.WriteString(f, strconv.FormatInt(str, 10)+"\n") //写入文件(字符串)

		if err12 != nil {
			panic(err12)
		}
	}
	fmt.Println(time.Now())

}
