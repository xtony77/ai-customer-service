## ADDED Requirements

### Requirement: Batch Q&A import endpoint
The system SHALL expose `POST /api/v1/qa` that accepts a JSON body with an `items` array, each item containing `question` (string, required) and `answer` (string, required).

#### Scenario: Valid batch import
- **WHEN** client sends `POST /api/v1/qa` with body `{"items": [{"question": "Q1", "answer": "A1"}, {"question": "Q2", "answer": "A2"}]}`
- **THEN** system returns HTTP 200 with JSON body `{"code": 200, "data": {"count": 2}}`

#### Scenario: Empty items array
- **WHEN** client sends `POST /api/v1/qa` with empty `items` array
- **THEN** system returns HTTP 400 with JSON body `{"code": 400, "message": "bad request"}`

#### Scenario: Missing or empty question in items
- **WHEN** any item in the array has an empty or whitespace-only `question`
- **THEN** system returns HTTP 400 with JSON body `{"code": 400, "message": "bad request"}`

#### Scenario: Missing or empty answer in items
- **WHEN** any item in the array has an empty or whitespace-only `answer`
- **THEN** system returns HTTP 400 with JSON body `{"code": 400, "message": "bad request"}`

### Requirement: Server-side embedding generation
The system SHALL call OpenAI Embedding API to generate a 1536-dimension embedding for each Q&A item's `question` field before storing.

#### Scenario: Embeddings generated for all items
- **WHEN** all items' embeddings are generated successfully
- **THEN** system inserts all items (question, answer, embedding) into `qa_embeddings` table in batches of 100

#### Scenario: OpenAI Embedding API fails during import
- **WHEN** OpenAI Embedding API returns an error for any item
- **THEN** system returns HTTP 500 with JSON body `{"code": 500, "message": "Server Error"}`

### Requirement: Store Q&A with embedding in pgvector
The system SHALL insert each Q&A item into the `qa_embeddings` table with columns: question, answer, embedding (vector(1536)), created_at, updated_at.

#### Scenario: Successful storage
- **WHEN** all embeddings are generated and database insert succeeds
- **THEN** all Q&A items are persisted in `qa_embeddings` with HNSW-indexed embeddings

### Requirement: No authentication required
The endpoint SHALL be publicly accessible without JWT authentication.

#### Scenario: Request without auth header
- **WHEN** client sends `POST /api/v1/qa` without Authorization header
- **THEN** system processes the request normally
