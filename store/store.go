package store

import (
	"sync"

	"github.com/jpgomesr/webhook-tester/request"
)

type WebhooksReceived struct {
	mu       sync.Mutex
	webhooks []request.WebhookRequest
}

type Clients struct {
	mu       sync.Mutex
	channels map[chan request.WebhookRequest]struct{}
}

func NewWebhooksReceived() *WebhooksReceived {
	return &WebhooksReceived{}
}

func NewClients() *Clients {
	return &Clients{
		channels: make(map[chan request.WebhookRequest]struct{}),
	}
}

func (whs *WebhooksReceived) Add(wHook request.WebhookRequest) {
	whs.mu.Lock()
	defer whs.mu.Unlock()
	whs.webhooks = append(whs.webhooks, wHook)
}

func (whs *WebhooksReceived) ListAll() []request.WebhookRequest {
	whs.mu.Lock()
	defer whs.mu.Unlock()

	webhooks := make([]request.WebhookRequest, len(whs.webhooks))
	copy(webhooks, whs.webhooks)
	return webhooks
}

func (clt *Clients) Add(ch chan request.WebhookRequest) {
	clt.mu.Lock()
	defer clt.mu.Unlock()
	clt.channels[ch] = struct{}{}
}

func (clt *Clients) Remove(ch chan request.WebhookRequest) {
	clt.mu.Lock()
	defer clt.mu.Unlock()
	delete(clt.channels, ch)
	close(ch)
}

func (clt *Clients) Broadcast(req request.WebhookRequest) {
	clt.mu.Lock()
	defer clt.mu.Unlock()

	for ch := range clt.channels {
		select {
		case ch <- req:
		default:
		}
	}
}
