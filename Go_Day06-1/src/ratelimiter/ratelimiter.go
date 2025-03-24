package ratelimiter

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	Counters    map[string]int
	Mu          sync.Mutex
	Window      time.Duration
	MaxRequests int
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rl.Mu.Lock()
		rl.Counters[ip]++
		rl.Mu.Unlock()

		if rl.Counters[ip] > rl.MaxRequests {
			http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			return
		}

		go func() {
			time.Sleep(rl.Window)
			rl.Mu.Lock()
			delete(rl.Counters, ip)
			rl.Mu.Unlock()
		}()

		next.ServeHTTP(w, r)
	})
}
