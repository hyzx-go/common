package logging

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

// GetLogger returns the singleton logger instance.
func GetLogger() *logrus.Entry {
	if logger == nil {
		InitLogger(DefaultConfig()) // Use default config if not initialized
	}
	return logrus.NewEntry(logger)
}

// logWrapper 结构体用于封装日志相关的方法
type logWrapper struct {
	log     *logrus.Logger
	ctx     *gin.Context
	traceId string
	once    sync.Once
}

// Ctx 创建一个新的 logWrapper 实例
func Ctx(ctx ...*gin.Context) *logWrapper {
	// 如果传入了 gin.Context，则使用第一个，否则使用 nil
	interiorCtx := (func() *gin.Context {
		if len(ctx) > 0 && ctx[0] != nil {
			return ctx[0]
		}
		return &gin.Context{}
	})()

	return &logWrapper{
		log:     logger,
		ctx:     interiorCtx,
		traceId: GetTraceId(interiorCtx),
	}
}

// InitLogger initializes the logger with the provided configuration.
func InitLogger(config Config) {
	once.Do(func() {
		// 检查并创建日志目录
		if err := os.MkdirAll(config.LogFileDir, 0755); err != nil {
			logger.Fatal("Failed to create log directory:", err)
		}

		logger = logrus.New()

		// 设置日志级别
		logger.SetLevel(logrus.DebugLevel)

		// 启用 ReportCaller 以显示文件名和行号
		logger.SetReportCaller(true)

		// 配置滚动日志文件
		jsonLogFile := &lumberjack.Logger{
			Filename:   config.LogFilePath,
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		}

		// 配置文本格式和 JSON 格式的 Hook
		hook := &jSONAndTextFormatterHook{
			// 添加全局字段
			fields: logrus.Fields{
				"app_name":  config.AppName,
				"version":   config.Version,
				"host_name": config.HostName,
			},
			JSONWriter: jsonLogFile,
			JSONFormat: &OrderedJSONFormatter{
				TimestampFormat: time.RFC3339,
			},
		}

		if config.EnableFileOutput {
			// 将 Hook 添加到 logrus
			logger.AddHook(hook)
		}

		if config.EnableTerminalOutput {
			// 将日志的主输出设置为终端
			logger.SetOutput(os.Stdout)
		} else {
			logger.SetOutput(ioutil.Discard)
		}
	})
}

// Info 封装 Info 级别的日志打印
func (lw *logWrapper) Info(keyword string, messages ...interface{}) {

	var message interface{}
	if len(messages) > 0 {
		message = messages[0]
	}

	lw.log.WithFields(logrus.Fields{
		"message":  message,
		"trace-id": lw.traceId, // 假设你在 gin.Context 中存储了 trace_id
	}).Info(keyword)
}

// Warn 封装 Warn 级别的日志打印
func (lw *logWrapper) Warn(keyword string, messages ...interface{}) {
	var message interface{}
	if len(messages) > 0 {
		message = messages[0]
	}

	// 检查 err 是否为 error 类型
	if e, ok := message.(error); ok {
		message = e.Error() // 将 error 信息转换为字符串并赋值给 err
	}

	lw.log.WithFields(logrus.Fields{
		"message":  message,
		"trace-id": lw.traceId,
	}).Warn(keyword)
}

// Error 封装 Error 级别的日志打印
func (lw *logWrapper) Error(keyword string, messages ...interface{}) {
	var message interface{}
	if len(messages) > 0 {
		message = messages[0]
	}

	// 检查 err 是否为 error 类型
	if e, ok := message.(error); ok {
		// 需要报警的错误列表
		if IsWarnError(e) {
			lw.Warn(keyword, e)
			return
		}

		message = e.Error() // 将 error 信息转换为字符串并赋值给 err
	}

	lw.log.WithFields(logrus.Fields{
		"message":  message,
		"trace-id": lw.traceId,
	}).Error(keyword)
}
