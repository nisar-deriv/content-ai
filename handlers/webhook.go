package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/regentmarkets/ContentAI/config"
	"github.com/regentmarkets/ContentAI/data"
	"github.com/regentmarkets/ContentAI/nlp"
)

var cfg config.Config

func InitHandlers() {
	cfg = config.LoadConfig()
}

type SlackPayload struct {
	Text string `json:"text"`
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload SlackPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	teamName, err := parseTeamName(payload.Text)
	if err != nil {
		http.Error(w, "Error parsing team name", http.StatusBadRequest)
		return
	}

	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
	weekEnd := weekStart.AddDate(0, 0, 4)
	weekFolder := fmt.Sprintf("Week %s to %s", weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))

	if _, err := os.Stat(weekFolder); os.IsNotExist(err) {
		err := os.Mkdir(weekFolder, 0755)
		if err != nil {
			http.Error(w, "Error creating directory", http.StatusInternalServerError)
			return
		}
	}

	var enhancedText string
	if cfg.UseOllama {
		enhancedText, err = nlp.EnhanceTextWithOllama(payload.Text)
	} else {
		enhancedText, err = nlp.EnhanceTextWithOpenAI(payload.Text, cfg.OpenAIKey)
	}
	if err != nil {
		http.Error(w, "Error processing text with NLP", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s/%s.txt", weekFolder, teamName)
	err = data.WriteToFile(filename, enhancedText)
	if err != nil {
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Update processed and stored successfully in %s", filename)
}

func parseTeamName(text string) (string, error) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Team") {
			return strings.Fields(line)[1], nil
		}
	}
	return "", fmt.Errorf("team name not found")
}
