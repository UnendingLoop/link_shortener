package cache

import (
	"context"
	"log"
	"time"

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
	// data, err := json.Marshal(pointer)
	// if err != nil {
	// 	return err
	// }

	// return r.client.SetWithExpiration(ctx, uid, data, r.ttl)
	return nil
}

func (r *RedisCache) GetByShortkey(ctx context.Context, uid string) (string, error) {
	// data, err := r.client.Get(ctx, uid)
	// if err != nil {
	// 	return nil, err
	// }
	// if len(data) == 0 {
	// 	return nil, repository.ErrNotFound
	// }

	// var candidate repository.Notification
	// if err := json.Unmarshal([]byte(data), &candidate); err != nil {
	// 	return nil, err
	// }
	// return &candidate, nil
	return "nil", nil
}

func (r *RedisCache) DeleteByShortkey(ctx context.Context, uid string) error {
	return r.client.Del(ctx, uid)
}
