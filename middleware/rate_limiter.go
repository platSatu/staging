package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// limiterStore menyimpan limiter per IP
var limiterStore = struct {
	sync.RWMutex
	clients map[string]*rate.Limiter
}{clients: make(map[string]*rate.Limiter)}

// getLimiterByIP mengembalikan limiter untuk IP tertentu
func getLimiterByIP(ip string, rps int, burst int) *rate.Limiter {
	limiterStore.RLock()
	limiter, exists := limiterStore.clients[ip]
	limiterStore.RUnlock()

	if !exists {
		limiter = rate.NewLimiter(rate.Limit(rps), burst)
		limiterStore.Lock()
		limiterStore.clients[ip] = limiter
		limiterStore.Unlock()
	}

	return limiter
}

// RateLimiter membuat middleware rate limiter per IP
// rps = request per second, burst = maksimal request sekaligus
func RateLimiter(rps int, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := getLimiterByIP(clientIP, rps, burst)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}

		c.Next()
	}
}
