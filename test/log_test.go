package test

import (
	"errors"
	log2 "github.com/hyzx-go/common-b2c/log"
	"testing"
	"time"
)

func TestLogging(t *testing.T) {
	// Initialize the logger
	log2.InitLogger(log2.Config{
		DefaultConf:          log2.DefaultConfig(),
		EnableTerminalOutput: false,
		EnableGormOutput:     false,
		AppName:              "1.1.1",
		Version:              "2.2.2",
		HostName:             "3.3.3",
	})

	log := log2.GetLogger()

	// Example usage
	log.WithField("module", "main").Info("Application started")
	log.Warn("This is a warning")
	log.WithError(errors.New("all err")).Error("An error occurred")

	log2.Ctx().Info("main err:", "lalalaallalalalla")
	time.Sleep(10)
}
