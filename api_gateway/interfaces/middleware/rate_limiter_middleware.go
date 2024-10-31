package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Mutex to protect the userLimiters map
var mu sync.Mutex

// Store rate limiters for each user or IP address.
var userLimiters = make(map[string]*rate.Limiter)


func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use the client's IP address as the key.
		userKey := c.ClientIP()

		limiter := getLimiter(userKey)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// getLimiter returns a rate limiter for the given user or IP.
func getLimiter(key string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()


	limiter, exists := userLimiters[key]
	if !exists {
	
		limiter = rate.NewLimiter(1, 5)
		userLimiters[key] = limiter
	}
	return limiter
}
