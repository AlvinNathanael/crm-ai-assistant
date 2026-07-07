package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIAPIKey string
	PostgresDSN  string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
		PostgresDSN:  os.Getenv("POSTGRES_DSN"),
	}
}
