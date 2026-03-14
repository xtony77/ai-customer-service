package csai_repository

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/models"

	"github.com/pgvector/pgvector-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CSAIRepository struct {
	Psql *gorm.DB
}

func NewCSAIRepository(psql *gorm.DB) domain.CSAIRepositoryInterface {
	return &CSAIRepository{
		Psql: psql,
	}
}

func (r *CSAIRepository) BatchInsert(items []domain.QAEmbedding) error {
	if len(items) == 0 {
		return nil
	}

	records := make([]models.QAEmbeddingModel, 0, len(items))
	for _, item := range items {
		records = append(records, models.QAEmbeddingModel{
			Question:  item.Question,
			Answer:    item.Answer,
			Embedding: pgvector.NewVector(item.Embedding),
		})
	}

	if err := r.Psql.CreateInBatches(records, 100).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *CSAIRepository) SearchSimilar(embedding []float32, limit int) ([]domain.QAEmbedding, error) {
	if limit <= 0 {
		limit = 3
	}

	var results []domain.QAEmbedding
	err := r.Psql.Raw(
		`SELECT id, question, answer, created_at, updated_at
			FROM qa_embeddings
			ORDER BY embedding <=> ?
			LIMIT ?`,
		pgvector.NewVector(embedding),
		limit,
	).
		Scan(&results).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return results, nil
}
