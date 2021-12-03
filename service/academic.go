package service

func MostSimplifyPaper(m map[string]interface{}) (ret map[string]interface{}) {
	ret["paper_id"] = m["paper_id"]
	ret["paper_title"] = m["paper_title"]
	return ret
}
func SimplifyPapers(inter []interface{}) []interface{} {
	ret_list := make([]interface{}, len(inter))
	for _, v := range inter {
		ret_list = append(ret_list, MostSimplifyPaper(v.(map[string]interface{})))
	}
	return ret_list
}
