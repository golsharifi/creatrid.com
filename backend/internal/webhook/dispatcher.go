package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/creatrid/creatrid/internal/store"
)

// Dispatch sends webhook events to all active endpoints for the user subscribed to the event.
// It runs asynchronously and does not block the caller.
func Dispatch(st *store.Store, eventType, userID string, payload interface{}) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			log.Printf("webhook dispatch: failed to marshal payload: %v", err)
			return
		}

		endpoints, err := st.ListActiveWebhookEndpointsForEvent(ctx, userID, eventType)
		if err != nil {
			log.Printf("webhook dispatch: failed to list endpoints: %v", err)
			return
		}

		for _, ep := range endpoints {
			deliverToEndpoint(ctx, st, ep, eventType, payloadJSON)
		}
	}()
}

func deliverToEndpoint(ctx context.Context, st *store.Store, ep *store.WebhookEndpoint, eventType string, payload json.RawMessage) {
	deliveryID, err := st.CreateWebhookDelivery(ctx, ep.ID, eventType, payload)
	if err != nil {
		log.Printf("webhook dispatch: failed to create delivery record: %v", err)
		return
	}

	// Sign the payload with HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(ep.Secret))
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, ep.URL, bytes.NewReader(payload))
	if err != nil {
		log.Printf("webhook dispatch: failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Creatrid-Event", eventType)
	req.Header.Set("X-Creatrid-Signature", signature)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errMsg := err.Error()
		_ = st.UpdateWebhookDelivery(ctx, deliveryID, 0, errMsg)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	_ = st.UpdateWebhookDelivery(ctx, deliveryID, resp.StatusCode, string(bodyBytes))
}
