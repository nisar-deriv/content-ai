package handlers

import (
	"encoding/json"
	"fmt"
	"log"
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
	log.Println("Starting report generation")
	filename, err := GenerateWeeklyReport()
	if err != nil {
		log.Printf("Error generating weekly report: %v", err)
		http.Error(w, fmt.Sprintf("Error generating weekly report: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Weekly report generated successfully: %s", filename)
	fmt.Fprintf(w, "Weekly report generated successfully: %s", filename)
}

func GenerateWeeklyReport() (string, error) {
	weekFolder, err := prepareWeekFolder()
	if err != nil {
		return "", err
	}

	progress, problems, plan, insights, err := compileReportSections(weekFolder)
	if err != nil {
		return "", err
	}

	finalReport := createFinalReport(
		strings.Join(progress, "\n\n"),
		strings.Join(problems, "\n\n"),
		strings.Join(plan, "\n\n"),
		strings.Join(insights, "\n\n"),
	)
	finalReportFilename := fmt.Sprintf("%s/final_report.html", weekFolder)
	log.Printf("Writing the final report to %s", finalReportFilename)
	if err := data.WriteToFile(finalReportFilename, finalReport); err != nil {
		return "", fmt.Errorf("error writing final report: %v", err)
	}

	return finalReportFilename, nil
}

func prepareWeekFolder() (string, error) {
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
	weekEnd := weekStart.AddDate(0, 0, 4)
	weekFolder := fmt.Sprintf("Week %s to %s", weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))

	if _, err := os.Stat(weekFolder); os.IsNotExist(err) {
		log.Printf("Creating week folder: %s", weekFolder)
		if err := os.Mkdir(weekFolder, 0755); err != nil {
			return "", fmt.Errorf("failed to create week folder: %v", err)
		}
	}
	return weekFolder, nil
}

func compileReportSections(weekFolder string) ([]string, []string, []string, []string, error) {
	files, err := os.ReadDir(weekFolder)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error reading week folder: %v", err)
	}

	var progress, problems, plan, insights []string

	for _, file := range files {
		if !file.IsDir() {
			content, err := data.ReadFromFile(fmt.Sprintf("%s/%s", weekFolder, file.Name()))
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("error reading file %s: %v", file.Name(), err)
			}

			parsedData, err := parseDetailedPayload(content)
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("error parsing payload from file %s: %v", file.Name(), err)
			}

			enhancedProgress, err := enhanceText(parsedData.Progress)
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("error enhancing progress text: %v", err)
			}
			progress = append(progress, formatTeamSection(file.Name(), enhancedProgress))

			// Repeat similar process for problems, plan, and insights
		}
	}

	return progress, problems, plan, insights, nil
}

func enhanceText(text string) (string, error) {
	if cfg.UseOllama {
		return nlp.EnhanceTextWithOllama(text)
	}
	return nlp.EnhanceTextWithOpenAI(text, cfg.OpenAIKey)
}

func parseDetailedPayload(content string) (DetailedPayload, error) {
	var payload DetailedPayload
	err := json.Unmarshal([]byte(content), &payload)
	if err != nil {
		return DetailedPayload{}, fmt.Errorf("error parsing payload: %v", err)
	}
	return payload, nil
}

func formatTeamSection(teamName, text string) string {
	return fmt.Sprintf("## %s\n%s", teamName, text)
}

func createFinalReport(progress, problems, plan, insights string) string {
	return fmt.Sprintf(`
        <html>
        <head><title>Weekly Report</title></head>
        <body>
            <h1>Weekly Report</h1>
            <h2>Progress</h2>
            <p>%s</p>
            <h2>Problems</h2>
            <p>%s</p>
            <h2>Plan</h2>
            <p>%s</p>
            <h2>Insights</h2>
            <p>%s</p>
        </body>
        </html>
    `, progress, problems, plan, insights)
}
