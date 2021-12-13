package service

import (
	"gitee.com/online-publish/slime-scholar-go/utils"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func InitRedis() {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     utils.REDIS_HOST,
		Password: utils.REDIS_PASSWORD,
		DB:       0,
	})
}

// 使用key存储为列表，用于关注使用
func RedisSaveValue(key string, value string) {
	redisClient.SAdd(key, value)
}

// 使用redis 取出数值,找不到范湖空列表
func RedisGetValueList(key string) []string {
	list, err := redisClient.SMembers(key).Result()
	if err != nil {
		panic(err)
	}
	//fmt.Println(list)
	return list
}

// 查看key中的list是否已经含有了该value
func RedisKeyIsExistValue(key string, value string) bool {
	list := RedisGetValueList(key)
	for _, val := range list {
		if val == value {
			return true
		}
	}
	return false
}

// redis基础操作，设置一个key，以list存储，设置数值与取值
//func main() {
//	fmt.Println("golang连接redis")
//
//	redisClient := redis.NewClient(&redis.Options{
//		Addr:     utils.REDIS_HOST,
//		Password: utils.REDIS_PASSWORD,
//		DB:       0,
//	})
//
//	pong, err := redisClient.Ping().Result()
//	fmt.Println(pong, err)
//
//	setKey := "golang_test_set"
//	redisClient.SAdd(setKey, 1)
//	redisClient.SAdd(setKey, 2)
//	setList, _ := redisClient.SMembers(setKey).Result()
//	fmt.Println("GetSet", setList)
//
//}
