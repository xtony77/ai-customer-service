package openai_service

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/configs"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	Client         *openai.Client
	Model          string
	EmbeddingModel string
}

var conn *openai.Client
var openaiSvcOnce sync.Once

func NewOpenAI(config *configs.Config) domain.OpenAIInterface {
	openaiSvcOnce.Do(func() {
		conn = openai.NewClient(config.OpenAI.APIKey)
	})

	return &OpenAIService{
		Client:         conn,
		Model:          config.OpenAI.Model,
		EmbeddingModel: config.OpenAI.EmbeddingModel,
	}
}

func (c *OpenAIService) GenerateEmbedding(input string) ([]float32, error) {
	resp, err := c.Client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: []string{input},
		Model: openai.EmbeddingModel(c.EmbeddingModel),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(resp.Data) == 0 {
		return nil, errors.WithStack(fmt.Errorf("openai embedding response is empty"))
	}

	return resp.Data[0].Embedding, nil
}

func (c *OpenAIService) ChatCompletion(question string, references []domain.QAEmbedding) (string, error) {
	resp, err := c.Client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: c.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: buildSystemPrompt(references),
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
	})
	if err != nil {
		return "", errors.WithStack(err)
	}
	if len(resp.Choices) == 0 {
		return "", errors.WithStack(fmt.Errorf("openai chat completion response is empty"))
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func buildSystemPrompt(references []domain.QAEmbedding) string {
	var builder strings.Builder

	builder.WriteString("You are an AI customer service assistant. ")
	builder.WriteString("Only answer based on the provided reference Q&A. ")
	builder.WriteString("Keep the answer within 150 Chinese characters. ")
	builder.WriteString("If the references do not contain enough information, reply that you cannot answer based on the available information.\n\n")
	builder.WriteString("Reference Q&A:\n")

	for idx, reference := range references {
		builder.WriteString(fmt.Sprintf("%d. Q: %s\nA: %s\n", idx+1, reference.Question, reference.Answer))
	}

	return builder.String()
}
