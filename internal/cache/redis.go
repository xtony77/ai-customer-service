package cache

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/configs"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Config *configs.Config
	Conn   *redis.Client
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
		Config: config,
		Conn:   rdb,
	}
}
