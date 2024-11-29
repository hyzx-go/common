package resp

import "fmt"

// ErrorCode 类型定义
type ErrorCode int

// 通用错误码
const (
	// 通用成功与客户端错误码
	Success      ErrorCode = 0   // 成功
	BadRequest   ErrorCode = 400 // 请求参数错误
	Unauthorized ErrorCode = 401 // 未授权
	Forbidden    ErrorCode = 403 // 禁止访问
	NotFound     ErrorCode = 404 // 资源未找到
	Conflict     ErrorCode = 409 // 资源冲突

	// 服务端错误码
	InternalError      ErrorCode = 500 // 服务器内部错误
	ServiceUnavailable ErrorCode = 503 // 服务不可用
	Timeout            ErrorCode = 504 // 请求超时

	// 自定义业务错误码 (09 作为通用业务模块代码)
	ParamsError         ErrorCode = 901 // 参数校验错误
	DatabaseError       ErrorCode = 902 // 数据库操作错误
	AuthenticationError ErrorCode = 903 // 认证失败
	PermissionError     ErrorCode = 904 // 权限不足
	ResourceExists      ErrorCode = 905 // 资源已存在
	OperationFailed     ErrorCode = 906 // 操作失败

	StandError ErrorCode = 999
)

// 用户模块错误码 (01)
const (
	UserNotFound ErrorCode = 10001 // 用户未找到
)

type ErrorCodeModule uint

const (
	ErrorModuleGeneral ErrorCodeModule = 1 //"General"
	ErrorModuleUser    ErrorCodeModule = 2 //"User"
	ErrorModuleUnknown ErrorCodeModule = 3 //"Unknown"
)

func (e ErrorCodeModule) String() string {
	switch e {
	case ErrorModuleGeneral:
		return "General"
	case ErrorModuleUser:
		return "UserModule"
	case ErrorModuleUnknown:
		return "Unknown"
	}
	return "Unknown"
}

// 解析错误码所属模块
func ParseErrorCode(code ErrorCode) (module ErrorCodeModule, detailCode int) {
	codeStr := fmt.Sprintf("%06d", code) // 保证长度一致
	moduleCode := codeStr[:2]            // 取前两位模块编号
	detailCode = int(code) % 100000      // 剩下部分为详细错误码

	switch moduleCode {
	case "00":
		module = ErrorModuleGeneral
	case "01":
		module = ErrorModuleUser
	default:
		module = ErrorModuleUnknown
	}
	return
}

// 错误码多语言映射
var langMap = map[string]map[ErrorCode]string{
	"en": {
		Success:             "Success",
		BadRequest:          "Invalid request parameters",
		Unauthorized:        "Unauthorized",
		Forbidden:           "Forbidden",
		NotFound:            "Resource not found",
		Conflict:            "Resource conflict",
		InternalError:       "Internal server error",
		ServiceUnavailable:  "Service unavailable",
		Timeout:             "Request timeout",
		ParamsError:         "Parameter validation error",
		DatabaseError:       "Database operation error",
		AuthenticationError: "Authentication failed",
		PermissionError:     "Permission denied",
		ResourceExists:      "Resource already exists",
		OperationFailed:     "Operation failed",
	},
	"zh": {
		Success:             "成功",
		BadRequest:          "请求参数错误",
		Unauthorized:        "未授权",
		Forbidden:           "禁止访问",
		NotFound:            "资源未找到",
		Conflict:            "资源冲突",
		InternalError:       "服务器内部错误",
		ServiceUnavailable:  "服务不可用",
		Timeout:             "请求超时",
		ParamsError:         "参数校验错误",
		DatabaseError:       "数据库操作错误",
		AuthenticationError: "认证失败",
		PermissionError:     "权限不足",
		ResourceExists:      "资源已存在",
		OperationFailed:     "操作失败",
	},
}

// 获取错误信息
func GetErrorMessage(code ErrorCode, lang string) string {
	if messages, exists := langMap[lang]; exists {
		if msg, found := messages[code]; found {
			return msg
		}
	}
	return "Unknown error"
}
