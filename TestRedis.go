package main

import (
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/service"
)

func main() {
	service.InitRedis()
	service.FollowUser(10086, 100887)
	fmt.Println(service.GetUserFollowingList(10086))
	fmt.Println(service.GetUserFollowedList(100887))
	service.CanCelFollowUser(10086, 100887)
	fmt.Println(service.GetUserFollowingList(10086))
	fmt.Println(service.GetUserFollowedList(100887))

}
