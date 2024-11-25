package initialize

import (
	core "github.com/hyzx-go/common-b2c"
)

func BuildAppStarter() *core.Starter {
	return core.NewsStartService()
}
