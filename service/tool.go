package service

import "encoding/json"

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
