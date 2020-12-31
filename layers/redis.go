package layers

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	Cache *redis.Client
	Write *redis.Client
	Ttl   time.Duration
}

func (b *Redis) Get(key string) ([]byte, error) {
	data := b.Cache.Get(context.Background(), key)
	if data.Err() != nil {
		return nil, data.Err()
	} else {
		return []byte(data.Val()), nil
	}
}

func (b *Redis) Set(key string, value []byte) error {
	return b.Write.Set(context.Background(), key, string(value), b.Ttl).Err()
}
