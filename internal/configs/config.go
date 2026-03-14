package configs

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Gin        *Gin
	JWT        *JWT
	PostgreSQL *PostgreSQL
	Redis      *Redis
	OpenAI     *OpenAI
}

type Gin struct {
	Host string
	Port string
	Mode string
}

type JWT struct {
	Secret string
}

type PostgreSQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Timezone string
}

type Redis struct {
	Host     string
	Port     string
	Password string
	Database int
}

type OpenAI struct {
	APIKey         string
	Model          string
	EmbeddingModel string
}

var config *Config
var configOnce sync.Once

func NewConfig() *Config {
	configOnce.Do(func() {
		_ = godotenv.Load()

		config = &Config{
			Gin: &Gin{
				Host: getEnv("GIN_HOST", "0.0.0.0"),
				Port: getEnv("GIN_PORT", "8080"),
				Mode: getEnv("GIN_MODE", "release"),
			},
			JWT: &JWT{
				Secret: getEnv("JWT_SECRET", ""),
			},
			PostgreSQL: &PostgreSQL{
				Host:     getEnv("POSTGRESQL_HOST", "0.0.0.0"),
				Port:     getEnv("POSTGRESQL_PORT", "5432"),
				Username: getEnv("POSTGRESQL_USER", ""),
				Password: getEnv("POSTGRESQL_PASSWORD", ""),
				Database: getEnv("POSTGRESQL_DATABASE", ""),
				Timezone: getEnv("POSTGRESQL_TIMEZONE", "Etc/UTC"),
			},
			Redis: &Redis{
				Host:     getEnv("REDIS_HOST", "0.0.0.0"),
				Port:     getEnv("REDIS_PORT", "6379"),
				Password: getEnv("REDIS_PASSWORD", ""),
				Database: getEnvInt("REDIS_DATABASE", 0),
			},
			OpenAI: &OpenAI{
				APIKey:         getEnv("OPENAI_API_KEY", ""),
				Model:          getEnv("OPENAI_MODEL", "gpt-4o-mini"),
				EmbeddingModel: getEnv("OPENAI_EMBEDDING_MODEL", "text-embedding-3-small"),
			},
		}
	})
	return config
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		res, _ := strconv.ParseInt(val, 10, 32)
		return int(res)
	}
	return defaultVal
}
