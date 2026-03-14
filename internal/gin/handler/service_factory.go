package handler

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/wire"
)

var defaultNewCSAIService = func() domain.CSAIServiceInterface {
	return wire.NewCSAIService()
}

var newCSAIService = defaultNewCSAIService
