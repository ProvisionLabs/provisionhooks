package main

import (
	"fmt"
	"net/http"

	"github.com/jpgomesr/webhook-tester/handler"
	"github.com/jpgomesr/webhook-tester/store"
)

func main() {
	s := store.NewWebhooksReceived()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/hooks", func(w http.ResponseWriter, r *http.Request) {
		handler.WebhookReceiver(w, r, s)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("failed to start server: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
