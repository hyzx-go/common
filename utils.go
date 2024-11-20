package logging

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

func GetTraceId(ctx ...*gin.Context) string {
	// 如果传入了 gin.Context，则尝试从上下文中获取或设置 trace-id
	if len(ctx) > 0 && ctx[0] != nil {
		if traceId := ctx[0].GetString("trace-id"); traceId != "" {
			return traceId
		}
		// 如果 trace-id 不存在，则生成一个新的 trace-id 并存入上下文
		traceId := strings.ReplaceAll(uuid.New().String(), "-", "")
		ctx[0].Set("trace-id", traceId)
		return traceId
	}

	// 如果没有提供 gin.Context，直接返回一个新的 trace-id
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
