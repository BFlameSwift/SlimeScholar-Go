package main

import (
	"fmt"

	"github.com/BFlameSwift/SlimeScholar-Go/service"
	"golang.org/x/net/context"
)

func main() {
	service.Init()
	fmt.Println("123")
	// var map_param map[string]string = make(map[string]string)
	// e1, _ := json.Marshal(model.ValueString{Value: "132"})

	// map_param["index"], map_param["id"], map_param["bodyJson"] = "conference", "1158167855", string(e1)
	ret, _ := service.Client.Get().Index("conference").Id("1158167855").Do(context.Background())
	fmt.Printf(string(ret.Source))
}
