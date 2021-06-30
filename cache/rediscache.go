package cache

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"os"
)

var RedisDb *redis.Client

func Start()  {
	redisDbAddr := "192.168.99.101:6379"
	if os.Getenv("REDIS_DB_ADDR") != "" {
		redisDbAddr = os.Getenv("REDIS_DB_ADDR")
	}
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     redisDbAddr,
		Password: os.Getenv("REDIS_DB_AUTH"),
		DB:       0,  // use default DB
	})

	if ping := RedisDb.Ping(); ping.Err() != nil {
		fmt.Println("Fail to connect to redis", ping.Err().Error())
		panic(ping.Err())
	}
	fmt.Println("Connected to Redis", redisDbAddr)
}



func Exists(key string) (bool, error) {
	if r, err := RedisDb.Exists(key).Result(); err != nil {
		return false, err
	} else {
		return r == 1, nil
	}
}

func Save(key, value string) (error) {
	return RedisDb.Set(key, value, 0).Err()
}

func Get(key string) (string, error)  {
	return RedisDb.Get(key).Result()
}
