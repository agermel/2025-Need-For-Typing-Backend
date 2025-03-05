package dao

import (
	"context"
	"errors"
	"time"
	"type/database"

	"github.com/redis/go-redis/v9"
)

func SaveCodeWithEmail(ctx context.Context, email, code string, ttl time.Duration, rdb *database.RedisClient) error {
	// 设置键值对
	err := rdb.Client.Set(ctx, email, code, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetCodeByEmail(ctx context.Context, email string, code string, rdb *database.RedisClient) error {
	storedCode, err := rdb.Client.Get(ctx, email).Result()
	if err == redis.Nil {
		return errors.New("code过期或不存在")
	} else if err != nil {
		return err
	}

	if storedCode != code {
		return errors.New("验证码错误")
	}

	return nil
}

func DeleteCode(ctx context.Context, email string, rdb *database.RedisClient) error {
	err := rdb.Client.Del(ctx, email).Err()
	if err != nil {
		return err
	}
	return nil
}
