package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var clientRateLimitCleanup sync.Once

// CleanRateLimitClientMap каждую секунду чистит ограничение на наш список пользователей
func CleanRateLimitClientMap(clientRateLimiter *rpsOptions) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		clientRateLimiter.rateLimitClientMap = sync.Map{}
	}
}

// ClientRateLimit установка лимита RPS на клиента
func ClientRateLimit(clientRateLimiter *rpsOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientRateLimitCleanup.Do(func() {
			go CleanRateLimitClientMap(clientRateLimiter)
		})

		userId, err := extractUserId(ctx)
		if err != nil {
			ctx.Next()
			return
		}

		rateLimiter := clientRateLimiter.getLimiter(userId)
		if !rateLimiter.Allow() {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}

func extractUserId(ctx *gin.Context) (string, error) {
	// TODO с идентификацией userId нужно ставить в сессию и доставать соответственно из сессии.
	return "", nil
}
