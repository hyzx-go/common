package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

// IsEmpty checks if a string is empty
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// ToUpperFirst converts the first letter of a string to upper case
func ToUpperFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// Reverse returns the reversed version of a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsNumeric checks if the string is numeric
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// StringToInt64 将字符串转换为 int64 类型，并处理错误
func StringToInt64(s string) (int64, error) {
	// 如果字符串为空，返回错误
	if s == "" {
		return 0, errors.New("input string is empty")
	}

	// 使用 strconv.ParseInt 转换字符串为 int64
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.New("failed to convert string to int64: " + err.Error())
	}

	return result, nil
}

// StringsToInt64s 将字符串数组转换为 int64 数组
func StringsToInt64s(strs []string) ([]int64, error) {
	// 初始化一个 int64 切片，用于存放转换结果
	int64s := make([]int64, len(strs))

	// 遍历字符串数组并转换
	for i, s := range strs {
		// 尝试将字符串转换为 int64
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			// 返回具体错误，指明转换失败的位置
			return nil, errors.New("failed to convert element at index " + strconv.Itoa(i) + ": " + err.Error())
		}
		int64s[i] = val
	}

	return int64s, nil
}

// ReplaceLogPrefix replaces all occurrences of oldPrefix with newPrefix in a log message.
func ReplaceLogPrefix(logMessage, oldPrefix, newPrefix string) string {
	return strings.Replace(logMessage, oldPrefix, newPrefix, -1)
}

// RemovePrefixFromURLs 去掉每个 URL 的前缀
func RemovePrefixFromURLs(input, prefix string) string {
	// 按逗号分割字符串
	urls := strings.Split(input, ",")

	// 去掉每个字符串的前缀
	for i, url := range urls {
		if strings.HasPrefix(url, prefix) {
			urls[i] = strings.TrimPrefix(url, prefix)
		}
	}

	// 再将处理后的字符串用逗号连接
	return strings.Join(urls, ",")
}

func ConvertToInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case []byte:
		return *(*string)(unsafe.Pointer(&v))
	}
	return ""
}

func Format(format string, s ...interface{}) string {
	var params []interface{}
	for _, value := range s {
		if value != nil {
			t := reflect.TypeOf(value)
			if t.Kind() == reflect.Ptr || t.Kind() == reflect.Struct || t.Kind() == reflect.Slice {
				if marshal, err := json.Marshal(&value); err == nil {
					params = append(params, ToString(marshal))
					continue
				}
			}
		}
		params = append(params, value)
	}
	return fmt.Sprintf(format, params...)
}
