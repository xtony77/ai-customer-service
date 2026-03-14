package handler

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/logger"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Question string `json:"question"`
}

type ChatResponse struct {
	Answer string `json:"answer"`
}

func Ask(c *gin.Context) {
	var request ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		Failed(c, domain.ErrorBadRequest)
		return
	}

	request.Question = strings.TrimSpace(request.Question)
	if request.Question == "" {
		Failed(c, domain.ErrorBadRequest)
		return
	}

	svc := newCSAIService()
	answer, err := svc.Ask(request.Question)
	if err != nil {
		if errors.Is(err, domain.ErrPromptGuardRejected) {
			Failed(c, domain.ErrorBadRequest)
			return
		}
		if logger.Slog != nil {
			logger.Slog.Error("ai ask failed", "error", err)
		}
		Failed(c, domain.ErrorServer)
		return
	}

	Success(c, ChatResponse{Answer: answer})
}
