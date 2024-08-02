package nlp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Text string `json:"text"`
}

func EnhanceTextWithOpenAI(text, apiKey string) (string, error) {
	if apiKey == "" {
		log.Println("Error: OpenAI API key not provided")
		return "", fmt.Errorf("OpenAI API key not set")
	}

	prompt := fmt.Sprintf("Enhance the following team update with summaries and insights:\n\n%s", text)
	requestBody, err := json.Marshal(OpenAIRequest{
		Model:     "text-davinci-003",
		Prompt:    prompt,
		MaxTokens: 150,
	})
	if err != nil {
		log.Printf("Error marshalling OpenAI request: %v", err)
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating request to OpenAI: %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to OpenAI: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body from OpenAI: %v", err)
		return "", err
	}

	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshalling OpenAI response: %v", err)
		return "", err
	}

	return response.Choices[0].Text, nil
}

func EnhanceTextWithOllama(text string) (string, error) {
	prompt := fmt.Sprintf("Enhance the following team update with summaries and insights:\n\n%s", text)
	requestBody, err := json.Marshal(OllamaRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		log.Printf("Error marshalling Ollama request: %v", err)
		return "", err
	}

	// Log the request body
	log.Printf("Request body being sent: %s", string(requestBody))

	req, err := http.NewRequest("POST", "http://host.docker.internal:11434/api/generate", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating request to Ollama: %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to Ollama: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body from Ollama: %v", err)
		return "", err
	}

	// Log the response body
	log.Printf("Response body received: %s", string(body))

	var response OllamaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshalling Ollama response: %v", err)
		return "", err
	}

	return response.Text, nil
}
