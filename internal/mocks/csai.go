package mocks

import (
	"ai-customer-service/domain"
	"time"

	"github.com/stretchr/testify/mock"
)

type OpenAI struct {
	mock.Mock
}

func (m *OpenAI) GenerateEmbedding(input string) ([]float32, error) {
	args := m.Called(input)
	if embedding := args.Get(0); embedding != nil {
		return embedding.([]float32), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *OpenAI) ChatCompletion(question string, references []domain.QAEmbedding) (string, error) {
	args := m.Called(question, references)
	return args.String(0), args.Error(1)
}

type CSAIRepository struct {
	mock.Mock
}

func (m *CSAIRepository) BatchInsert(items []domain.QAEmbedding) error {
	args := m.Called(items)
	return args.Error(0)
}

func (m *CSAIRepository) SearchSimilar(embedding []float32, limit int) ([]domain.QAEmbedding, error) {
	args := m.Called(embedding, limit)
	if results := args.Get(0); results != nil {
		return results.([]domain.QAEmbedding), args.Error(1)
	}

	return nil, args.Error(1)
}

type Redis struct {
	mock.Mock
}

func (m *Redis) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *Redis) Set(key string, value string, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

type CSAIService struct {
	mock.Mock
}

func (m *CSAIService) Ask(question string) (string, error) {
	args := m.Called(question)
	return args.String(0), args.Error(1)
}

func (m *CSAIService) QABatchCreate(items []domain.QAItem) (int, error) {
	args := m.Called(items)
	return args.Int(0), args.Error(1)
}
