package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/regentmarkets/ContentAI/config"
	"github.com/slack-go/slack"
	"gopkg.in/yaml.v2"
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

func (s *SlackHandler) FetchMessagesFromChannels(teamNames []string) (map[string][]string, error) { // Change return type
	log.Printf("Fetching messages for teams: %v", teamNames)
	messages := make(map[string][]string) // Change map type
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))

	for _, teamName := range teamNames {
		channelID, exists := config.SlackChannelIDs[teamName]
		if !exists {
			log.Printf("No channel ID found for team: %s", teamName)
			continue
		}

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
				messages[channelID] = append(messages[channelID], message.Text) // Append only Text field
			}
		}
	}
	log.Printf("Filtered messages ready for processing: %v", messages)
	return messages, nil
}

func FetchUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchUpdatesHandler triggered")
	slackHandler := NewSlackHandler()

	teamNames := make([]string, 0, len(config.SlackChannelIDs))
	for teamName := range config.SlackChannelIDs {
		teamNames = append(teamNames, teamName)
	}

	messages, err := slackHandler.FetchMessagesFromChannels(teamNames)
	if err != nil {
		log.Printf("Failed to fetch messages: %v", err)
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	for teamName, channelID := range config.SlackChannelIDs {
		teamMessages, ok := messages[channelID]
		if !ok {
			log.Printf("No messages found for team %s in channel %s", teamName, channelID)
			continue
		}

		log.Printf("Fetched messages for team %s: %v", teamName, teamMessages)

		weekFolder := getWeekFolder()
		if err := ensureDirectory(weekFolder); err != nil {
			log.Printf("Error creating directory %s for team %s: %v", weekFolder, teamName, err)
			http.Error(w, fmt.Sprintf("Error creating directory for team %s: %v", teamName, err), http.StatusInternalServerError)
			continue
		}

		yamlData, err := yaml.Marshal(teamMessages)
		if err != nil {
			log.Printf("Error marshaling YAML for team %s: %v", teamName, err)
			http.Error(w, fmt.Sprintf("Error marshaling YAML for team %s: %v", teamName, err), http.StatusInternalServerError)
			continue
		}

		filename := fmt.Sprintf("%s/%s.yaml", weekFolder, teamName)
		if err := os.WriteFile(filename, yamlData, 0644); err != nil {
			log.Printf("Error writing to file %s for team %s: %v", filename, teamName, err)
			http.Error(w, fmt.Sprintf("Error writing to file for team %s: %v", teamName, err), http.StatusInternalServerError)
			continue
		}

		log.Printf("Update processed and stored successfully for team %s in %s", teamName, filename)
	}

	fmt.Fprintf(w, "Updates processed and stored successfully for all teams")
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
