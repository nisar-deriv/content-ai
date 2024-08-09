package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/regentmarkets/ContentAI/nlp"
)

func ReportGenerationHandlerAi(w http.ResponseWriter, r *http.Request) {
	err := GenerateWeeklyReportsAi()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating weekly reports: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Weekly reports generated successfully")
}

func GenerateWeeklyReportsAi() error {
	weekFolder := getWeekFolder()
	files, err := os.ReadDir(weekFolder)
	if err != nil {
		return fmt.Errorf("error reading week folder: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := fmt.Sprintf("%s/%s", weekFolder, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", file.Name(), err)
			}

			enhancedContent, err := enhanceFullContent(string(content))
			if err != nil {
				return fmt.Errorf("error enhancing content for file %s: %v", file.Name(), err)
			}

			enhancedFilename := fmt.Sprintf("%s/enhanced_%s", weekFolder, file.Name())
			err = os.WriteFile(enhancedFilename, []byte(enhancedContent), 0644)
			if err != nil {
				return fmt.Errorf("error writing enhanced content to file %s: %v", enhancedFilename, err)
			}

		}
	}

	return nil
}

func enhanceFullContent(content string) (string, error) {
	return nlp.EnhanceTextWithOllama(content)
}
