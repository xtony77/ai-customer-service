package routes

import (
	"ai-customer-service/internal/configs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(config *configs.Config) *gin.Engine {
	if config.Gin.Mode == "RELEASE" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "pong")
		})
	}

	return r
}
