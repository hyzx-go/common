package logging

import (
	"errors"
	"testing"
	"time"
)

func TestLogging(t *testing.T) {
	// Initialize the logger
	InitLogger(Config{
		DefaultConf:          DefaultConfig(),
		EnableTerminalOutput: false,
		EnableGormOutput:     false,
		AppName:              "1.1.1",
		Version:              "2.2.2",
		HostName:             "3.3.3",
	})

	log := GetLogger()

	// Example usage
	log.WithField("module", "main").Info("Application started")
	log.Warn("This is a warning")
	log.WithError(errors.New("all err")).Error("An error occurred")

	Ctx().Info("main err:", "lalalaallalalalla")
	time.Sleep(10)
}
