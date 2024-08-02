package main

import (
	"net/http"

	"github.com/regentmarkets/ContentAI/handlers"
)

func main() {
	handlers.InitHandlers()

	http.HandleFunc("/webhook", handlers.WebhookHandler)
	http.HandleFunc("/fetch-updates", handlers.FetchUpdatesHandler)
	http.HandleFunc("/generate-report", handlers.ReportGenerationHandler)
	http.ListenAndServe(":8080", nil)
}
