package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/regentmarkets/ContentAI/data"
	"github.com/regentmarkets/ContentAI/nlp"
)

type DetailedPayload struct {
	Progress string `json:"progress"`
	Problems string `json:"problems"`
	Plan     string `json:"plan"`
	Insights string `json:"insights"`
}

func ReportGenerationHandler(w http.ResponseWriter, r *http.Request) {
	err := GenerateWeeklyReports()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating weekly reports: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Weekly reports generated successfully")
}

func GenerateWeeklyReports() error {
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

			parsedData, err := parseDetailedPayload(content)
			if err != nil {
				return fmt.Errorf("error parsing payload from file %s: %v", file.Name(), err)
			}

			enhancedContent, err := enhanceAndFormatContent(file.Name(), parsedData)
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

func enhanceAndFormatContent(filename string, data DetailedPayload) (string, error) {
	enhancedProgress, err := enhanceText(data.Progress)
	if err != nil {
		return "", fmt.Errorf("error enhancing progress text: %v", err)
	}
	enhancedProblems, err := enhanceText(data.Problems)
	if err != nil {
		return "", fmt.Errorf("error enhancing problems text: %v", err)
	}
	enhancedPlan, err := enhanceText(data.Plan)
	if err != nil {
		return "", fmt.Errorf("error enhancing plan text: %v", err)
	}
	enhancedInsights, err := enhanceText(data.Insights)
	if err != nil {
		return "", fmt.Errorf("error enhancing insights text: %v", err)
	}

	return createEnhancedContent(
		filename,
		enhancedProgress,
		enhancedProblems,
		enhancedPlan,
		enhancedInsights,
	), nil
}

func enhanceText(text string) (string, error) {
	if cfg.UseOllama {
		return nlp.EnhanceTextWithOllama(text)
	}
	return nlp.EnhanceTextWithOpenAI(text, cfg.OpenAIKey)
}

func parseDetailedPayload(content string) (DetailedPayload, error) {
	var payload DetailedPayload
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Progress:") {
			payload.Progress = strings.TrimSpace(strings.TrimPrefix(line, "Progress:"))
		} else if strings.HasPrefix(line, "Problems:") {
			payload.Problems = strings.TrimSpace(strings.TrimPrefix(line, "Problems:"))
		} else if strings.HasPrefix(line, "Plan:") {
			payload.Plan = strings.TrimSpace(strings.TrimPrefix(line, "Plan:"))
		} else if strings.HasPrefix(line, "Insights:") {
			payload.Insights = strings.TrimSpace(strings.TrimPrefix(line, "Insights:"))
		}
	}
	return payload, nil
}

func createEnhancedContent(filename, progress, problems, plan, insights string) string {
	return fmt.Sprintf(`
        Team: %s
        Progress:
        %s
        
        Problems:
        %s
        
        Plan:
        %s
        
        Insights:
        %s
    `, filename, progress, problems, plan, insights)
}
