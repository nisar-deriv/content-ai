package config

import (
	"os"
)

var SlackChannelIDs = map[string]string{
	"cloudplatform": "C07F60QARM3", // test-cloudplatform
	"k8s":           "C07F8PXH1CK", // test-k8s
	"prod":          "C07EUAWCFB9", // test-prod
	"winops":        "C07F8PW8LJX", // test-winops

	// Add more team name to channel ID mappings as needed
}

type Config struct {
	UseOllama     bool
	OpenAIKey     string
	GitHubRepoURL string
}

func LoadConfig() Config {
	useOllama := os.Getenv("USE_OLLAMA") == "true"
	openAIKey := os.Getenv("OPENAI_API_KEY")
	gitHubRepoURL := os.Getenv("GITHUB_REPO_URL")

	return Config{
		UseOllama:     useOllama,
		OpenAIKey:     openAIKey,
		GitHubRepoURL: gitHubRepoURL,
	}
}
