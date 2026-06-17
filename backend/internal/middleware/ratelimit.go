package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimit implements a simple fixed-window limiter per client IP, suitable
// for protecting auth endpoints (PRD: max 10 login attempts / IP / 15min).
func RateLimit(max int, window time.Duration) func(http.Handler) http.Handler {
	type bucket struct {
		count     int
		windowEnd time.Time
	}
	var mu sync.Mutex
	buckets := make(map[string]*bucket)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)

			mu.Lock()
			b, ok := buckets[ip]
			now := time.Now()
			if !ok || now.After(b.windowEnd) {
				b = &bucket{count: 0, windowEnd: now.Add(window)}
				buckets[ip] = b
			}
			b.count++
			exceeded := b.count > max
			mu.Unlock()

			if exceeded {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"rate_limited","message":"too many attempts, try again later"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
