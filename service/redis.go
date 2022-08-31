package service

import (
	"github.com/BFlameSwift/SlimeScholar-Go/utils"
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

func RedisRemoveValue(key string, value string) {
	redisClient.SRem(key, value)
}

func RedisSaveValueSorted(key string, value string) {
	redisClient.RPush(key, value)
}
func RedisGetValueSorted(key string) []string {
	length, err := redisClient.LLen(key).Result()
	if err != nil {
		panic(err)
	}
	vals, err := redisClient.LRange(key, 0, length).Result()
	if err != nil {
		panic(err)
	}
	return vals
}
