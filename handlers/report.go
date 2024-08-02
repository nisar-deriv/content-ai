package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/regentmarkets/ContentAI/config"
	"github.com/regentmarkets/ContentAI/data"
	"github.com/regentmarkets/ContentAI/nlp"
)

var cfg config.Config

func InitHandlers() {
	log.Println("Loading configuration")
	cfg = config.LoadConfig() // Corrected to call without parameters
}

type DetailedPayload struct {
	Progress string `json:"progress"`
	Problems string `json:"problems"`
	Plan     string `json:"plan"`
	Insights string `json:"insights"`
}

func ReportGenerationHandlerAi(w http.ResponseWriter, r *http.Request) {
	err := GenerateWeeklyReportsAi()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating weekly reports: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Weekly reports generated successfully")
}

func GenerateWeeklyReportsAi() error {
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
	weekEnd := weekStart.AddDate(0, 0, 4)
	weekFolder := fmt.Sprintf("Week %s to %s", weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))

	files, err := os.ReadDir(weekFolder)
	if err != nil {
		return fmt.Errorf("error reading week folder: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			content, err := data.ReadFromFile(fmt.Sprintf("%s/%s", weekFolder, file.Name()))
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", file.Name(), err)
			}

			enhancedContent, err := enhanceFullContent(content)
			if err != nil {
				return fmt.Errorf("error enhancing content for file %s: %v", file.Name(), err)
			}

			enhancedFilename := fmt.Sprintf("%s/enhanced_%s", weekFolder, file.Name())
			err = data.WriteToFile(enhancedFilename, enhancedContent)
			if err != nil {
				return fmt.Errorf("error writing enhanced content to file %s: %v", enhancedFilename, err)
			}
		}
	}

	return nil
}

func enhanceFullContent(content string) (string, error) {
	if cfg.UseOllama {
		return nlp.EnhanceTextWithOllama(content)
	}
	return nlp.EnhanceTextWithOpenAI(content, cfg.OpenAIKey)
}
