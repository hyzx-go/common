package test

import (
	"github.com/gin-gonic/gin"
	"github.com/hyzx-go/common-b2c/initialize"
	"testing"
)

// UserModule 路由注册
func UserModule(r *gin.RouterGroup) {
	r.GET("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "user info"})
	})
}

// ProductModule 路由注册
func ProductModule(r *gin.RouterGroup) {
	r.GET("/product", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "product info"})
	})
}

func TestLoadConfig(t *testing.T) {

	// 模块化路由
	modules := []func(r *gin.RouterGroup){
		UserModule,
		ProductModule,
	}

	initialize.BuildAppStarter(modules).Start()
}
