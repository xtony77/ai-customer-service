package handler

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsk(t *testing.T) {
	t.Cleanup(func() {
		newCSAIService = defaultNewCSAIService
	})

	t.Run("valid request returns 200", func(t *testing.T) {
		service := &mocks.CSAIService{}
		newCSAIService = func() domain.CSAIServiceInterface {
			return service
		}
		service.On("Ask", "營業時間是？").Return("週一到週五 09:00-18:00", nil).Once()

		response := performJSONRequest(t, Ask, "/api/v1/chat", `{"question":"營業時間是？"}`)

		assert.Equal(t, 200, response.Code)
		assert.JSONEq(t, `{"answer":"週一到週五 09:00-18:00"}`, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("missing question returns 400", func(t *testing.T) {
		response := performJSONRequest(t, Ask, "/api/v1/chat", `{}`)

		assert.Equal(t, 400, response.Code)
		assert.JSONEq(t, `{"code":400,"message":"bad request"}`, response.Body.String())
	})

	t.Run("service error returns 500", func(t *testing.T) {
		service := &mocks.CSAIService{}
		newCSAIService = func() domain.CSAIServiceInterface {
			return service
		}
		service.On("Ask", "系統壞了嗎？").Return("", errors.New("service failed")).Once()

		response := performJSONRequest(t, Ask, "/api/v1/chat", `{"question":"系統壞了嗎？"}`)

		assert.Equal(t, 500, response.Code)
		assert.JSONEq(t, `{"code":500,"message":"Server Error"}`, response.Body.String())
		service.AssertExpectations(t)
	})
}
