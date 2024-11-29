package resp

import "github.com/gin-gonic/gin"

type Response struct {
	Code       ErrorCode     `json:"code"`
	Module     ErrorCodeType `json:"module"`
	DetailCode int           `json:"detail_code"`
	Message    string        `json:"message"`
	Data       interface{}   `json:"data,omitempty"`
}

// 成功响应
func OkResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, Response{
		Code:       Success,
		Module:     "General",
		DetailCode: 0,
		Message:    GetErrorMessage(Success, ctx.GetHeader("Accept-Language")),
		Data:       data,
	})
}

func ErrorResp(ctx *gin.Context, code ErrorCode, err error) {
	langRes := ctx.GetHeader("Accept-Language")
	if langRes == "" {
		langRes = En.String()
	}

	module, detailCode := ParseErrorCode(code)

	ctx.JSON(200, Response{
		Code:       code,
		Module:     module,
		DetailCode: detailCode,
		Message:    GetErrorMessage(code, langRes),
		Data:       err.Error(),
	})
}
