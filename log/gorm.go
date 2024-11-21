package log

import (
	"context"
	"github.com/sirupsen/logrus"
	gLog "gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	Logger *logrus.Logger
}

// NewGormLogger 创建一个新的 zap 日志记录器
func NewGormLogger() *GormLogger {
	return &GormLogger{Logger: logger}
}

// LogMode 实现 gorm.Logger 接口，配置日志级别
func (z *GormLogger) LogMode(level gLog.LogLevel) gLog.Interface {
	// LogLevel 设置日志级别
	return &GormLogger{Logger: z.Logger}
}

// Info 实现 gorm.Logger 接口，用于记录普通信息
func (z *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	// 获取 trace ID（例如从 Gin Context 或 ctx 中提取）
	traceID, ok := ctx.Value("trace-id").(string)
	if !ok || traceID == "" {
		traceID = "unknown"
	}

	// Info 方法
	z.Logger.WithFields(logrus.Fields{"trace-id": traceID}).Infof(msg, args...)
}

// Warn 实现 gorm.Logger 接口，用于记录警告信息
func (z *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	traceID, ok := ctx.Value("trace-id").(string)
	if !ok || traceID == "" {
		traceID = "unknown"
	}
	// Warn 方法
	z.Logger.WithFields(logrus.Fields{"trace-id": traceID}).Warnf(msg, args...)
}

// Error 实现 gorm.Logger 接口，用于记录错误信息
func (z *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	traceID, ok := ctx.Value("trace-id").(string)
	if !ok || traceID == "" {
		traceID = "unknown"
	}
	// Error 方法
	z.Logger.WithFields(logrus.Fields{"trace-id": traceID}).Errorf(msg, args...)
}

// Trace 实现 gorm.Logger 接口，用于记录 SQL 执行日志
func (z *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	duration := time.Since(begin)
	sql, rows := fc()
	traceID, ok := ctx.Value("trace-id").(string)
	if !ok || traceID == "" {
		traceID = "unknown"
	}
	// 记录 SQL 执行的日志
	if err != nil {
		// 错误日志
		z.Logger.WithFields(logrus.Fields{"trace-id": traceID}).Errorf("SQL query failed: %v | Duration: %v | Rows: %v | Error: %v", sql, duration, rows, err)
	} else {
		// 正常执行日志
		z.Logger.WithFields(logrus.Fields{"trace-id": traceID}).Infof("SQL executed: %v | Duration: %v | Rows: %v", sql, duration, rows)
	}
}
