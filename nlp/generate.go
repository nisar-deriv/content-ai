package nlp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const apiKeyEnv = "OPENAI_API_KEY"

type OpenAIRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func EnhanceText(text string) (string, error) {
	apiKey := os.Getenv(apiKeyEnv)
	if apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not set")
	}

	prompt := fmt.Sprintf("Enhance the following team update with summaries and insights:\n\n%s", text)
	requestBody, _ := json.Marshal(OpenAIRequest{
		Model:     "text-davinci-003",
		Prompt:    prompt,
		MaxTokens: 150,
	})
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Choices[0].Text, nil
}
