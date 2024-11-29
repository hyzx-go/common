package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	rand_chars  = "0123456789abcedfhigklmnopqrstuvwxyz!@#$%^&*()_+-=,.?/~`"
	charsIdBits = 6
	charsIdMask = 1<<charsIdBits - 1
	charsIdMax  = 63 / charsIdBits
)

func MakeUuid() string {
	uuid, err := uuid.NewV6()
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(uuid.String(), "-", "")
}

func MakePhoneCode() string {
	min, max := 10000, 99999
	code := rand.Intn(max-min+1) + min

	return strconv.Itoa(code)
}

func Md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func InArray(need string, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

func RandString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), charsIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), charsIdMax
		}
		if idx := int(cache & charsIdMask); idx < len(rand_chars) {
			sb.WriteByte(rand_chars[idx])
			i--
		}
		cache >>= charsIdBits
		remain--
	}
	return sb.String()
}

type MapStr map[string]interface{}

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

func GoSafeWithRetry(fn func(), retryCount int) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in goroutine: %v\n", r)
				debug.PrintStack()
				if retryCount > 0 {
					log.Printf("Retrying task... (%d retries left)\n", retryCount)
					GoSafeWithRetry(fn, retryCount-1)
				}
			}
		}()
		fn()
	}()
}
