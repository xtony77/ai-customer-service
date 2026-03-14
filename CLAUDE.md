# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AI Customer Service API — a RAG (Retrieval Augmented Generation) chatbot built in Go. Takes user questions, generates embeddings via OpenAI, searches a pgvector database for similar Q&A pairs, and returns GPT-generated answers. Includes prompt injection protection and Redis caching.

## Common Commands

```bash
# Run the server
go run cmd/gin/main.go

# Run all tests
go test ./...

# Run a single package's tests
go test ./internal/services/csai_service/

# Regenerate Wire dependency injection (after changing providers/constructors)
make wire

# Run database migrations
make migrate-up

# Format code
gofmt -w .
```

## Infrastructure

Docker Compose provides PostgreSQL 16 (with pgvector) and Redis 7:
```bash
docker-compose up -d
```

Copy `.env.example` to `.env` and fill in credentials before running.

## Architecture

Layered architecture with constructor-based dependency injection (Google Wire):

```
cmd/gin/main.go          → entrypoint, wiring only
domain/                   → shared interfaces, DTOs, error definitions
internal/
  configs/                → env-based config (singleton, sync.Once)
  gin/handler/            → thin Gin handlers (bind → service → respond)
  gin/routes/             → route registration
  services/csai_service/  → business logic orchestration (ask, batch import)
  services/openai_service/→ OpenAI API (embeddings, chat completion)
  repositories/           → data access (pgvector similarity search)
  cache/                  → Redis adapter
  database/               → PostgreSQL/GORM adapter
  wire/                   → Wire DI configuration
  mocks/                  → testify mocks
  models/                 → GORM models
  logger/                 → slog-based structured logging (singleton)
  jwt/                    → JWT auth service
```

**Request flow (Chat):** Handler → CSAIService (prompt guard → Redis cache check → OpenAI embedding → pgvector top-3 similarity search → GPT chat completion → cache result) → Handler response

**Key interfaces in `domain/`:** CSAIServiceInterface, CSAIRepositoryInterface, OpenAIInterface, RedisInterface. All layers depend on these interfaces, not concrete types.

## Code Conventions

See `AGENTS.md` for full details. Key points:

- **Response helpers:** Always use `Success()` / `Failed()` from `handler/helpers.go`. Never call `c.JSON()` directly.
- **Config access:** Only through `configs.NewConfig()`. No direct `os.Getenv` elsewhere.
- **Error wrapping:** Use `github.com/pkg/errors` with `errors.WithStack`.
- **Package naming:** underscore-separated directories (`csai_service`, `openai_service`).
- **Interfaces:** defined in `domain/`, suffixed with `Interface`.
- **Singletons:** config, logger, database, cache, OpenAI client all use `sync.Once`.

## Testing

Tests use testify mocks (`internal/mocks/`). Handler tests create HTTP requests against Gin test contexts. OpenAI service tests stub the HTTP client. Run `make wire` before tests if dependency wiring changed.

## Database

Single migration in `migrations/` creates `qa_embeddings` table with a `VECTOR(1536)` column and HNSW index for cosine similarity search (`<=>` operator).
