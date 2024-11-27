package config

import (
	"errors"
	"fmt"
	"github.com/hyzx-go/common-b2c/global"
	"github.com/hyzx-go/common-b2c/log"
	"github.com/hyzx-go/common-b2c/utils"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	// DefaultMaxMsgSize maximum message that client can receive (200 MB).
	defaultMaxMsgSize = 1024 * 1024 * 200

	// DefaultMaxSendMsgSize maximum message that client can send (200 MB).
	defaultMaxSendMsgSize = 1024 * 1024 * 200
)

type BeanFactory interface {
	// Initialize Initialize Config
	Initialize(inConfig bool, parser *parser) error

	// Destroy Connect
	Destroy() error
}

const (
	_defaultLogKey        = "log"
	_defaultMysqlKey      = "mysql"
	_defaultRedisKey      = "redis"
	_defaultSystemKey     = "system"
	_defaultOssKey        = "oss"
	_defaultHttpClientKey = "httpClient"
)

func (p *parser) initBeanKeys() {
	p.beanKeys = []string{
		_defaultSystemKey,
		_defaultLogKey,
		_defaultMysqlKey,
		_defaultRedisKey,
		_defaultHttpClientKey,
		_defaultOssKey,
	}
}

func getBeanFactory(key string) BeanFactory {
	switch key {
	case _defaultSystemKey:
		return &SystemConf{}
	case _defaultLogKey:
		return &LogConf{}
	case _defaultMysqlKey:
		return &MysqlList{}
	case _defaultRedisKey:
		return &RedisList{}
	case _defaultOssKey:
		return &OssConf{}
	case _defaultHttpClientKey:
		return &HttpClientConf{}
	default:
		log.GetLogger().Error(fmt.Sprintf("cannot find this key %s's beanFactory", key))
	}
	return nil
}

func (c *SystemConf) Initialize(inConfig bool, p *parser) error {
	if !inConfig {
		panic(errors.New("please check system config"))
	}
	p.systemConf = c

	// set date time zone
	utils.SetSystemDateTimeZone(c.TimeZone)

	hostname, err := os.Hostname()
	if err != nil {
		log.Ctx(nil).Error(fmt.Sprintf("system config initialize get hostname:%s", hostname), err)
		return errors.New(fmt.Sprintf("system config initialize get hostname err:%s", err.Error()))
	}
	p.systemConf.HostName = hostname
	log.Ctx(nil).Info(fmt.Sprintf("system config initialize successful :%v", p.systemConf))

	global.LogPreInfo = logrus.Fields{
		"app_name":  p.systemConf.ServiceName,
		"version":   p.systemConf.Version,
		"host_name": p.systemConf.HostName,
	}

	return nil
}

func (c *SystemConf) Destroy() error {
	return nil
}

func (c *LogConf) Initialize(inConfig bool, p *parser) error {
	if !inConfig {
		panic(errors.New("please check log config"))
	}

	p.logConf = c
	log.InitLogger(log.Config{
		EnableTerminalOutput: p.logConf.EnableTerminalOutput,
		EnableGormOutput:     p.logConf.EnableGormOutput,
	})
	return nil
}

func (c *LogConf) Destroy() error {
	return nil
}

func (c *MysqlList) Initialize(inConfig bool, p *parser) error {
	if !inConfig {
		return nil
	}
	p.mysqlConf = *c

	p.mysqlDB = p.mysqlConf.ConnsMysql()
	return nil
}
func (c *MysqlList) Destroy() error {
	return nil
}

func (r *RedisList) Initialize(inConfig bool, p *parser) error {
	if !inConfig {
		return nil
	}

	p.redisConf = *r
	rdsPools, err := p.redisConf.InitRedis()
	if err != nil {
		return err
	}
	p.redisDB = rdsPools
	return err
}

func (r *RedisList) Destroy() error {
	for _, item := range r.List {
		client, err := GetRedisIns(item.InsName)
		if err != nil {
			return err
		}
		return client.Close()
	}
	return nil
}
