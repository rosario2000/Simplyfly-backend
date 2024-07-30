package redis

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)

var (
	client *redis.Client
)

func Instance() *redis.Client {
	opt, _ := redis.ParseURL("rediss://default:AeIHAAIjcDE4ODczNmU0OTdlMWI0ZTA3OGEwOTBlMmQ1OGJiYjFkMHAxMA@clever-meerkat-57863.upstash.io:6379")
	client = redis.NewClient(opt)
	return client
}

func Get(ctx context.Context, key string) (string, error) {
	return Instance().Get(key).Result()
}
func Set(ctx context.Context, key string, val string, expiration time.Duration) error {
	return Instance().Set(key, val, expiration).Err()
}

func Delete(ctx context.Context, key string) error {
	return Instance().Del(key).Err()
}
