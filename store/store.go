package store

import (
	"sync"

	"github.com/jpgomesr/webhook-tester/request"
)

type WebhooksReceived struct {
	mu       sync.Mutex
	webhooks []request.WebhookRequest
}

func NewWebhooksReceived() *WebhooksReceived {
	return &WebhooksReceived{}
}

func (whs *WebhooksReceived) Add(wHook request.WebhookRequest) {
	whs.mu.Lock()
	defer whs.mu.Unlock()
	whs.webhooks = append(whs.webhooks, wHook)
}
