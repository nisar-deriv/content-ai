package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/regentmarkets/ContentAI/config"
	"github.com/slack-go/slack"
)

type SlackHandler struct {
	client *slack.Client
}

func NewSlackHandler() *SlackHandler {
	token := os.Getenv("SLACK_API_TOKEN")
	client := slack.New(token)
	return &SlackHandler{client: client}
}

func (s *SlackHandler) FetchMessagesFromChannels(channelIDs []string) (map[string][]slack.Message, error) {
	messages := make(map[string][]slack.Message)
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))

	for _, channelID := range channelIDs {
		historyParams := slack.GetConversationHistoryParameters{
			ChannelID: channelID,
			Limit:     100,
			Oldest:    strconv.FormatInt(startOfWeek.Unix(), 10), // Convert Unix timestamp to string
		}
		history, err := s.client.GetConversationHistory(&historyParams)
		if err != nil {
			log.Printf("Error fetching messages from channel %s: %v", channelID, err)
			return nil, err
		}
		for _, message := range history.Messages {
			if message.Text == "weekly update" {
				messages[channelID] = append(messages[channelID], message)
			}
		}
	}
	return messages, nil
}

func FetchUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	slackHandler := NewSlackHandler()
	messages, err := slackHandler.FetchMessagesFromChannels(config.SlackChannelIDs)
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	// Serialize messages to JSON
	jsonData, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, "Failed to serialize messages", http.StatusInternalServerError)
		return
	}

	// Parse team name from messages
	teamName, err := parseTeamName(string(jsonData))
	if err != nil {
		log.Printf("Error parsing team name: %v", err)
		http.Error(w, fmt.Sprintf("Error parsing team name: %v", err), http.StatusBadRequest)
		return
	}

	// Create directory for the current week
	weekFolder := getWeekFolder()
	if err := ensureDirectory(weekFolder); err != nil {
		log.Printf("Error creating directory %s: %v", weekFolder, err)
		http.Error(w, fmt.Sprintf("Error creating directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Write JSON data to a file
	filename := fmt.Sprintf("%s/%s.json", weekFolder, teamName)
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
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
