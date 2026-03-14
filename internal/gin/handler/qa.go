package handler

import (
	"ai-customer-service/domain"
	"ai-customer-service/internal/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

type QABatchCreateRequest struct {
	Items []domain.QAItem `json:"items"`
}

type QABatchCreateResponse struct {
	Count int `json:"count"`
}

func QABatchCreate(c *gin.Context) {
	var request QABatchCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		Failed(c, domain.ErrorBadRequest)
		return
	}
	if len(request.Items) == 0 {
		Failed(c, domain.ErrorBadRequest)
		return
	}

	items := make([]domain.QAItem, 0, len(request.Items))
	for _, item := range request.Items {
		item.Question = strings.TrimSpace(item.Question)
		item.Answer = strings.TrimSpace(item.Answer)

		if item.Question == "" {
			Failed(c, domain.ErrorBadRequest)
			return
		}
		if item.Answer == "" {
			Failed(c, domain.ErrorBadRequest)
			return
		}

		items = append(items, item)
	}

	svc := newCSAIService()
	count, err := svc.QABatchCreate(items)
	if err != nil {
		if logger.Slog != nil {
			logger.Slog.Error("qa batch create failed", "error", err)
		}
		Failed(c, domain.ErrorServer)
		return
	}

	Success(c, QABatchCreateResponse{Count: count})
}
