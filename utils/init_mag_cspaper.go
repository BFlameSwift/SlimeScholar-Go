package utils

// "io"

// func main() {
// 	s := time.Now()
// 	fmt.Println(s)
// 	file, err := os.Open("cs_fields_closure.txt")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	scanner := bufio.NewScanner(file)
// 	var field_ids map[string]int = make(map[string]int)
// 	field_num := 0
// 	for scanner.Scan() {
// 		field_num++
// 		line := scanner.Text()
// 		field_ids[line] = 1
// 		// _,ok := field_ids[line]
// 		// fmt.Println(line)
// 	}
// 	fmt.Println("field num ", field_num)
// 	file2, err := os.Open("PaperFieldsOfStudy.txt")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	paper_map := make(map[int64]int)
// 	scanner = bufio.NewScanner(file2)
// 	// cs_paper_ids := make([]string,0,30000000)
// 	i := 0
// 	for scanner.Scan() {
// 		if i++; i%10000000 == 0 {
// 			fmt.Println(i, (time.Now()), len(paper_map))
// 		}
// 		this_list := strings.Split(scanner.Text(), "\t")
// 		if _, ok := field_ids[this_list[1]]; ok {
// 			paper_map[strconv.ParseInt(this_list[0], 10, 64)] = 1
// 			// len_data := len(cs_paper_ids)
// 			// if len_data == 0 || cs_paper_ids[len_data-1] != this_list[0]{
// 			// 	cs_paper_ids = append(cs_paper_ids,this_list[0])
// 			// }

// 		}
// 		//append(cs_paper_ids,scanner.Text())
// 	}

// 	// f, err1 := os.OpenFile("cspaper_ids.txt", os.O_APPEND, 0666) //打开文件
// 	// if err1!=nil{ panic(err1)}
// 	// for str := range(paper_map){
// 	// 	_, err12 := io.WriteString(f, str+"\n") //写入文件(字符串)
// 	// 	delete(paper_map,str)
// 	// 	if err12 != nil {panic(err12)}
// 	// }
// 	fmt.Println(time.Now())

// }
