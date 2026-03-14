package domain

import (
	"time"
)

type CSAIServiceInterface interface {
	Ask(question string) (string, error)
	QABatchCreate(items []QAItem) (int, error)
}

type QAItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type QAEmbedding struct {
	ID        int64     `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Embedding []float32 `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CSAIRepositoryInterface interface {
	BatchInsert(items []QAEmbedding) error
	SearchSimilar(embedding []float32, limit int) ([]QAEmbedding, error)
}
