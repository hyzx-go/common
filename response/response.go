package response

import (
	"github.com/gin-gonic/gin"
	"github.com/hyzx-go/common-b2c/utils"
	"io"
	"net/http"
)

type Response struct {
	TraceId    string      `json:"trace-id"`
	Code       ErrorCode   `json:"code"`
	Module     string      `json:"module"`
	DetailCode int         `json:"detail_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

type ErrResp struct {
	Code ErrorCode   `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

func Resp(code ErrorCode, data interface{}, msg string, c *gin.Context) {
	if c.IsAborted() {
		return
	}
	module, DetailCode := ParseErrorCode(code)
	// 开始时间
	response := Response{
		Code:       code,
		Module:     module.String(),
		DetailCode: DetailCode,
		Message:    msg,
		Data:       data,
	}
	c.JSON(http.StatusOK, response)
	c.Abort()
}

func Fail(errResp ErrResp, c *gin.Context) {
	if c.IsAborted() {
		return
	}
	module, detailCode := ParseErrorCode(errResp.Code)

	c.JSON(http.StatusOK, Response{
		TraceId:    utils.GetTraceId(c),
		Code:       errResp.Code,
		Module:     module.String(),
		DetailCode: detailCode,
		Message:    GetErrorMessage(Success, Lang(c.GetHeader("Accept-Language"))),
		Data:       errResp.Data,
	})
}

// 成功响应
func Ok(data interface{}, c *gin.Context) {
	module, DetailCode := ParseErrorCode(Success)
	c.JSON(http.StatusOK, Response{
		Code:       Success,
		Module:     module.String(),
		DetailCode: DetailCode,
		Message:    GetErrorMessage(Success, Lang(c.GetHeader("Accept-Language"))),
		Data:       data,
	})
}

func File(filename string, length int, reader io.Reader, c *gin.Context) {
	if c.IsAborted() {
		return
	}
	headers := make(map[string]string)
	headers["content-disposition"] = "attachment; filename=\"" + filename + "\""
	c.DataFromReader(http.StatusOK, int64(length), "application/octet-stream", reader, headers)
	c.Abort()
}

func OkWithMessage(message string, c *gin.Context) {
	Resp(Success, map[string]interface{}{}, message, c)
}
func OkWithData(data interface{}, c *gin.Context) {
	Resp(Success, data, "SUCCESS", c)
}
func OkDetailed(data interface{}, message string, c *gin.Context) {
	Resp(Success, data, message, c)
}
func FailWithMessage(message string, c *gin.Context) {
	Resp(StandError, map[string]interface{}{}, message, c)
}
func FailWithDetailed(code ErrorCode, data interface{}, message string, c *gin.Context) {
	Resp(code, data, message, c)
}
func FailWithCodeMsg(code ErrorCode, message string, c *gin.Context) {
	Resp(code, "", message, c)
}
