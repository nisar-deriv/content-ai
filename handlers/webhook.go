package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/regentmarkets/ContentAI/config"
)

var cfg config.Config

func InitHandlers() {
	log.Println("Loading configuration")
	cfg = config.LoadConfig() // Corrected to call without parameters
}

type SlackPayload struct {
	Text string `json:"text"`
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s", r.Method)
	if r.Method != http.MethodPost {
		log.Printf("Invalid request method: %s", r.Method)
		http.Error(w, "Invalid request method. Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload SlackPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, fmt.Sprintf("Error parsing request body: %v", err), http.StatusBadRequest)
		return
	}

	teamName, err := parseTeamName(payload.Text)
	if err != nil {
		log.Printf("Error parsing team name: %v", err)
		http.Error(w, fmt.Sprintf("Error parsing team name: %v", err), http.StatusBadRequest)
		return
	}

	weekFolder := getWeekFolder()
	if err := ensureDirectory(weekFolder); err != nil {
		log.Printf("Error creating directory %s: %v", weekFolder, err)
		http.Error(w, fmt.Sprintf("Error creating directory: %v", err), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s/%s.txt", weekFolder, teamName)
	if err := os.WriteFile(filename, []byte(payload.Text), 0644); err != nil {
		log.Printf("Error writing to file %s: %v", filename, err)
		http.Error(w, fmt.Sprintf("Error writing to file: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Update processed and stored successfully in %s", filename)
	fmt.Fprintf(w, "Update processed and stored successfully in %s", filename)
}

func parseTeamName(text string) (string, error) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "Team") {
			return strings.Fields(line)[1], nil
		}
	}
	log.Println("Team name not found in the text")
	return "", fmt.Errorf("team name not found")
}

func getWeekFolder() string {
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
	weekEnd := weekStart.AddDate(0, 0, 4)
	return fmt.Sprintf("Week %s to %s", weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))
}

func ensureDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Creating directory %s", path)
		return os.Mkdir(path, 0755)
	}
	return nil
}
