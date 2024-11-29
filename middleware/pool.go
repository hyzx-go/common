package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/hyzx-go/common-b2c/pool"
	"time"
)

// PoolMiddleware 使用线程池的 Gin 中间件
func PoolMiddleware(pool *pool.GoroutinePool, retryCount int, retryDelay time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建任务
		task := NewExampleTask(c.Request.URL.Path)

		// 提交任务到线程池
		pool.Submit(task, retryCount, retryDelay)

		// 为后续 Handler 设置上下文标记
		c.Set("processed", true)

		// 继续处理后续 Handler
		c.Next()
	}
}
