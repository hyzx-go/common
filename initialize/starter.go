package initialize

import (
	"github.com/gin-gonic/gin"
	core "github.com/hyzx-go/common-b2c"
)

func BuildAppStarter(routers []func(r *gin.RouterGroup)) *core.Starter {
	return core.NewsStartService(routers)
}
