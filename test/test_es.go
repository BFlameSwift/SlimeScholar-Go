package main

import (
	"encoding/json"
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/model"

	"gitee.com/online-publish/slime-scholar-go/service"
)


func main() {
	service.Init()
	fmt.Println("123")
	var map_param map[string]string = make(map[string]string)
	e1, _ := json.Marshal(model.ValueString{Value: "132"})

	map_param["index"], map_param["type"], map_param["id"], map_param["bodyJson"] = "megacorp", "employee", "5f", string(e1)
	ret := service.Create(map_param)
	fmt.Printf(ret)
}
