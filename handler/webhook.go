package handler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jpgomesr/webhook-tester/request"
	"github.com/jpgomesr/webhook-tester/store"
)

func WebhookReceiver(w http.ResponseWriter, r *http.Request, s *store.WebhooksReceived) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("failed to read body: %v", err)
		return
	}
	headers := r.Header
	method := r.Method

	go func() {
		s.Add(request.WebhookRequest{
			Method:     method,
			Headers:    headers,
			Body:       body,
			ReceivedAt: time.Now(),
		})
	}()

	w.WriteHeader(http.StatusOK)
}
