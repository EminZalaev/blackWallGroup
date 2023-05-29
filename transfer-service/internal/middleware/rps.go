package middleware

import (
	"sync"

	"github.com/spf13/viper"
	"golang.org/x/time/rate"

	"transfer-service/internal/config"
)

// rpsOptions Настройки RSP
type rpsOptions struct {
	// rateLimitClientMap map[string]*rate.Limiter - хранит в себе jwt токен клиента и счетчик лимита по запросам в 1с
	rateLimitClientMap sync.Map
	// rateLimit лимит по запросам
	rateLimit rate.Limit
	// burst распределение запросов
	burst int
}

func NewClientRateLimiter() *rpsOptions {
	clientRateLimiter := &rpsOptions{
		rateLimitClientMap: sync.Map{},
		rateLimit:          rate.Limit(viper.GetFloat64(config.ClientRateRps)),
		burst:              viper.GetInt(config.ClientRateBurst),
	}

	return clientRateLimiter
}

func (c *rpsOptions) addClientLimiter(ownerId string) *rate.Limiter {
	limiter := rate.NewLimiter(c.rateLimit, c.burst)
	c.rateLimitClientMap.Store(ownerId, limiter)
	return limiter
}

func (c *rpsOptions) getLimiter(ownerId string) *rate.Limiter {
	limiter, exists := c.rateLimitClientMap.Load(ownerId)
	if !exists {
		return c.addClientLimiter(ownerId)
	}

	return limiter.(*rate.Limiter)
}
