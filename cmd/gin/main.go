package main

import (
	"ai-customer-service/internal/configs"
	"ai-customer-service/internal/gin/routes"
	"ai-customer-service/internal/logger"
	"fmt"
)

func main() {
	config := configs.NewConfig()
	logger.NewSlog(config)
	ginServer := routes.NewRouter(config)
	ginServer.Run(fmt.Sprintf("%s:%s", config.Gin.Host, config.Gin.Port))
}
