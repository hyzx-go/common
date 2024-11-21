package log

import (
	"errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

const SlowApiThreshold = 5 * time.Second

// 定义需要警告的错误数组
var WarnErrorSlice = []error{
	//gorm
	gorm.ErrInvalidDB,          // 无效的数据库连接
	gorm.ErrDuplicatedKey,      // 唯一键冲突
	gorm.ErrForeignKeyViolated, // 外键冲突
	gorm.ErrInvalidTransaction, // 无效事务
	gorm.ErrMissingWhereClause, // 缺少 WHERE 子句
	//mysql connect
	errors.New("connect: connection refused"), // 数据库连接拒绝
	errors.New("driver: bad connection"),
}

func IsWarnError(err error) bool {
	for _, warnErr := range WarnErrorSlice {
		// 判断 err 是否和 warnErr 完全匹配，或 err 包含 warnErr 的字符串内容
		if errors.Is(err, warnErr) || strings.Contains(err.Error(), warnErr.Error()) {
			return true
		}
	}
	return false
}
