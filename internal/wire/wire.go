//go:build wireinject

package wire

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/cache"
	"ai-customer-service/internal/configs"
	"ai-customer-service/internal/database"
	"ai-customer-service/internal/jwt"

	"github.com/google/wire"
)

func NewConfig() *configs.Config {
	return configs.NewConfig()
}

func NewJWT() domain.JWTInterface {
	panic(wire.Build(jwt.NewJWT, NewConfig))
}

func NewPostgreSQL() (domain.PostgreSQLInterface, error) {
	panic(wire.Build(database.NewPostgreSQL, NewConfig))
}

func NewRedis() (domain.RedisInterface, error) {
	panic(wire.Build(cache.NewRedis, NewConfig))
}
