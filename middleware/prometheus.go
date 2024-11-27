package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

/*
//group.GET("/metrics", gin.WrapH(promhttp.Handler()))
*/

// 定义 Prometheus 指标
var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

// 初始化 Prometheus 指标
func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(requestDuration)
}

// PrometheusMiddleware 是用于统计请求的 Gin 中间件
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 请求结束后记录指标
		statusCode := c.Writer.Status()
		duration := time.Since(startTime).Seconds()

		// 使用标准库的 http.StatusText 获取状态文本
		httpRequests.WithLabelValues(c.Request.Method, c.FullPath(), http.StatusText(statusCode)).Inc()
		requestDuration.WithLabelValues(c.Request.Method, c.FullPath(), http.StatusText(statusCode)).Observe(duration)
	}
}

// PrometheusHandler 返回用于处理 /metrics 路由的 handler
func PrometheusHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
