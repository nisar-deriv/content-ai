package handlers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/regentmarkets/ContentAI/data"
)

type ParsedData struct {
	Progress string
	Problems string
	Plan     string
	Insights string
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

			parsedData := parseSlackPayload(content)
			progress = append(progress, formatTeamSection(file.Name(), parsedData.Progress))
			problems = append(problems, formatTeamSection(file.Name(), parsedData.Problems))
			plan = append(plan, formatTeamSection(file.Name(), parsedData.Plan))
			insights = append(insights, formatTeamSection(file.Name(), parsedData.Insights))
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

func createFinalReport(progress, problems, plan, insights string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
        }
        .grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
        }
        .section {
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 10px;
            box-shadow: 0 0 10px rgba(0,0,0,0.1);
        }
        .section h2 {
            font-size: 18px;
            border-bottom: 1px solid #ddd;
            padding-bottom: 10px;
        }
        .team-section {
            margin: 10px 0;
        }
        .team-section img {
            vertical-align: middle;
        }
        .team-section h3 {
            display: inline;
            margin-left: 10px;
            font-size: 16px;
        }
    </style>
</head>
<body>
    <div class="grid cards">
        <div class="section">
            <h2>🚀 Progress</h2>
            %s
        </div>
        <div class="section">
            <h2>⚠️ Problems</h2>
            %s
        </div>
        <div class="section">
            <h2>📋 Plan</h2>
            %s
        </div>
        <div class="section">
            <h2>💡 Insights</h2>
            %s
        </div>
    </div>
</body>
</html>
`, progress, problems, plan, insights)
}

func formatTeamSection(teamName, content string) string {
	teamName = strings.TrimSuffix(teamName, ".txt")
	icon := getTeamIcon(teamName)
	return fmt.Sprintf(`
<div class="team-section">
    <img src="%s" alt="%s" width="24" height="24">
    <h3>Team %s</h3>
    <p>%s</p>
</div>
`, icon, teamName, teamName, formatList(content))
}

func getTeamIcon(teamName string) string {
	switch teamName {
	case "WinOps":
		return "https://path/to/assets/images/winops.png"
	case "PE Production", "PE Development":
		return "https://path/to/assets/images/aws.png"
	case "DBA":
		return "https://path/to/assets/images/dba.png"
	case "Kubernetes Core":
		return "https://path/to/assets/images/kubernetes.png"
	default:
		return "https://path/to/assets/images/default.png"
	}
}

func formatList(content string) string {
	items := strings.Split(content, "\n")
	var formattedItems []string
	for _, item := range items {
		formattedItems = append(formattedItems, fmt.Sprintf("<li>%s</li>", item))
	}
	return fmt.Sprintf("<ul>%s</ul>", strings.Join(formattedItems, "\n"))
}

func parseSlackPayload(text string) ParsedData {
	lines := strings.Split(text, "\n")
	var progress, problems, plan, insights []string
	section := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Cloud Dev Platform Kubernetes progress") {
			section = "progress"
		} else if strings.HasPrefix(line, "Problems") {
			section = "problems"
		} else if strings.HasPrefix(line, "Plan") {
			section = "plan"
		} else if strings.HasPrefix(line, "Insights") {
			section = "insights"
		} else if line == "" {
			continue
		} else {
			switch section {
			case "progress":
				progress = append(progress, line)
			case "problems":
				problems = append(problems, line)
			case "plan":
				plan = append(plan, line)
			case "insights":
				insights = append(insights, line)
			}
		}
	}

	return ParsedData{
		Progress: strings.Join(progress, "\n"),
		Problems: strings.Join(problems, "\n"),
		Plan:     strings.Join(plan, "\n"),
		Insights: strings.Join(insights, "\n"),
	}
}
