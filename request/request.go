package request

import "time"

type WebhookRequest struct {
	Method     string
	Headers    map[string][]string
	Body       []byte
	ReceivedAt time.Time
}
