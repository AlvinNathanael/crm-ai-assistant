package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIAPIKey    string
	OpenAIChatModel string
	PostgresDSN     string

	DocumentConfig *DocumentConfig
}

type DocumentConfig struct {
	DocumentChunkSize    int
	DocumentChunkOverlap int
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		OpenAIAPIKey:    os.Getenv("OPENAI_API_KEY"),
		OpenAIChatModel: getEnv("OPENAI_CHAT_MODEL", "gpt-6-luna"),
		PostgresDSN:     os.Getenv("POSTGRES_DSN"),

		// Document
		DocumentConfig: &DocumentConfig{
			DocumentChunkSize:    getEnvAsInt("DOCUMENT_CHUNK_SIZE", 1000),
			DocumentChunkOverlap: getEnvAsInt("DOCUMENT_CHUNK_OVERLAP", 150),
		},
	}
}

func getEnv(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
