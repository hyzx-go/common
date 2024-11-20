package logging

import (
	"github.com/sirupsen/logrus"
)

// Config defines the configuration for the logger.
type Config struct {
	LogLevel     logrus.Level // Logging level
	LogFileDir   string       // Log file path
	LogFilePath  string       // Log file path
	MaxSize      int          // Maximum size of a log file (in MB)
	MaxBackups   int          // Maximum number of backup files
	MaxAge       int          // Maximum age of a log file (in days)
	Compress     bool         // Whether to compress old log files
	JSONFormat   bool         // Whether to use JSON format
	ReportCaller bool         // Whether to include caller info

	EnableFileOutput     bool
	EnableTerminalOutput bool

	AppName  string
	Version  string
	HostName string
}

// DefaultConfig returns a default configuration for the logger.
func DefaultConfig() Config {
	return Config{
		LogLevel:     logrus.DebugLevel,
		LogFilePath:  "logs/app.log",
		MaxSize:      10,
		MaxBackups:   3,
		MaxAge:       7,
		Compress:     true,
		JSONFormat:   false,
		ReportCaller: true,
	}
}
