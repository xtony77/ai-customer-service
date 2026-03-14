package cache

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/configs"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Conn *redis.Client
}

var rdb *redis.Client
var redisOnce sync.Once

func NewRedis(config *configs.Config) domain.RedisInterface {
	redisOnce.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
			Password: config.Redis.Password,
			DB:       config.Redis.Database,
		})
	})

	return &RedisService{
		Conn: rdb,
	}
}

func (r *RedisService) Get(key string) (string, error) {
	value, err := r.Conn.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", errors.WithStack(err)
	}

	return value, nil
}

func (r *RedisService) Set(key string, value string, expiration time.Duration) error {
	if err := r.Conn.Set(context.Background(), key, value, expiration).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
