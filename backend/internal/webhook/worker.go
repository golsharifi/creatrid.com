package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/creatrid/creatrid/internal/store"
)

// retryIntervals defines exponential backoff durations for each retry attempt.
var retryIntervals = []time.Duration{
	30 * time.Second,
	2 * time.Minute,
	10 * time.Minute,
	30 * time.Minute,
	2 * time.Hour,
}

// Worker is a background worker that polls for pending webhook deliveries
// and dispatches them to their target endpoints with retries.
type Worker struct {
	store    *store.Store
	client   *http.Client
	interval time.Duration
	stopCh   chan struct{}
}

// NewWorker creates a new webhook delivery worker.
func NewWorker(st *store.Store) *Worker {
	return &Worker{
		store: st,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		interval: 5 * time.Second,
		stopCh:   make(chan struct{}),
	}
}

// Start begins the polling loop. It blocks until the context is canceled or Stop is called.
func (w *Worker) Start(ctx context.Context) {
	log.Println("Webhook delivery worker started")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Webhook delivery worker stopped (context canceled)")
			return
		case <-w.stopCh:
			log.Println("Webhook delivery worker stopped")
			return
		case <-ticker.C:
			w.processPending(ctx)
		}
	}
}

// Stop signals the worker to stop its polling loop.
func (w *Worker) Stop() {
	close(w.stopCh)
}

// processPending fetches up to 10 pending deliveries and delivers each one.
func (w *Worker) processPending(ctx context.Context) {
	deliveries, err := w.store.ListPendingDeliveries(ctx, 10)
	if err != nil {
		log.Printf("Webhook worker: failed to list pending deliveries: %v", err)
		return
	}

	for _, d := range deliveries {
		w.deliver(ctx, d)
	}
}

// deliver attempts to send a single webhook delivery to its endpoint.
func (w *Worker) deliver(ctx context.Context, delivery *store.WebhookDelivery) {
	// Look up the endpoint to get URL and secret
	ep, err := w.store.FindWebhookEndpointByID(ctx, delivery.EndpointID)
	if err != nil || ep == nil {
		log.Printf("Webhook worker: endpoint %s not found for delivery %d, marking dead", delivery.EndpointID, delivery.ID)
		_ = w.store.MarkDeliveryDead(ctx, delivery.ID)
		return
	}

	// Sign the payload with HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(ep.Secret))
	mac.Write(delivery.Payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	// Create the HTTP request
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, ep.URL, bytes.NewReader(delivery.Payload))
	if err != nil {
		log.Printf("Webhook worker: failed to create request for delivery %d: %v", delivery.ID, err)
		w.handleFailure(ctx, delivery, 0, err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", fmt.Sprintf("sha256=%s", signature))
	req.Header.Set("X-Webhook-Event", delivery.EventType)
	req.Header.Set("X-Webhook-ID", fmt.Sprintf("%d", delivery.ID))

	// Execute the request
	resp, err := w.client.Do(req)
	if err != nil {
		w.handleFailure(ctx, delivery, 0, err.Error())
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	bodyStr := string(bodyBytes)

	// Check if the response indicates success (2xx)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		_ = w.store.IncrementDeliveryAttempt(ctx, delivery.ID, "success", resp.StatusCode, bodyStr, nil)
	} else {
		w.handleFailure(ctx, delivery, resp.StatusCode, bodyStr)
	}
}

// handleFailure records a failed delivery attempt and schedules retries or marks as dead.
func (w *Worker) handleFailure(ctx context.Context, delivery *store.WebhookDelivery, responseStatus int, responseBody string) {
	nextAttempt := delivery.Attempts + 1 // current attempt count after this failure
	maxAttempts := delivery.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 5
	}

	if nextAttempt >= maxAttempts {
		// All retries exhausted
		_ = w.store.IncrementDeliveryAttempt(ctx, delivery.ID, "dead", responseStatus, responseBody, nil)
		log.Printf("Webhook worker: delivery %d exhausted all %d attempts, marked dead", delivery.ID, maxAttempts)
		return
	}

	// Calculate next retry with exponential backoff
	retryIdx := nextAttempt - 1
	if retryIdx < 0 {
		retryIdx = 0
	}
	if retryIdx >= len(retryIntervals) {
		retryIdx = len(retryIntervals) - 1
	}
	nextRetry := time.Now().Add(retryIntervals[retryIdx])

	_ = w.store.IncrementDeliveryAttempt(ctx, delivery.ID, "pending", responseStatus, responseBody, &nextRetry)
	log.Printf("Webhook worker: delivery %d attempt %d failed, retrying at %s", delivery.ID, nextAttempt, nextRetry.Format(time.RFC3339))
}
