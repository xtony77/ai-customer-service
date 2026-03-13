package jwt

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/configs"
)

type JWTService struct {
	Config *configs.Config
}

func NewJWT(config *configs.Config) domain.JWTInterface {
	return &JWTService{
		Config: config,
	}
}
