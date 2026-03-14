## ADDED Requirements

### Requirement: Keyword-based prompt injection filter
The system SHALL check the user question against a list of known prompt injection patterns (keywords and regex) before processing. The check SHALL happen before cache lookup.

Blocked keywords (case-insensitive substring match):
- "ignore previous instructions"
- "ignore all previous instructions"
- "system prompt"
- "developer message"
- "reveal your prompt"
- "你是什麼模型"
- "忽略之前的指令"
- "系統提示"

Blocked regex patterns (case-insensitive):
- `(?i)ignore\s+(all\s+)?previous\s+instructions`
- `(?i)system\s+prompt`
- `(?i)developer\s+message`
- `(?i)reveal\s+(your|the)\s+prompt`
- `(?i)(你是什麼模型|忽略之前的指令|系統提示)`

#### Scenario: Injection pattern detected
- **WHEN** user question matches any keyword or regex pattern
- **THEN** system returns HTTP 400 with JSON body `{"code": 400, "message": "bad request"}` and does NOT proceed to cache lookup or AI API calls

#### Scenario: Clean question passes filter
- **WHEN** user question does not match any injection patterns
- **THEN** system proceeds normally to cache lookup

### Requirement: AI system prompt constrains responses
The system SHALL include a system prompt that instructs the AI to: (1) only answer based on the provided top-3 reference Q&A, (2) limit response to 150 Chinese characters (150字), (3) refuse to answer questions unrelated to the reference material by replying that it cannot answer based on available information.

#### Scenario: Question related to reference material
- **WHEN** GPT receives a question with relevant top-3 context
- **THEN** GPT generates an answer based solely on the provided reference Q&A within 150 Chinese characters

#### Scenario: Question unrelated to reference material
- **WHEN** GPT receives a question with no relevant context in top-3 results
- **THEN** GPT responds indicating it cannot answer the question based on available information
