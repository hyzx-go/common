package logging

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"path/filepath"
	"runtime"
	"strings"
)

// JSONAndTextFormatterHook 自定义 Hook，用于同时输出 JSON 和文本格式
type jSONAndTextFormatterHook struct {
	JSONWriter io.Writer
	JSONFormat *OrderedJSONFormatter
	fields     logrus.Fields
}

func (hook *jSONAndTextFormatterHook) Fire(entry *logrus.Entry) error {
	for k, v := range hook.fields {
		entry.Data[k] = v
	}

	// 生成 JSON 格式日志
	jsonData, err := hook.JSONFormat.Format(entry)
	if err == nil {
		hook.JSONWriter.Write(jsonData)
	}

	return nil
}

func (hook *jSONAndTextFormatterHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// OrderedJSONFormatter 自定义 JSON 格式化器，确保字段顺序
type OrderedJSONFormatter struct {
	TimestampFormat string
}

type OrderedLogEntry struct {
	Time       string      `json:"time"`
	Level      string      `json:"level"` // 日志等级
	TraceID    interface{} `json:"trace-id,omitempty"`
	Keyword    interface{} `json:"keyword,omitempty"`
	Message    interface{} `json:"message,omitempty"`
	Path       interface{} `json:"path,omitempty"`
	Method     interface{} `json:"method,omitempty"`
	Params     interface{} `json:"params,omitempty"`
	StatusCode interface{} `json:"status_code,omitempty"`
	ClientIP   interface{} `json:"client_ip,omitempty"`
	Latency    interface{} `json:"latency,omitempty"`
	Err        interface{} `json:"err,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	// 额外字段可动态追加到 Extra 字段中
	Extra     map[string]interface{} `json:"-"`
	CallStack []CallerInfo           `json:"call_stack,omitempty"` // 完整调用栈信息
}
type CallerInfo struct {
	File     string `json:"file"`
	Function string `json:"function"`
	Line     int    `json:"line"`
}

func (f *OrderedJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 创建有序日志条目并设置时间和消息
	logEntry := OrderedLogEntry{
		Time:    entry.Time.Format(f.TimestampFormat),
		Level:   entry.Level.String(), // 获取日志等级并设置
		Keyword: entry.Message,
		Extra:   make(map[string]interface{}),
	}

	// 仅在 err 与 warn 级别打印堆栈信息
	if entry.Level == logrus.ErrorLevel || entry.Level == logrus.WarnLevel {
		var callStack []CallerInfo
		for i := 4; ; i++ { // 从第4层开始，跳过日志库内部调用
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			// 过滤掉不属于本地项目的调用信息
			if !strings.Contains(file, "vendor") {
				function := runtime.FuncForPC(pc).Name()
				callStack = append(callStack, CallerInfo{
					File:     filepath.Dir(file),
					Function: function,
					Line:     line,
				})
			}
		}
		logEntry.CallStack = callStack
	}

	// 动态设置字段
	for key, value := range entry.Data {
		switch key {
		case "params":
			logEntry.Params = value
		case "client_ip":
			logEntry.ClientIP = value
		case "trace-id":
			logEntry.TraceID = value
		case "path":
			logEntry.Path = value
		case "method":
			logEntry.Method = value
		case "status_code":
			logEntry.StatusCode = value
		case "latency":
			logEntry.Latency = value
		case "message":
			logEntry.Message = value
		default:
			logEntry.Extra[key] = value
		}
	}

	// 将 `Extra` 中的动态字段追加到 JSON 数据
	type alias OrderedLogEntry
	finalEntry := struct {
		alias
		Extra map[string]interface{} `json:"extra,omitempty"`
	}{
		alias: alias(logEntry),
		Extra: logEntry.Extra,
	}

	// 序列化 `struct` 为 JSON，保证字段顺序
	jsonData, err := json.Marshal(finalEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entry: %v", err)
	}

	return append(jsonData, '\n'), nil
}
