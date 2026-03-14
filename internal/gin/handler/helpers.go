package handler

import (
	"ai-customer-service/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

func Failed(c *gin.Context, err domain.ErrorFormat) {
	c.JSON(err.HttpStatus, gin.H{
		"code":    err.Code,
		"message": err.Message,
	})
}
