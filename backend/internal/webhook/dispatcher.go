package webhook

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/creatrid/creatrid/internal/store"
)

// Dispatcher queues webhook events for delivery by the background worker.
type Dispatcher struct {
	store *store.Store
}

// NewDispatcher creates a new webhook dispatcher.
func NewDispatcher(st *store.Store) *Dispatcher {
	return &Dispatcher{store: st}
}

// Dispatch queues a webhook event for all matching endpoints of a user.
// It creates delivery records for the background worker to process.
func (d *Dispatcher) Dispatch(ctx context.Context, userID, eventType string, payload interface{}) {
	endpoints, err := d.store.ListActiveWebhookEndpointsForEvent(ctx, userID, eventType)
	if err != nil || len(endpoints) == 0 {
		return
	}

	payloadJSON, err := json.Marshal(map[string]interface{}{
		"event":     eventType,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"data":      payload,
	})
	if err != nil {
		log.Printf("Webhook dispatch marshal error: %v", err)
		return
	}

	for _, ep := range endpoints {
		_, err := d.store.CreateWebhookDelivery(ctx, ep.ID, eventType, payloadJSON)
		if err != nil {
			log.Printf("Webhook delivery creation error for endpoint %s: %v", ep.ID, err)
		}
	}
}
