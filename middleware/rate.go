package middlewares

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

// 限流器管理
type rateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu        sync.Mutex
	limiters  = make(map[string]*rateLimiter)
	rateLimit = rate.Every(50 * time.Millisecond) // 每 100ms 一个请求
	burst     = 5                                 // 允许的突发请求数
)

// 获取限流器
func getLimiter(key string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	// 如果限流器已存在且最近 3 分钟内访问过，直接返回
	if limiter, exists := limiters[key]; exists && time.Since(limiter.lastSeen) <= 3*time.Minute {
		limiter.lastSeen = time.Now()
		return limiter.limiter
	}

	// 否则创建新的限流器
	limiter := rate.NewLimiter(rateLimit, burst)
	limiters[key] = &rateLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter
}

// 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用客户端 IP 作为限流的唯一标识
		clientIP := c.ClientIP()

		limiter := getLimiter(clientIP)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "Please slow down",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
