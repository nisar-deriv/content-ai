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
	if token == "" {
		log.Println("SLACK_API_TOKEN is not set")
	}
	client := slack.New(token)
	log.Println("Slack client initialized")
	return &SlackHandler{client: client}
}

func (s *SlackHandler) FetchMessagesFromChannels(channelIDs []string) (map[string][]slack.Message, error) {
	log.Printf("Fetching messages from channels: %v", channelIDs)
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
		log.Printf("Messages fetched for channel %s: %d messages", channelID, len(history.Messages))
		for _, message := range history.Messages {
			if strings.HasPrefix(message.Text, "weekly update") {
				messages[channelID] = append(messages[channelID], message)
			}
		}
	}
	log.Printf("Filtered messages ready for processing: %v", messages)
	return messages, nil
}

func FetchUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchUpdatesHandler triggered")
	slackHandler := NewSlackHandler()
	messages, err := slackHandler.FetchMessagesFromChannels(config.SlackChannelIDs)
	if err != nil {
		log.Printf("Failed to fetch messages: %v", err)
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}
	// Print messages here for debugging purposes
	log.Println("Fetched messages:", messages)

	// Serialize messages to JSON
	jsonData, err := json.Marshal(messages)
	if err != nil {
		log.Printf("Failed to serialize messages: %v", err)
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
	log.Println("Parsing team name from text")
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "| Team:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				teamName := strings.TrimSpace(parts[1])
				log.Printf("Team name found: %s", teamName)
				return teamName, nil
			}
		}
	}
	log.Println("Team name not found in the text")
	return "", fmt.Errorf("team name not found")
}

func getWeekFolder() string {
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
	weekEnd := weekStart.AddDate(0, 0, 4)
	weekFolder := fmt.Sprintf("Week %s to %s", weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))
	log.Printf("Week folder calculated: %s", weekFolder)
	return weekFolder
}

func ensureDirectory(path string) error {
	log.Printf("Checking directory %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Creating directory %s", path)
		return os.Mkdir(path, 0755)
	}
	log.Printf("Directory already exists: %s", path)
	return nil
}
