package config

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"sync"
)

const (
	_defaultConfigFile = "config.yaml"
	//_defaultI18nFile   = "i18n.yaml"
)

type parser struct {
	options       *Options
	serverConf    *ServerConf
	mysqlMasterDB *gorm.DB
	mysqlSlaveDB  *gorm.DB
	redisClient   *redis.Client

	env        string
	systemConf *SystemConf
	logConf    *LogConf
	mysqlConf  *MySQLConf
	//httpClientConf *HttpClientConf
}

type Options struct {
	mu                sync.Mutex
	watchConfigSwitch bool
	rawVal            map[string]interface{}
	confFilepath      struct {
		dir  string
		file string
	}
}
