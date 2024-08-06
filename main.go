package main

import (
	"net/http"

	"github.com/regentmarkets/ContentAI/handlers"
)

func main() {
	handlers.InitHandlers()

	http.HandleFunc("/fetch-updates", handlers.FetchUpdatesHandler)
	//http.HandleFunc("/generate-report", handlers.ReportGenerationHandler)
	http.HandleFunc("/generate-report-ai", handlers.ReportGenerationHandlerAi)
	//http.HandleFunc("/ghpr", handlers.GitHubPullRequestHandler)
	http.ListenAndServe(":8080", nil)
}
