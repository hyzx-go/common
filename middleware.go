package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

// RequestLogger 是一个记录请求日志的中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取或生成 trace ID
		traceID := GetTraceId(c)
		requestParams := extractRequestParams(c) // 提取请求参数
		startTime := time.Now()

		// 包装 ResponseWriter 以捕获响应内容
		responseWriter := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		// 记录请求开始日志
		logRequestStart(c, traceID, requestParams)

		// 执行请求
		c.Next()

		// 请求结束后记录请求日志，包括响应内容
		logRequestEnd(c, traceID, startTime, responseWriter.body.String())
	}
}

// logRequestStart 记录请求开始时的日志
func logRequestStart(c *gin.Context, traceID string, params map[string]interface{}) {
	logger.WithFields(logrus.Fields{
		"trace-id":  traceID,
		"method":    c.Request.Method,
		"path":      c.Request.URL.Path,
		"client_ip": c.ClientIP(),
		"params":    params,
	}).Info("request received start")
}

// logRequestEnd 记录请求结束时的日志
func logRequestEnd(c *gin.Context, traceID string, startTime time.Time, responseBody string) {
	duration := time.Since(startTime)

	fields := logrus.Fields{
		"trace-id":    traceID,
		"status_code": c.Writer.Status(),
		"latency":     fmt.Sprintf("%.3f", duration.Seconds()),
		"response":    responseBody, // 捕获的响应内容
	}

	if duration > SlowApiThreshold {
		fields["path"] = c.Request.URL.Path
		fields["method"] = c.Request.Method
		logger.WithFields(fields).Warn("slow request")
	}
	logger.WithFields(fields).Info("request completed")
}

// ResponseWriter 包装器，用于捕获响应数据
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)                  // 将响应写入缓存
	return w.ResponseWriter.Write(b) // 将响应写入实际的 ResponseWriter
}

// extractRequestParams 提取 GET 和 POST 请求的参数
func extractRequestParams(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{})

	switch c.Request.Method {
	case "GET":
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}
	case "POST":
		contentType := c.Request.Header.Get("Content-Type")
		if contentType == "application/json" {
			params = parseJSONBody(c)
		} else {
			params = parseFormBody(c)
		}
	}
	return params
}

// parseJSONBody 解析 JSON 格式的 POST 请求体
func parseJSONBody(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{})
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Failed to read request body:", err)
		return params
	}

	if err := json.Unmarshal(bodyBytes, &params); err != nil {
		logger.Error("Failed to parse JSON data:", err)
	}

	// 重置请求体以便后续处理中使用
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return params
}

// parseFormBody 解析表单格式的 POST 请求体
func parseFormBody(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{})
	if err := c.Request.ParseForm(); err != nil {
		logger.Error("Failed to parse form data:", err)
		return params
	}

	for key, values := range c.Request.Form {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	return params
}

// GinRecovery 是一个用于捕获 panic 并记录日志的中间件
func GinRecovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(os.Stderr, func(c *gin.Context, recovered interface{}) {
		traceID, _ := c.Get("trace-id")
		logger.WithFields(logrus.Fields{
			"trace_id": traceID,
			"error":    recovered,
			"path":     c.Request.URL.Path,
			"method":   c.Request.Method,
		}).Warn("panic recovered")

		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
	})
}
