package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jpgomesr/webhook-tester/handler"
	"github.com/jpgomesr/webhook-tester/store"
)

func main() {
	sw := store.NewWebhooksReceived()
	sc := store.NewClients()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/hooks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.WebhookReceiver(w, r, sw, sc)
		case http.MethodGet:
			handler.WebhooksPage(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, err := fmt.Fprintf(w, "Method %s not allowed", r.Method)
			if err != nil {
				return
			}
		}
	})
	mux.HandleFunc("/hooks/data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, err := fmt.Fprintf(w, "Method %s not allowed", r.Method)
			if err != nil {
				return
			}
			return
		}

		handler.WebhooksReceived(w, r, sw)
	})
	mux.HandleFunc("/hooks/events", func(w http.ResponseWriter, r *http.Request) {
		handler.SseHandler(w, r, sc, sw)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Printf("failed to start server: %v", err)
	}

	log.Printf("[server] (%s) - Server started on port 8080\n", time.Now().Format(time.RFC3339))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
