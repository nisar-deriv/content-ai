package main

import (
	"net/http"

	"github.com/regentmarkets/ContentAI/handlers"
)

func main() {
	http.HandleFunc("/webhook", handlers.WebhookHandler)
	http.ListenAndServe(":8080", nil)
}
