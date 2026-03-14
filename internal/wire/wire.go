//go:build wireinject

package wire

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/cache"
	"ai-customer-service/internal/configs"
	"ai-customer-service/internal/database"
	"ai-customer-service/internal/jwt"
	"ai-customer-service/internal/repositories/csai_repository"
	"ai-customer-service/internal/services/csai_service"
	"ai-customer-service/internal/services/openai_service"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewConfig() *configs.Config {
	return configs.NewConfig()
}

func NewJWT() domain.JWTInterface {
	panic(wire.Build(jwt.NewJWT, NewConfig))
}

func NewPostgreSQL() *gorm.DB {
	panic(wire.Build(database.NewPostgreSQL, NewConfig))
}

func NewRedis() domain.RedisInterface {
	panic(wire.Build(cache.NewRedis, NewConfig))
}

func NewOpenAI() domain.OpenAIInterface {
	panic(wire.Build(openai_service.NewOpenAI, NewConfig))
}

func NewCSAIRepository() domain.CSAIRepositoryInterface {
	panic(wire.Build(csai_repository.NewCSAIRepository, NewPostgreSQL))
}

func NewCSAIService() domain.CSAIServiceInterface {
	panic(wire.Build(csai_service.NewCSAIService, NewRedis, NewOpenAI, NewCSAIRepository))
}
