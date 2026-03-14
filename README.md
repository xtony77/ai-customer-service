# AI Customer Service

AI 客服 API，基於 RAG 架構，透過向量搜尋知識庫中的相似 Q&A，搭配 GPT-4o mini 生成回答。

## 架構流程

```
用戶問題 → Prompt Injection 防護 → Redis 快取查詢
  ├─ 命中 → 回傳快取回答
  └─ 未命中 → OpenAI Embedding → pgvector Top-3 相似 Q&A → GPT-4o mini 生成回答 → 存入快取 → 回傳
```

## 技術棧

- **Go** + **Gin** — HTTP API
- **PostgreSQL** + **pgvector** (HNSW) — 向量搜尋
- **Redis** — 精確匹配快取 (SHA256)
- **OpenAI API** — Embedding (text-embedding-3-small) + Chat Completion (GPT-4o mini)
- **GORM** — ORM
- **golang-migrate** — 資料庫 Migration
- **Wire** — 依賴注入

## API

| Method | Path | 說明 |
|--------|------|------|
| POST | `/api/v1/chat` | 客服問答 |
| POST | `/api/v1/qa` | 批次匯入 Q&A 知識庫 |

## 快速開始

```bash
# 啟動 PostgreSQL (pgvector) + Redis
docker compose up -d

# 設定環境變數
cp .env.example .env
# 編輯 .env 填入 OPENAI_API_KEY

# 執行 Migration
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# 啟動服務
go run cmd/gin/main.go
```
