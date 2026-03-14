# AI Customer Service Agent Guide

This file defines the project-specific coding style and working rules for this repository. Prefer following the existing architecture and naming in the codebase over generic Go conventions when they differ.

## Architecture

- Keep `domain/` for shared interfaces, DTOs, and API/domain error definitions.
- Keep concrete implementations under `internal/` by responsibility:
  - `internal/services/` for business logic
  - `internal/repositories/` for persistence
  - `internal/gin/` for HTTP routing and handlers
  - `internal/cache/`, `internal/database/`, `internal/jwt/` for infrastructure adapters
  - `internal/wire/` for dependency injection
- Keep `cmd/` as the process entrypoint only. Startup wiring belongs there, not business logic.

## Naming And Layout

- Follow the repository's current package and directory naming scheme: lower-case names, including underscore-separated package directories such as `openai_service`, `csai_service`, and `csai_repository`.
- Use `NewXxx` constructor names consistently.
- Service structs end with `Service`; repository structs end with `Repository`; interfaces in `domain/` end with `Interface`.
- Keep request/response structs near the HTTP handler that uses them.
- Use exported Go field names with explicit `json` tags for API payloads.
- Keep database table and column names in snake_case. GORM models should use explicit `gorm` tags and a `TableName()` method when needed.

## Dependency Injection

- Prefer constructor injection. Do not instantiate infrastructure dependencies directly inside handlers or service methods.
- Use `wire` for injector assembly in `internal/wire/`.
- After changing providers, constructors, or dependency graphs, regenerate wire output with `make wire`.
- Keep constructor signatures and injector return types aligned. Do not mix interface-returning providers with constructors that require unrelated concrete types.

## Interfaces And Context

- Keep service and repository contracts defined in `domain/`.
- Prefer depending on domain interfaces across layers instead of concrete implementations.
- When a handler call can be canceled or timed out, include `context.Context` in service and repository method signatures and pass `c.Request.Context()` downward.
- Keep domain interfaces and implementation method signatures identical. If one side uses `context.Context`, the other side must too.

## Error Handling And Logging

- For infrastructure and service errors, use `github.com/pkg/errors` and wrap returned errors with `errors.WithStack`.
- Keep transport-facing error payloads normalized through `domain.ErrorFormat` and the handler helpers in `internal/gin/handler/helpers.go`.
- Handlers should log internal errors and return normalized API errors instead of leaking raw dependency errors to clients.
- Prefer early returns for validation and error paths.

## HTTP Layer

- Keep Gin handlers thin:
  - bind and validate request payloads
  - normalize input with `strings.TrimSpace` where applicable
  - call a service
  - convert the result to a consistent JSON response
- **Always** use `Success()` and `Failed()` response helpers for all handler responses. Never call `c.JSON()` directly in handlers — all responses must go through these helpers to maintain the standard response shape:
  - Success: `{"code": 200, "data": <payload>}`
  - Failed: `{"code": <err.Code>, "message": <err.Message>}`
- Route registration stays in `internal/gin/routes/`; handlers should not build routers.

## Config And Globals

- Read environment configuration through `internal/configs.NewConfig()` only. Do not scatter direct `os.Getenv` access across the codebase.
- Shared singletons already follow `sync.Once` patterns in config, logger, database, cache, and OpenAI setup. Preserve that pattern when extending those components.
- Keep environment variable names uppercase with explicit defaults in config.

## Style Preferences

- Prefer small functions and early exits over nested control flow.
- Add comments only when the intent is not obvious from the code.
- Keep imports grouped by standard library, internal packages, and third-party packages as produced by `gofmt`.
- Preserve ASCII unless an existing file already uses Chinese text for user-facing content.

## Verification

- Run `gofmt` on edited Go files.
- Run `go test ./...` after meaningful code changes when the package graph is expected to compile.
- If dependency wiring changed, run `make wire` before tests.
