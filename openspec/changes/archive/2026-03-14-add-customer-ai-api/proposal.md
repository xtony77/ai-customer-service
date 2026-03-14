## Why

建立 AI 客服 API，讓用戶可以透過 API 提問，系統自動從知識庫中搜尋相關 Q&A 並透過 GPT-4o mini 生成回答。透過 semantic cache 和向量搜尋減少重複的 AI API 呼叫成本。

## What Changes

- 新增 `POST /api/v1/chat` 端點：接收用戶問題，經過 prompt injection 防護 → Redis 精確快取 → pgvector 向量搜尋 top-3 Q&A → GPT-4o mini 生成回答 → 存入 Redis 快取
- 新增 `POST /api/v1/qa` 端點：批次匯入 Q&A 資料，server 端自動呼叫 OpenAI Embedding API 產生 embedding 後存入 pgvector
- 新增 `qa_embeddings` 資料表（pgvector + HNSW index），使用 golang-migrate/migrate 管理 migration
- Docker Compose 的 PostgreSQL image 改為 `pgvector/pgvector:pg16` 以支援 vector extension
- `.env` 新增 `OPENAI_API_KEY`、`OPENAI_MODEL`（預設 gpt-4o-mini）、`OPENAI_EMBEDDING_MODEL`（預設 text-embedding-3-small）
- 整合 OpenAI API client（chat completion + embedding）

## Capabilities

### New Capabilities
- `chat-api`: 用戶提問端點，包含 prompt injection 防護、Redis 精確快取、pgvector RAG top-3、GPT-4o mini 生成回答
- `qa-import`: 批次匯入 Q&A 資料端點，自動產生 embedding 並存入 pgvector
- `prompt-guard`: Prompt injection 基本防護（關鍵字/正則過濾）

### Modified Capabilities
（無既有 capability 需修改）

## Impact

- **新增套件**: `github.com/golang-migrate/migrate/v4`、`github.com/pgvector/pgvector-go`、`github.com/sashabaranov/go-openai`
- **Config 變更**: `internal/configs/config.go` 新增 OpenAI 相關設定
- **Docker**: `docker-compose.yml` PostgreSQL image 變更
- **資料庫**: 新增 `qa_embeddings` 表 + pgvector extension + HNSW index
- **API**: 新增兩個公開端點（無 JWT 認證）
