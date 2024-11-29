package middlewares

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

// 定义全局的限流器存储
type ipLimiter struct {
	limiter  *rate.Limiter // 具体的限流器
	lastSeen time.Time     // 最后一次访问时间
}

var (
	limiterStore = make(map[string]*ipLimiter) // 存储每个 IP 的限流器
	storeMutex   sync.Mutex                    // 保护共享数据的互斥锁
	rateLimit    = rate.Limit(5)               // 每秒允许的请求数
	burstLimit   = 10                          // 令牌桶的容量
	expireAfter  = 5 * time.Minute             // 限流器过期时间
)

// 获取或创建 IP 的限流器
func getLimiter(ip string) *rate.Limiter {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	// 如果限流器存在并未过期，则返回
	if limiter, exists := limiterStore[ip]; exists {
		limiter.lastSeen = time.Now()
		return limiter.limiter
	}

	// 创建新的限流器
	limiter := rate.NewLimiter(rateLimit, burstLimit)
	limiterStore[ip] = &ipLimiter{limiter: limiter, lastSeen: time.Now()}

	return limiter
}

// 清理过期的限流器
func cleanupLimiters() {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	for ip, limiter := range limiterStore {
		if time.Since(limiter.lastSeen) > expireAfter {
			delete(limiterStore, ip)
		}
	}
}

// 中间件实现
func RateLimitMiddleware() gin.HandlerFunc {
	go func() {
		for {
			time.Sleep(time.Minute) // 定时清理
			cleanupLimiters()
		}
	}()

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := getLimiter(clientIP)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
