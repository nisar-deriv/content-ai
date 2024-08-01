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
	filename, err := GenerateWeeklyReport()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating weekly report: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Weekly report generated successfully: %s", filename)
}

func GenerateWeeklyReport() (string, error) {
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
	weekEnd := weekStart.AddDate(0, 0, 4)
	weekFolder := fmt.Sprintf("Week %s to %s", weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))

	files, err := os.ReadDir(weekFolder)
	if err != nil {
		return "", fmt.Errorf("error reading week folder: %v", err)
	}

	var progress, problems, plan, insights []string

	for _, file := range files {
		if !file.IsDir() {
			content, err := data.ReadFromFile(fmt.Sprintf("%s/%s", weekFolder, file.Name()))
			if err != nil {
				return "", fmt.Errorf("error reading file %s: %v", file.Name(), err)
			}

			parsedData, err := parseDetailedPayload(content)
			if err != nil {
				return "", fmt.Errorf("error parsing payload from file %s: %v", file.Name(), err)
			}

			enhancedProgress, err := enhanceText(parsedData.Progress)
			if err != nil {
				return "", fmt.Errorf("error enhancing progress text: %v", err)
			}
			enhancedProblems, err := enhanceText(parsedData.Problems)
			if err != nil {
				return "", fmt.Errorf("error enhancing problems text: %v", err)
			}
			enhancedPlan, err := enhanceText(parsedData.Plan)
			if err != nil {
				return "", fmt.Errorf("error enhancing plan text: %v", err)
			}
			enhancedInsights, err := enhanceText(parsedData.Insights)
			if err != nil {
				return "", fmt.Errorf("error enhancing insights text: %v", err)
			}

			progress = append(progress, formatTeamSection(file.Name(), enhancedProgress))
			problems = append(problems, formatTeamSection(file.Name(), enhancedProblems))
			plan = append(plan, formatTeamSection(file.Name(), enhancedPlan))
			insights = append(insights, formatTeamSection(file.Name(), enhancedInsights))
		}
	}

	finalReport := createFinalReport(
		strings.Join(progress, "\n"),
		strings.Join(problems, "\n"),
		strings.Join(plan, "\n"),
		strings.Join(insights, "\n"),
	)

	finalReportFilename := fmt.Sprintf("%s/final_report.html", weekFolder)
	err = data.WriteToFile(finalReportFilename, finalReport)
	if err != nil {
		return "", fmt.Errorf("error writing final report: %v", err)
	}

	return finalReportFilename, nil
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
