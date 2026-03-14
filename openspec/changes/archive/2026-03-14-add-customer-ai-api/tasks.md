## 1. Infrastructure & Config

- [x] 1.1 修改 `docker-compose.yml`：PostgreSQL image 改為 `pgvector/pgvector:pg16`
- [x] 1.2 新增 `.env` / `.env.example` 的 OpenAI 設定（`OPENAI_API_KEY`、`OPENAI_MODEL`、`OPENAI_EMBEDDING_MODEL`）
- [x] 1.3 更新 `internal/configs/config.go`：新增 `OpenAI` struct（APIKey、Model、EmbeddingModel）
- [x] 1.4 安裝新套件：`golang-migrate/migrate/v4`、`pgvector/pgvector-go`、`sashabaranov/go-openai`

## 2. Database Migration

- [x] 2.1 建立 migration 檔案 `migrations/000001_create_qa_embeddings.up.sql`：CREATE EXTENSION vector + CREATE TABLE qa_embeddings + HNSW index
- [x] 2.2 建立 migration 檔案 `migrations/000001_create_qa_embeddings.down.sql`：DROP TABLE + DROP EXTENSION

## 3. Domain & Models

- [x] 3.1 定義 `domain/qa.go`：QAEmbedding model struct、QARepository interface、QAService interface
- [x] 3.2 定義 `domain/openai.go`：OpenAIClient interface（GenerateEmbedding、ChatCompletion）
- [x] 3.3 定義 `domain/chat.go`：ChatRequest/ChatResponse struct、ChatService interface

## 4. OpenAI Client

- [x] 4.1 實作 `internal/openai/client.go`：OpenAI client wrapper，實作 GenerateEmbedding 和 ChatCompletion

## 5. Prompt Guard

- [x] 5.1 實作 `internal/guard/prompt_guard.go`：關鍵字/正則 prompt injection 過濾器

## 6. Repository Layer

- [x] 6.1 實作 `internal/repository/qa_repository.go`：BatchInsert（批次寫入 Q&A + embedding）、SearchSimilar（pgvector cosine top-3）
- [x] 6.2 實作 `internal/repository/cache_repository.go`：Redis GET/SET（sha256 key，永久快取）

## 7. Service Layer

- [x] 7.1 實作 `internal/service/qa_service.go`：批次匯入邏輯（呼叫 embedding API → 寫入 pgvector）
- [x] 7.2 實作 `internal/service/chat_service.go`：完整 chat 流程（prompt guard → cache → embedding → RAG → GPT → cache 回寫）

## 8. Handler Layer (Gin)

- [x] 8.1 實作 `internal/handler/qa_handler.go`：`POST /api/v1/qa` handler，request binding + validation
- [x] 8.2 實作 `internal/handler/chat_handler.go`：`POST /api/v1/chat` handler，request binding + validation
- [x] 8.3 更新 `internal/gin/routes/route.go`：註冊新的 route

## 9. Wire DI

- [x] 9.1 更新 `internal/wire/wire.go`：新增 OpenAI client、repository、service、handler 的 provider
- [x] 9.2 執行 `wire gen` 重新產生 `wire_gen.go`

## 10. Unit Tests

- [x] 10.1 安裝測試套件：`testify`（assert + mock）
- [x] 10.2 建立 mock 檔案 `internal/mocks/`：mock OpenAIInterface、CSAIRepositoryInterface、RedisInterface
- [x] 10.3 測試 Prompt Guard：驗證 injection 關鍵字被攔截、正常問題通過（`internal/services/csai_service/main_test.go` 或獨立檔案）
- [x] 10.4 測試 Chat Service（Ask）：cache hit 直接回傳、cache miss 走完整流程（embedding → RAG → GPT → cache 寫入）、OpenAI 錯誤處理、無 Q&A 資料場景
- [x] 10.5 測試 QA Service（QABatchCreate）：批次匯入成功、embedding API 失敗、空 items 處理
- [x] 10.6 測試 Chat Handler：valid request 200、missing question 400、service error 500
- [x] 10.7 測試 QA Handler：valid batch 200、empty items 400、missing fields 400
- [x] 10.8 測試 OpenAI Client：mock HTTP response 驗證 embedding 和 chat completion 解析正確
