package logger

import (
	"ai-customer-service/internal/configs"
	"log/slog"
	"os"
	"strings"
	"sync"
)

var Slog *slog.Logger
var slogOnce sync.Once

func NewSlog(config *configs.Config) *slog.Logger {
	slogOnce.Do(func() {
		var handler slog.Handler

		switch strings.ToLower(config.Gin.Mode) {
		case "debug":
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			})

		default:
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
		}

		Slog = slog.New(handler)
	})

	return Slog
}
