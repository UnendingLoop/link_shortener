package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"shortener/internal/repository"

	"github.com/wb-go/wbf/redis"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(client *redis.Client, ttl time.Duration) ShortCache {
	var rc RedisCache
	if ttl <= 0 {
		rc.ttl = time.Minute
		log.Println("Got incorrect TTL for storing notifications in Redis. Continue with default value 1m(60 seconds)...")
	} else {
		rc.ttl = ttl
	}
	rc.client = client
	return &rc
}

func (r *RedisCache) SetByShortkey(ctx context.Context, key string, link string) error {
	data, err := json.Marshal(link)
	if err != nil {
		return err
	}

	return r.client.SetWithExpiration(ctx, key, data, r.ttl)
}

func (r *RedisCache) GetByShortkey(ctx context.Context, key string) (string, error) {
	data, err := r.client.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", repository.ErrNotFound
	}

	result := ""
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return "", err
	}
	return result, nil
}

func (r *RedisCache) DeleteByShortkey(ctx context.Context, key string) error {
	return r.client.Del(ctx, key)
}
