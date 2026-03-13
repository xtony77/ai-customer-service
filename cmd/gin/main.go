package main

import (
	"ai-customer-service/internal/configs"
	"ai-customer-service/internal/gin/routes"
	"fmt"
)

func main() {
	config := configs.NewConfig()
	ginServer := routes.NewRouter(config)
	ginServer.Run(fmt.Sprintf("%s:%s", config.Gin.Host, config.Gin.Port))
}
