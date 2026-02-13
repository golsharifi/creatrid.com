package handler

import "sync"

// SSEHub manages Server-Sent Event connections. Each authenticated user can
// have multiple concurrent SSE channels (e.g. multiple browser tabs).
type SSEHub struct {
	mu      sync.RWMutex
	clients map[string][]chan []byte // userID -> channels
}

// NewSSEHub creates and returns a new SSEHub instance.
func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients: make(map[string][]chan []byte),
	}
}

// Subscribe registers a new channel for the given user and returns it. The
// caller must call Unsubscribe when done.
func (h *SSEHub) Subscribe(userID string) chan []byte {
	ch := make(chan []byte, 64)
	h.mu.Lock()
	h.clients[userID] = append(h.clients[userID], ch)
	h.mu.Unlock()
	return ch
}

// Unsubscribe removes a specific channel from the given user's channel list
// and closes it.
func (h *SSEHub) Unsubscribe(userID string, ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	channels := h.clients[userID]
	for i, c := range channels {
		if c == ch {
			h.clients[userID] = append(channels[:i], channels[i+1:]...)
			close(ch)
			break
		}
	}
	if len(h.clients[userID]) == 0 {
		delete(h.clients, userID)
	}
}

// Notify sends data to all active channels belonging to the given user.
// The send is non-blocking: if a channel's buffer is full the message is
// dropped for that particular channel.
func (h *SSEHub) Notify(userID string, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, ch := range h.clients[userID] {
		select {
		case ch <- data:
		default:
			// Channel buffer full, skip to avoid blocking
		}
	}
}
