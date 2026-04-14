package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ipRecord tracks the request count and the window start time for a single IP.
type ipRecord struct {
	count    int
	windowStart time.Time
	mu       sync.Mutex
}

// rateLimiter holds per-IP records.
type rateLimiter struct {
	records sync.Map
	limit   int
	window  time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{limit: limit, window: window}
	// Periodically clean up stale IP records to prevent unbounded memory growth.
	go func() {
		ticker := time.NewTicker(window * 10)
		defer ticker.Stop()
		for range ticker.C {
			rl.records.Range(func(key, value any) bool {
				record := value.(*ipRecord)
				record.mu.Lock()
				if time.Since(record.windowStart) > window {
					rl.records.Delete(key)
				}
				record.mu.Unlock()
				return true
			})
		}
	}()
	return rl
}

func (rl *rateLimiter) allow(ip string) bool {
	now := time.Now()
	val, _ := rl.records.LoadOrStore(ip, &ipRecord{windowStart: now})
	record := val.(*ipRecord)

	record.mu.Lock()
	defer record.mu.Unlock()

	if now.Sub(record.windowStart) > rl.window {
		// Reset the window
		record.count = 0
		record.windowStart = now
	}
	record.count++
	return record.count <= rl.limit
}

// AuthRateLimiter returns a Gin middleware that allows at most `limit` requests
// per `window` duration per unique client IP. Intended for auth endpoints only.
func AuthRateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	rl := newRateLimiter(limit, window)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests — please wait before trying again",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
