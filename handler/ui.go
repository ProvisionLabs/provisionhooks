package handler

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/jpgomesr/webhook-tester/request"
	"github.com/jpgomesr/webhook-tester/store"
)

func WebhooksPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pagePath := filepath.Join("static", "index.html")
	http.ServeFile(w, r, pagePath)
}

func WebhooksReceived(w http.ResponseWriter, r *http.Request, s *store.WebhooksReceived) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	webhooks := s.ListAll()
	response := make([]request.WebhookRequestResponse, 0, len(webhooks))

	for _, wh := range webhooks {
		response = append(response, request.WebhookRequestResponse{
			Method:     wh.Method,
			Headers:    wh.Headers,
			Body:       string(wh.Body),
			ReceivedAt: wh.ReceivedAt,
		})
	}

	_ = json.NewEncoder(w).Encode(response)
}
