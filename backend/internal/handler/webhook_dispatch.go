package handler

import (
	"context"

	"github.com/creatrid/creatrid/internal/webhook"
)

// webhookDispatcher is the global webhook dispatcher, set from main.go.
var webhookDispatcher *webhook.Dispatcher

// SetWebhookDispatcher sets the global webhook dispatcher for all handlers.
func SetWebhookDispatcher(d *webhook.Dispatcher) {
	webhookDispatcher = d
}

// dispatchWebhook is a convenience function for handlers to queue webhook events.
// It runs the dispatch in a background goroutine so it does not block the request.
func dispatchWebhook(userID, eventType string, payload interface{}) {
	if webhookDispatcher != nil {
		go webhookDispatcher.Dispatch(context.Background(), userID, eventType, payload)
	}
}
