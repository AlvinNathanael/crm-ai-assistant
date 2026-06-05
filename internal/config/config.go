package config

import "os"

type Config struct {
	OpenAIAPIKey string
}

func Load() Config {
	return Config{
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
	}
}
