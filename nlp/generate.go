package nlp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

type TeamUpdate struct {
	Team     string   `yaml:"Team"`
	Problems []string `yaml:"Problems"`
	Progress []string `yaml:"Progress"`
	Insights []string `yaml:"Insights"`
	Plans    []string `yaml:"Plans"`
}

// Function to enhance text using Ollama
func EnhanceTextWithOllama(text string) (string, error) {
	prompt := fmt.Sprintf("Please enhance the text under the sections Problems, Progress, Insights, and Plans by adding more detail and using more sophisticated language while strictly maintaining the original bullet point format. Each item should remain as a single concise statement without breaking into sub-points or detailed explanations. Focus on using precise and technical terms appropriate for an IT professional audience and final output should be a yaml and no sub list converstion \n\n%s", text)
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

	req, err := http.NewRequest("POST", "http://ollama:11434/api/generate", bytes.NewBuffer(requestBody))
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

	var response OllamaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshalling Ollama response: %v", err)
		return "", err
	}

	if response.Done {
		return response.Response, nil
	}

	return "", fmt.Errorf("incomplete response from Ollama")
}
