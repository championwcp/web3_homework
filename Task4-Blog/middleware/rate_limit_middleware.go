package middleware

import (
	"net/http"
	"sync"
	"time"

	"blog-system/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 限流器
type RateLimiter struct {
	ips map[string]*rate.Limiter
	mu   sync.RWMutex
	r    rate.Limit
	b    int
}

// NewRateLimiter 创建限流器
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// AddIP 添加 IP 到限流器
func (rl *RateLimiter) AddIP(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter := rate.NewLimiter(rl.r, rl.b)
	rl.ips[ip] = limiter

	return limiter
}

// GetLimiter 获取 IP 的限流器
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.ips[ip]
	rl.mu.RUnlock()

	if !exists {
		return rl.AddIP(ip)
	}

	return limiter
}

// RateLimit 限流中间件
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.GetLimiter(ip)

		if !limiter.Allow() {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "请求过于频繁", "请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GlobalRateLimit 全局限流中间件
func GlobalRateLimit(requestsPerSecond int, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
	
	return func(c *gin.Context) {
		if !limiter.Allow() {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "系统繁忙", "请稍后再试")
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// TokenBucketRateLimit 令牌桶限流中间件
func TokenBucketRateLimit(fillInterval time.Duration, capacity int64) gin.HandlerFunc {
	tokenBucket := make(chan struct{}, capacity)
	
	// 初始化令牌桶
	for i := int64(0); i < capacity; i++ {
		tokenBucket <- struct{}{}
	}
	
	// 定时添加令牌
	go func() {
		ticker := time.NewTicker(fillInterval)
		defer ticker.Stop()
		
		for range ticker.C {
			select {
			case tokenBucket <- struct{}{}:
			default:
				// 桶已满，丢弃令牌
			}
		}
	}()
	
	return func(c *gin.Context) {
		select {
		case <-tokenBucket:
			c.Next()
		default:
			utils.ErrorResponse(c, http.StatusTooManyRequests, "请求过于频繁", "请稍后再试")
			c.Abort()
		}
	}
}