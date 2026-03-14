package handler

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQABatchCreate(t *testing.T) {
	t.Cleanup(func() {
		newCSAIService = defaultNewCSAIService
	})

	t.Run("valid batch returns 200", func(t *testing.T) {
		service := &mocks.CSAIService{}
		items := []domain.QAItem{
			{Question: "Q1", Answer: "A1"},
			{Question: "Q2", Answer: "A2"},
		}
		newCSAIService = func() domain.CSAIServiceInterface {
			return service
		}
		service.On("QABatchCreate", items).Return(2, nil).Once()

		response := performJSONRequest(t, QABatchCreate, "/api/v1/qa", `{"items":[{"question":"Q1","answer":"A1"},{"question":"Q2","answer":"A2"}]}`)

		assert.Equal(t, 200, response.Code)
		assert.JSONEq(t, `{"count":2}`, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("empty items returns 400", func(t *testing.T) {
		response := performJSONRequest(t, QABatchCreate, "/api/v1/qa", `{"items":[]}`)

		assert.Equal(t, 400, response.Code)
		assert.JSONEq(t, `{"code":400,"message":"bad request"}`, response.Body.String())
	})

	t.Run("missing required fields returns 400", func(t *testing.T) {
		response := performJSONRequest(t, QABatchCreate, "/api/v1/qa", `{"items":[{"question":"Q1","answer":""}]}`)

		assert.Equal(t, 400, response.Code)
		assert.JSONEq(t, `{"code":400,"message":"bad request"}`, response.Body.String())
	})
}
