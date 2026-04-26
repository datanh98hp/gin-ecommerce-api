package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dat19/gin-ecommerce-api/internal/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect(cfg *config.Config) error {
	Client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     50,
		MinIdleConns: 10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Client.Set(ctx, key, data, ttl).Err()
}

func Get(ctx context.Context, key string, dest interface{}) error {
	data, err := Client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func Delete(ctx context.Context, keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

func DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = Client.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := Client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		n += len(keys)
		if cursor == 0 {
			break
		}
	}
	return nil
}

func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}
