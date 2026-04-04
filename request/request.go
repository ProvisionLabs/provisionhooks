package request

import "time"

type WebhookRequest struct {
	Method     string
	Headers    map[string][]string
	Body       []byte
	ReceivedAt time.Time
}

type WebhookRequestResponse struct {
	Method     string              `json:"method"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	ReceivedAt time.Time           `json:"received_at"`
}
