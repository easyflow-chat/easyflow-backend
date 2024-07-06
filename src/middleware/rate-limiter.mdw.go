package middleware

import (
	"easyflow-backend/src/api"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"golang.org/x/time/rate"
)

var userLimiterMap = make(map[string]*rate.Limiter)
var userLimiterMapMutex sync.Mutex

func getUserLimiter(clientIPAddress string, seconds int, limit int) *rate.Limiter {
	userLimiterMapMutex.Lock()
	defer userLimiterMapMutex.Unlock()

	limiter, ok := userLimiterMap[clientIPAddress]
	if !ok {
		limiter = rate.NewLimiter(rate.Limit(seconds), limit)
		userLimiterMap[clientIPAddress] = limiter
	}
	return limiter
}

func RateLimiter(seconds int, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIPAddress := c.ClientIP()

		limiter := getUserLimiter(clientIPAddress, seconds, limit)
		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, api.ApiError{
				Code:    http.StatusTooManyRequests,
				Error:   "Rate limit exceeded",
				Details: "Slow down buddy",
			})
			c.Header("Retry-After", "1")
			c.Abort()
		}

	}
}
