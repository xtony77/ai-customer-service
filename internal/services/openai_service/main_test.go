package openai_service

import (
	"ai-customer-service/domain"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAIServiceGenerateEmbedding(t *testing.T) {
	service := newTestOpenAIService(func(request *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodPost, request.Method)
		assert.Equal(t, "/v1/embeddings", request.URL.Path)
		return jsonResponse(t, map[string]any{
			"object": "list",
			"model":  "text-embedding-3-small",
			"data": []map[string]any{
				{
					"object":    "embedding",
					"embedding": []float32{0.1, 0.2, 0.3},
					"index":     0,
				},
			},
			"usage": map[string]any{
				"prompt_tokens": 1,
				"total_tokens":  1,
			},
		}), nil
	})

	embedding, err := service.GenerateEmbedding("營業時間是？")

	require.NoError(t, err)
	assert.Equal(t, []float32{0.1, 0.2, 0.3}, embedding)
}

func TestOpenAIServiceChatCompletion(t *testing.T) {
	service := newTestOpenAIService(func(request *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodPost, request.Method)
		assert.Equal(t, "/v1/chat/completions", request.URL.Path)
		return jsonResponse(t, map[string]any{
			"id":      "chatcmpl-test",
			"object":  "chat.completion",
			"created": 1,
			"model":   "gpt-4o-mini",
			"choices": []map[string]any{
				{
					"index": 0,
					"message": map[string]any{
						"role":    "assistant",
						"content": "依據知識庫資料，七天內可退貨。",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]any{
				"prompt_tokens":     1,
				"completion_tokens": 1,
				"total_tokens":      2,
			},
		}), nil
	})

	answer, err := service.ChatCompletion("可以退貨嗎？", []domain.QAEmbedding{
		{Question: "可以退貨嗎？", Answer: "七天內可退貨"},
	})

	require.NoError(t, err)
	assert.Equal(t, "依據知識庫資料，七天內可退貨。", answer)
}

func newTestOpenAIService(do func(request *http.Request) (*http.Response, error)) *OpenAIService {
	config := openai.DefaultConfig("test-token")
	config.BaseURL = "http://openai.test/v1"
	config.HTTPClient = stubHTTPClient{do: do}

	return &OpenAIService{
		Client:         openai.NewClientWithConfig(config),
		Model:          "gpt-4o-mini",
		EmbeddingModel: "text-embedding-3-small",
	}
}

type stubHTTPClient struct {
	do func(request *http.Request) (*http.Response, error)
}

func (s stubHTTPClient) Do(request *http.Request) (*http.Response, error) {
	return s.do(request)
}

func jsonResponse(t *testing.T, payload any) *http.Response {
	t.Helper()

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(body)),
	}
}
