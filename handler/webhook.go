package handler

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jpgomesr/webhook-tester/request"
	"github.com/jpgomesr/webhook-tester/store"
)

func WebhookReceiver(w http.ResponseWriter, r *http.Request, sw *store.WebhooksReceived, sc *store.Clients) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		return
	}
	headers := r.Header
	method := r.Method

	req := request.WebhookRequest{
		Method:     method,
		Headers:    headers,
		Body:       body,
		ReceivedAt: time.Now(),
	}

	go func() {
		sw.Add(req)
		sc.Broadcast(req)
	}()

	log.Printf("[webhook] (%s) - %s %s %s", time.Now(), method, headers, body)

	w.WriteHeader(http.StatusOK)
}
