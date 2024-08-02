package config

import (
	"os"
)

var SlackChannelIDs = []string{
	"C07F60QARM3", // test-cloudplatform
	"C07F8PXH1CK", // test-k8s
	"C07EUAWCFB9", // test-prod
	"C07F8PW8LJX", // test-winops

	// Add more channel IDs as needed
}

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
