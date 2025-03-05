package database

import (
	"type/config"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

var Rdb = NewRedis()

func NewRedis() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.AllConfig.Server.RedisURL,
		Password: config.AllConfig.Database.RedisPassword,
		DB:       0, // 默认第0个数据库
	})
	return &RedisClient{Client: rdb}
}
