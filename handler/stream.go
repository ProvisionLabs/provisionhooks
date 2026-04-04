package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jpgomesr/webhook-tester/request"
	"github.com/jpgomesr/webhook-tester/store"
)

func SseHandler(w http.ResponseWriter, r *http.Request, sc *store.Clients, sw *store.WebhooksReceived) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := make(chan request.WebhookRequest, 8)

	sc.Add(ch)

	defer sc.Remove(ch)

	webhooks := sw.ListAll()
	snapshotResponse := make([]request.WebhookRequestResponse, 0, len(webhooks))
	for _, wh := range webhooks {
		snapshotResponse = append(snapshotResponse, toWebhookResponse(wh))
	}

	snapshot, _ := json.Marshal(snapshotResponse)

	_, err := fmt.Fprintf(w, "event: snapshot\ndata: %s\n\n", snapshot)
	if err != nil {
		return
	}
	flusher.Flush()

	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return

		case req, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(toWebhookResponse(req))
			_, err := fmt.Fprintf(w, "event: webhook\ndata: %s\n\n", data)
			if err != nil {
				return
			}
			flusher.Flush()

		case <-ticker.C:
			_, err := fmt.Fprintf(w, ": keep-alive\n\n")
			if err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func toWebhookResponse(reqWebhook request.WebhookRequest) request.WebhookRequestResponse {
	return request.WebhookRequestResponse{
		Method:     reqWebhook.Method,
		Headers:    reqWebhook.Headers,
		Body:       string(reqWebhook.Body),
		ReceivedAt: reqWebhook.ReceivedAt,
	}
}
