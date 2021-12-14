package service

import (
	"encoding/json"
	"sort"
	"strconv"
	"time"
)

func GetMapAllKey(m map[string]interface{}) []string {
	list := make([]string, 0, len(m))
	for key := range m {
		list = append(list, key)
	}
	return list
}

func GetMapAllContent(m map[string]interface{}) []interface{} {
	list := make([]interface{}, 0, len(m))
	for key := range m {
		list = append(list, m[key])
	}
	return list
}

func StructToMap(s interface{}) (ret_map map[string]interface{}) {
	str, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(str), &ret_map)
	return ret_map
}

// 对mao 按照key 排序

func RankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	//从小到大排序
	//sort.Sort(pl)
	//从大到小排序
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func GetTopNKey(m map[string]int, n int) (ret []string) {
	pl := RankByWordCount(m)
	for i := 0; i < len(pl) && i < n; i++ {
		ret = append(ret, pl[i].Key)
	}
	return ret
}

//对map按照value排序后返回
func GetAllSortedKey(m map[string]int) (ret []string) {
	pl := RankByWordCount(m)
	for i := 0; i < len(pl); i++ {
		ret = append(ret, pl[i].Key)
	}
	return ret
}

//忽略所恶的将字符串转化为数字
func PureAtoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

const (
	YearMonthDay     = "2006-01-02"
	HourMinuteSecond = "15:04:05"
	DefaultLayout    = YearMonthDay + " " + HourMinuteSecond
)

// 默认格式日期字符串转time
func TimeStrToTimeDefault(str string) time.Time {
	parseTime, _ := time.ParseInLocation(DefaultLayout, str, time.Local)
	return parseTime
}
