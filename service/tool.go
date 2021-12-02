package service

func GetMapAllContent(m map[string]interface{}) []interface{} {
	list := make([]interface{}, 0, len(m))
	for key := range m {
		list = append(list, m[key])
	}
	return list
}
