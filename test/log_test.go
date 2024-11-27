package test

import (
	"context"
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
	})

	log := log2.GetLogger()

	// Example usage
	log.WithField("module", "main").Info("Application started")
	log.Warn("This is a warning")
	log.WithError(errors.New("all err")).Error("An error occurred")

	log2.Ctx(nil).Info("main err:", "lalalaallalalalla")
	time.Sleep(10)
}

func TestCtx(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, log2.TraceId, "9999999999888888888")
	logW := log2.Ctx(ctx)
	logW.Info("1111", 222)
}
