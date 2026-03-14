package csai_service

import (
	"ai-customer-service/domain"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const noKnowledgeBaseMessage = "我無法回答這個問題，請聯繫人工客服處理。"

type CSAIService struct {
	Redis    domain.RedisInterface
	OpenAI   domain.OpenAIInterface
	CSAIRepo domain.CSAIRepositoryInterface
}

func NewCSAIService(redis domain.RedisInterface, openAIClient domain.OpenAIInterface, csAIRepo domain.CSAIRepositoryInterface) domain.CSAIServiceInterface {
	return &CSAIService{
		Redis:    redis,
		OpenAI:   openAIClient,
		CSAIRepo: csAIRepo,
	}
}

func (s *CSAIService) Ask(question string) (string, error) {
	question = strings.TrimSpace(question)
	cacheKey := hashQuestion(question)

	if err := promptValidate(question); err != nil {
		return "", err
	}

	cached, err := s.Redis.Get(cacheKey)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if cached != "" {
		return cached, nil
	}

	embedding, err := s.OpenAI.GenerateEmbedding(question)
	if err != nil {
		return "", errors.WithStack(err)
	}

	references, err := s.CSAIRepo.SearchSimilar(embedding, 3)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if len(references) == 0 {
		return noKnowledgeBaseMessage, nil
	}

	answer, err := s.OpenAI.ChatCompletion(question, references)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if err := s.Redis.Set(cacheKey, answer, 0); err != nil {
		return "", errors.WithStack(err)
	}

	return answer, nil
}

func promptValidate(question string) error {
	keywords := []string{
		"ignore previous instructions",
		"ignore all previous instructions",
		"system prompt",
		"developer message",
		"reveal your prompt",
		"你是什麼模型",
		"忽略之前的指令",
		"系統提示",
	}

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)ignore\s+(all\s+)?previous\s+instructions`),
		regexp.MustCompile(`(?i)system\s+prompt`),
		regexp.MustCompile(`(?i)developer\s+message`),
		regexp.MustCompile(`(?i)reveal\s+(your|the)\s+prompt`),
		regexp.MustCompile(`(?i)(你是什麼模型|忽略之前的指令|系統提示)`),
	}

	normalized := strings.ToLower(question)
	for _, keyword := range keywords {
		if strings.Contains(normalized, strings.ToLower(keyword)) {
			return errors.WithStack(domain.ErrPromptGuardRejected)
		}
	}

	for _, pattern := range patterns {
		if pattern.MatchString(question) {
			return errors.WithStack(domain.ErrPromptGuardRejected)
		}
	}

	return nil
}

func (s *CSAIService) QABatchCreate(items []domain.QAItem) (int, error) {
	if len(items) == 0 {
		return 0, nil
	}

	embeddings := make([]domain.QAEmbedding, 0, len(items))
	for _, item := range items {
		embedding, err := s.OpenAI.GenerateEmbedding(item.Question)
		if err != nil {
			return 0, errors.WithStack(err)
		}

		embeddings = append(embeddings, domain.QAEmbedding{
			Question:  item.Question,
			Answer:    item.Answer,
			Embedding: embedding,
		})
	}

	if err := s.CSAIRepo.BatchInsert(embeddings); err != nil {
		return 0, errors.WithStack(err)
	}

	return len(embeddings), nil
}

func hashQuestion(question string) string {
	hash := sha256.Sum256([]byte(question))
	return hex.EncodeToString(hash[:])
}
