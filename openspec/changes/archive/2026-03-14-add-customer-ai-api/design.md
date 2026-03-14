## Context

現有專案使用 Gin + GORM + Redis + Wire DI 架構，已有基本的 PostgreSQL 和 Redis 連線設定，但尚無業務邏輯。需要在此基礎上建立 AI 客服 API，包含向量搜尋、語意快取和 OpenAI 整合。

現有結構：
- `cmd/gin/main.go` → 啟動 Gin server
- `internal/configs/config.go` → 環境變數設定（Gin、JWT、PostgreSQL、Redis）
- `internal/database/postgresql.go` → GORM PostgreSQL 連線
- `internal/cache/redis.go` → go-redis 連線
- `internal/wire/` → Wire DI
- `domain/` → interface 定義
- `docker-compose.yml` → PostgreSQL 16 + Redis 7

## Goals / Non-Goals

**Goals:**
- 提供 `POST /api/v1/chat` 端點，用戶提問後經 prompt injection 防護 → Redis 精確快取 → pgvector top-3 RAG → GPT-4o mini 回答
- 提供 `POST /api/v1/qa` 端點，批次匯入 Q&A 並自動產生 embedding
- 使用 golang-migrate/migrate 管理資料庫 migration
- 基本的 prompt injection 防護（關鍵字/正則過濾）
- AI 回答限制 150 字內，僅根據 top-3 參考資料回答

**Non-Goals:**
- SSE streaming 回應
- JWT 認證（展示用，暫不加）
- Admin Q&A 管理 CRUD（只做批次匯入）
- 進階 prompt injection 防護（LLM 偵測、Moderation API）
- Semantic cache（語意相似度比對），僅做精確匹配快取

## Decisions

### 1. 快取策略：Redis 精確匹配

**選擇**: 用 `sha256(question)` 作為 Redis key，value 為 JSON 回答，不設 TTL。

**替代方案**: Redis Stack 向量搜尋做語意快取 → 過於複雜，且後面已有 pgvector 做向量搜尋兜底。

**理由**: 簡單有效，精確匹配的 cache hit 直接省掉 embedding + pgvector + GPT 三次呼叫。Cache miss 由 pgvector RAG 處理語意相似的問題。

### 2. 向量搜尋：pgvector + HNSW

**選擇**: pgvector extension + HNSW index，embedding 維度 1536（text-embedding-3-small）。

**替代方案**: IVFFlat index → 需要先 train，資料量小時 HNSW 表現更好。

**理由**: HNSW 不需 training，查詢速度快，適合中小規模 Q&A 知識庫。

### 3. OpenAI Client：go-openai

**選擇**: `github.com/sashabaranov/go-openai` 套件。

**理由**: Go 社群最成熟的 OpenAI client，支援 chat completion 和 embedding API。

### 4. Migration：golang-migrate/migrate

**選擇**: 使用 golang-migrate/migrate CLI 管理 schema migration，GORM 僅用於查詢。

**理由**: 用戶指定使用此套件。Migration 檔案放在 `migrations/` 目錄。

### 5. Embedding 產生：Server 端

**選擇**: Q&A 匯入時由 server 呼叫 OpenAI Embedding API 產生 embedding。

**理由**: 用戶端不需要處理 embedding 邏輯，匯入 API 只需提供 question + answer。

### 6. 專案分層架構

遵循現有 domain-driven 分層：
- `domain/` → interface + model 定義
- `internal/handler/` → Gin handler（HTTP 層）
- `internal/service/` → 業務邏輯層
- `internal/repository/` → 資料存取層（PostgreSQL、Redis）

### 7. Docker PostgreSQL Image

**選擇**: `pgvector/pgvector:pg16` 取代 `postgres:16-alpine`。

**理由**: 內建 pgvector extension，無需手動安裝。

## Risks / Trade-offs

- **[OpenAI API 延遲]** → Cache miss 時需要 2 次 OpenAI 呼叫（embedding + chat），延遲約 1-3 秒。透過 Redis 精確快取減少重複呼叫。
- **[Embedding 批次匯入效能]** → 大量 Q&A 匯入時，逐筆呼叫 OpenAI Embedding API 較慢。目前為展示用途，暫不做批次 embedding 優化。
- **[pgvector HNSW 記憶體]** → HNSW index 需要較多記憶體，但中小規模 Q&A（< 10 萬筆）不成問題。
- **[Prompt injection 防護有限]** → 基本關鍵字過濾可被繞過，但搭配 system prompt 限制（只根據參考資料回答、150 字限制），風險可控。
