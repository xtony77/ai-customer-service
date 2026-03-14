package routes

import (
	"ai-customer-service/internal/configs"
	"ai-customer-service/internal/gin/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(config *configs.Config) *gin.Engine {
	if config.Gin.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api/v1")
	{
		api.POST("/qa", handler.QABatchCreate)
		api.POST("/chat", handler.Ask)
	}

	return r
}
