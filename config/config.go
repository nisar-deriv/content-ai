package config

import (
	"os"
)

type Config struct {
	UseOllama bool
	OpenAIKey string
}

func LoadConfig() Config {
	useOllama := os.Getenv("USE_OLLAMA") == "true"
	openAIKey := os.Getenv("OPENAI_API_KEY")

	return Config{
		UseOllama: useOllama,
		OpenAIKey: openAIKey,
	}
}
