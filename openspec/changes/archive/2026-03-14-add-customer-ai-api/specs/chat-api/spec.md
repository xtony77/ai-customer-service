## ADDED Requirements

### Requirement: Chat endpoint accepts user question
The system SHALL expose `POST /api/v1/chat` that accepts a JSON body with a `question` field (string, required).

#### Scenario: Valid question submitted
- **WHEN** client sends `POST /api/v1/chat` with body `{"question": "營業時間是？"}`
- **THEN** system returns HTTP 200 with JSON body `{"code": 200, "data": {"answer": "..."}}`

#### Scenario: Missing question field
- **WHEN** client sends `POST /api/v1/chat` with empty body or missing `question`
- **THEN** system returns HTTP 400 with JSON body `{"code": 400, "message": "bad request"}`

#### Scenario: Empty question after trimming whitespace
- **WHEN** client sends `POST /api/v1/chat` with `question` containing only whitespace
- **THEN** system returns HTTP 400 with JSON body `{"code": 400, "message": "bad request"}`

### Requirement: Redis exact-match cache lookup
The system SHALL compute `sha256(question)` as cache key and check Redis before any AI API calls.

#### Scenario: Cache hit
- **WHEN** the question's sha256 hash exists as a Redis key
- **THEN** system returns the cached answer directly without calling OpenAI APIs

#### Scenario: Cache miss
- **WHEN** the question's sha256 hash does not exist in Redis
- **THEN** system proceeds to generate embedding and perform vector search

### Requirement: Generate embedding on cache miss
The system SHALL call OpenAI Embedding API (model from `OPENAI_EMBEDDING_MODEL` env, default `text-embedding-3-small`) to generate a 1536-dimension embedding for the user question.

#### Scenario: Embedding generated successfully
- **WHEN** cache miss occurs and OpenAI Embedding API returns successfully
- **THEN** system uses the embedding to query pgvector for top-3 similar Q&A

#### Scenario: OpenAI Embedding API fails
- **WHEN** OpenAI Embedding API returns an error
- **THEN** system returns HTTP 500 with JSON body `{"code": 500, "message": "Server Error"}`

### Requirement: Vector search top-3 Q&A via pgvector
The system SHALL query `qa_embeddings` table using cosine distance with HNSW index, returning the top 3 most similar Q&A pairs.

#### Scenario: Top-3 results found
- **WHEN** pgvector query returns results
- **THEN** system uses the top-3 question-answer pairs as context for GPT

#### Scenario: No Q&A data exists
- **WHEN** `qa_embeddings` table is empty or no results found
- **THEN** system returns HTTP 200 with `{"code": 200, "data": {"answer": "目前沒有可用的知識庫資料"}}` without calling GPT

### Requirement: GPT-4o mini generates answer
The system SHALL call OpenAI Chat Completion API (model from `OPENAI_MODEL` env, default `gpt-4o-mini`) with a system prompt that instructs: only answer based on the provided top-3 reference Q&A, limit response to 150 Chinese characters (150字), refuse to answer if references lack relevant information.

#### Scenario: GPT generates answer successfully
- **WHEN** OpenAI Chat Completion API returns successfully
- **THEN** system stores the answer in Redis (key = sha256(question), no TTL) and returns HTTP 200 with `{"code": 200, "data": {"answer": "..."}}`

#### Scenario: OpenAI Chat Completion API fails
- **WHEN** OpenAI Chat Completion API returns an error
- **THEN** system returns HTTP 500 with JSON body `{"code": 500, "message": "Server Error"}`

### Requirement: Cache answer permanently
The system SHALL store the GPT-generated answer in Redis with no TTL (permanent cache, expiration = 0).

#### Scenario: Answer cached after generation
- **WHEN** GPT successfully generates an answer
- **THEN** system stores it in Redis with key `sha256(question)` and expiration 0 (no expiration)
