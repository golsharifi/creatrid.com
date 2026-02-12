package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	t.Run("initial burst capacity passes", func(t *testing.T) {
		rl := &RateLimiter{
			visitors: make(map[string]*visitor),
			rate:     10,
			burst:    5,
		}

		for i := 0; i < 5; i++ {
			if !rl.allow("192.168.1.1") {
				t.Errorf("request %d should be allowed within burst capacity", i+1)
			}
		}
	})

	t.Run("exceeding burst returns false", func(t *testing.T) {
		rl := &RateLimiter{
			visitors: make(map[string]*visitor),
			rate:     0.1, // very slow refill
			burst:    3,
		}

		// Use all burst tokens
		for i := 0; i < 3; i++ {
			if !rl.allow("192.168.1.1") {
				t.Fatalf("request %d should be allowed within burst", i+1)
			}
		}

		// Next request should be rejected
		if rl.allow("192.168.1.1") {
			t.Error("request after burst exhaustion should be denied")
		}
	})

	t.Run("tokens refill over time", func(t *testing.T) {
		rl := &RateLimiter{
			visitors: make(map[string]*visitor),
			rate:     100, // 100 tokens per second
			burst:    5,
		}

		// Exhaust all tokens
		for i := 0; i < 5; i++ {
			rl.allow("192.168.1.1")
		}

		// Should be denied now
		if rl.allow("192.168.1.1") {
			t.Fatal("should be denied after exhausting burst")
		}

		// Simulate time passing by directly modifying the visitor's lastSeen
		rl.mu.Lock()
		v := rl.visitors["192.168.1.1"]
		v.lastSeen = v.lastSeen.Add(-100 * time.Millisecond) // 100ms ago -> 100ms * 100/s = 10 tokens refilled
		rl.mu.Unlock()

		// Should be allowed again after refill
		if !rl.allow("192.168.1.1") {
			t.Error("should be allowed after tokens refill")
		}
	})

	t.Run("different IPs have separate limits", func(t *testing.T) {
		rl := &RateLimiter{
			visitors: make(map[string]*visitor),
			rate:     0.1,
			burst:    2,
		}

		// Exhaust IP1
		for i := 0; i < 2; i++ {
			rl.allow("10.0.0.1")
		}
		if rl.allow("10.0.0.1") {
			t.Error("IP1 should be denied after exhaustion")
		}

		// IP2 should still have full burst
		if !rl.allow("10.0.0.2") {
			t.Error("IP2 should be allowed - separate bucket")
		}
		if !rl.allow("10.0.0.2") {
			t.Error("IP2 second request should still be allowed")
		}
	})

	t.Run("tokens do not exceed burst cap after long idle", func(t *testing.T) {
		rl := &RateLimiter{
			visitors: make(map[string]*visitor),
			rate:     1000,
			burst:    5,
		}

		// First request
		rl.allow("192.168.1.1")

		// Simulate long idle
		rl.mu.Lock()
		v := rl.visitors["192.168.1.1"]
		v.lastSeen = v.lastSeen.Add(-time.Hour) // 1 hour ago
		rl.mu.Unlock()

		// Should be allowed (tokens refilled to burst cap)
		if !rl.allow("192.168.1.1") {
			t.Error("should be allowed after long idle")
		}

		// Tokens should be capped at burst (5), so we can make burst-1 more requests
		// (we just consumed 1 above)
		for i := 0; i < 4; i++ {
			if !rl.allow("192.168.1.1") {
				t.Errorf("request %d after refill should be allowed", i+1)
			}
		}

		// Now should be denied (5 tokens used)
		if rl.allow("192.168.1.1") {
			t.Error("should be denied after using all burst tokens")
		}
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	t.Run("allowed requests return 200", func(t *testing.T) {
		handler := RateLimit(100, 10)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
	})

	t.Run("exceeded rate returns 429 with Retry-After header", func(t *testing.T) {
		handler := RateLimit(0.001, 1)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// First request uses the single burst token
		req1 := httptest.NewRequest("GET", "/test", nil)
		req1.RemoteAddr = "10.0.0.1:12345"
		rr1 := httptest.NewRecorder()
		handler.ServeHTTP(rr1, req1)

		if rr1.Code != http.StatusOK {
			t.Fatalf("first request should succeed, got %d", rr1.Code)
		}

		// Second request should be rate limited
		req2 := httptest.NewRequest("GET", "/test", nil)
		req2.RemoteAddr = "10.0.0.1:12345"
		rr2 := httptest.NewRecorder()
		handler.ServeHTTP(rr2, req2)

		if rr2.Code != http.StatusTooManyRequests {
			t.Errorf("expected status 429, got %d", rr2.Code)
		}
		if rr2.Header().Get("Retry-After") != "1" {
			t.Errorf("expected Retry-After: 1, got %s", rr2.Header().Get("Retry-After"))
		}
	})

	t.Run("X-Real-Ip header is used for rate limiting", func(t *testing.T) {
		handler := RateLimit(0.001, 1)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// First request from IP via X-Real-Ip
		req1 := httptest.NewRequest("GET", "/test", nil)
		req1.RemoteAddr = "proxy:8080"
		req1.Header.Set("X-Real-Ip", "1.2.3.4")
		rr1 := httptest.NewRecorder()
		handler.ServeHTTP(rr1, req1)

		if rr1.Code != http.StatusOK {
			t.Fatalf("first request should succeed, got %d", rr1.Code)
		}

		// Second request from same X-Real-Ip should be limited
		req2 := httptest.NewRequest("GET", "/test", nil)
		req2.RemoteAddr = "proxy:8080"
		req2.Header.Set("X-Real-Ip", "1.2.3.4")
		rr2 := httptest.NewRecorder()
		handler.ServeHTTP(rr2, req2)

		if rr2.Code != http.StatusTooManyRequests {
			t.Errorf("expected 429 for same X-Real-Ip, got %d", rr2.Code)
		}

		// Request from different X-Real-Ip should pass
		req3 := httptest.NewRequest("GET", "/test", nil)
		req3.RemoteAddr = "proxy:8080"
		req3.Header.Set("X-Real-Ip", "5.6.7.8")
		rr3 := httptest.NewRecorder()
		handler.ServeHTTP(rr3, req3)

		if rr3.Code != http.StatusOK {
			t.Errorf("different X-Real-Ip should succeed, got %d", rr3.Code)
		}
	})
}
