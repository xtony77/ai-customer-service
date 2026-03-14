package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
)

type QAEmbeddingModel struct {
	ID        int64           `gorm:"column:id;primaryKey;autoIncrement"`
	Question  string          `gorm:"column:question;type:text;not null"`
	Answer    string          `gorm:"column:answer;type:text;not null"`
	Embedding pgvector.Vector `gorm:"column:embedding;type:vector(1536);not null"`
	CreatedAt time.Time       `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime"`
	UpdatedAt time.Time       `gorm:"column:updated_at;type:timestamptz;not null;autoUpdateTime"`
}

func (QAEmbeddingModel) TableName() string {
	return "qa_embeddings"
}
