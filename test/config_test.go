package test

import (
	"github.com/gin-gonic/gin"
	common "github.com/hyzx-go/common-b2c"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	//var default_Config config.DefaultParserLoader
	//default_Config.Load()
	registerRoutes := func(r *gin.Engine) {
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		r.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "hello world"})
		})
	}

	var start common.Starter
	start.Start(registerRoutes)
}
