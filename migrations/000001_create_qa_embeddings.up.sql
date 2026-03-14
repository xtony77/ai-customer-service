CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS qa_embeddings (
    id BIGSERIAL PRIMARY KEY,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    embedding VECTOR(1536) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_qa_embeddings_embedding_hnsw
    ON qa_embeddings
    USING hnsw (embedding vector_cosine_ops);
