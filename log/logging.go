package log

import (
	"context"
	"github.com/hyzx-go/common-b2c/global"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	TraceId = "trace-id"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

// GetLogger returns the singleton logger instance.
func GetLogger() *logrus.Entry {
	if logger == nil {
		InitLogger(Config{DefaultConf: DefaultConfig()}) // Use default config if not initialized
	}
	return logrus.NewEntry(logger)
}

// logWrapper 结构体用于封装日志相关的方法
type logWrapper struct {
	log     *logrus.Entry
	ctx     context.Context
	traceId string
}

// Ctx 创建一个新的 logWrapper 实例
func Ctx(ctx context.Context) *logWrapper {
	if logger == nil {
		InitLogger(Config{DefaultConf: DefaultConfig()})
	}

	if ctx == nil {
		return &logWrapper{log: logrus.NewEntry(logger)}
	}

	traceId := ctx.Value(TraceId)
	if traceId == nil {
		traceId = "unknown"
	}
	return &logWrapper{log: logrus.NewEntry(logger), ctx: ctx, traceId: traceId.(string)}
}

// InitLogger initializes the logger with the provided configuration.
func InitLogger(config Config) {
	defaultConf := DefaultConfig()
	if config.DefaultConf == nil {
		config.DefaultConf = defaultConf
	} else {
		// 仅在当前配置为空时才使用默认值
		if config.DefaultConf.Dir == "" {
			config.DefaultConf.Dir = defaultConf.Dir
		}
		if config.DefaultConf.File == "" {
			config.DefaultConf.File = defaultConf.File
		}

		if config.LogLevel < logrus.DebugLevel {
			config.LogLevel = logrus.DebugLevel
		}

		if config.MaxSize == 0 {
			config.MaxSize = defaultConf.MaxSize
		}
		if config.MaxBackups == 0 {
			config.MaxBackups = defaultConf.MaxBackups
		}
		if config.MaxAge == 0 {
			config.MaxAge = defaultConf.MaxAge
		}
	}

	once.Do(func() {
		// 检查并创建日志目录
		if err := os.MkdirAll(config.Dir, 0755); err != nil {
			logger.Fatal("Failed to create log directory:", err)
		}

		logger = logrus.New()

		// 设置日志级别
		logger.SetLevel(config.LogLevel)

		// 启用 ReportCaller 以显示文件名和行号
		logger.SetReportCaller(config.ReportCaller)

		// 配置滚动日志文件
		jsonLogFile := &lumberjack.Logger{
			Filename:   config.Dir + "/" + config.File,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}

		// 配置文本格式和 JSON 格式的 Hook
		hook := &jSONAndTextFormatterHook{
			// 添加全局字段
			fields:     global.LogPreInfo,
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
