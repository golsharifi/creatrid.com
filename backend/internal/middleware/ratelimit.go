package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type visitor struct {
	tokens    float64
	lastSeen  time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     float64 // tokens per second
	burst    int     // max tokens
}

func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     requestsPerSecond,
		burst:    burst,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) allow(ip string) (bool, int, int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	now := time.Now()

	if !exists {
		remaining := rl.burst - 1
		rl.visitors[ip] = &visitor{tokens: float64(remaining), lastSeen: now}
		return true, remaining, rl.burst
	}

	elapsed := now.Sub(v.lastSeen).Seconds()
	v.tokens += elapsed * rl.rate
	if v.tokens > float64(rl.burst) {
		v.tokens = float64(rl.burst)
	}
	v.lastSeen = now

	if v.tokens < 1 {
		return false, 0, rl.burst
	}

	v.tokens--
	return true, int(v.tokens), rl.burst
}

func RateLimit(requestsPerSecond float64, burst int) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(requestsPerSecond, burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Real-Ip"); forwarded != "" {
				ip = forwarded
			}

			allowed, remaining, limit := limiter.allow(ip)
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", "1")

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitUser works like RateLimit but uses the authenticated user's ID as the
// rate-limit key when available, falling back to the client IP address for
// unauthenticated requests.
func RateLimitUser(requestsPerSecond float64, burst int) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(requestsPerSecond, burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := ""
			if user := UserFromContext(r.Context()); user != nil {
				key = "user:" + user.ID
			} else {
				key = r.RemoteAddr
				if forwarded := r.Header.Get("X-Real-Ip"); forwarded != "" {
					key = forwarded
				}
			}

			allowed, remaining, limit := limiter.allow(key)
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", "1")

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
