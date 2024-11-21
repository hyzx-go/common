package log

import (
	"github.com/sirupsen/logrus"
)

// Config defines the configuration for the logger.
type Config struct {
	DefaultConf
	EnableTerminalOutput bool
	EnableGormOutput     bool
	AppName              string
	Version              string
	HostName             string
}

type DefaultConf struct {
	LogLevel         logrus.Level // Logging level
	LogFileDir       string       // Log file path
	LogFilePath      string       // Log file path
	MaxSize          int          // Maximum size of a log file (in MB)
	MaxBackups       int          // Maximum number of backup files
	MaxAge           int          // Maximum age of a log file (in days)
	Compress         bool         // Whether to compress old log files
	ReportCaller     bool         // Whether to include caller info
	EnableFileOutput bool
}

// DefaultConfig returns a default configuration for the logger.
func DefaultConfig() DefaultConf {
	return DefaultConf{
		LogLevel:         logrus.DebugLevel,
		LogFileDir:       "./logs",
		LogFilePath:      "./logs/app.log",
		MaxSize:          10,
		MaxBackups:       5,
		MaxAge:           30,
		Compress:         true,
		ReportCaller:     true,
		EnableFileOutput: true,
	}
}
