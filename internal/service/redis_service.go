package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"userapi/internal/model"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisService {
	return &RedisService{client: client}
}

func (r *RedisService) SetToBlacklist(ctx context.Context, jti string, ttl time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", jti)

	return r.client.Set(ctx, key, "true", ttl).Err()
}

func (r *RedisService) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", jti)
	result, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return result == "true", nil
}

func (r *RedisService) SetCachedUsers(ctx context.Context, users []model.User) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, "users", data, 30*time.Minute).Err()
}

func (r *RedisService) GetCachedUsers(ctx context.Context) ([]model.User, error) {
	cache, err := r.client.Get(ctx, "users").Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var users []model.User
	if err := json.Unmarshal([]byte(cache), &users); err != nil {
		return nil, err
	}

	return users, nil
}
