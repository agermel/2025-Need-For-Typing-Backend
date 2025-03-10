package database

import (
	"context"
	"log"
	"type/config"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

var Rdb = NewRedis()

func NewRedis() *RedisClient {
	config.LoadConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.AllConfig.Server.RedisURL,
		Password: config.AllConfig.Database.RedisPassword,
		DB:       0, // 默认第0个数据库
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}

	return &RedisClient{Client: rdb}
}
