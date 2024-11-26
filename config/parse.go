package config

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"path"
	"strings"
	"sync"
)

const (
	_defaultConfigFile = "config.yaml"
)

var (
	ErrNotFind     = errors.New("resource not found")
	_once          sync.Once
	_parser        *parser
	_parserManager *ParserManager
)

type Option func(*Options)

type Parser interface {
	GetEnv() string
	//GetErrorMsg(code string) string
	GetSystemConf() (*SystemConf, error)
	GetLogConf() (*LogConf, error)
	//GetHttpClientConf() *HttpClientConf

	GetMysqlDnMap() (map[string]*gorm.DB, error)
	GetRedisDbMap() (map[string]*redis.Pool, error)
	//GetHTTPClient() rpc.Http
	GetParserManager() *ParserManager
}
type ParserManager struct {
	loadType                string
	beforeInitializeConfigs []func() error
	afterInitializeConfigs  []func(p Parser) error
	parserLoader            ParserLoader
}
type ParserLoader interface {
	Load()
	Destroy()
}

type parser struct {
	options    *Options
	serverConf *ServerConf
	mysqlDB    map[string]*gorm.DB
	redisDB    map[string]*redis.Pool

	beanKeys []string
	env      string

	factoryBeans []BeanFactory
	systemConf   *SystemConf
	logConf      *LogConf
	mysqlConf    MysqlList
	redisConf    RedisList
	//httpClientConf *HttpClientConf
}

func (p *parser) GetMysqlDnMap() (map[string]*gorm.DB, error) {
	if p.mysqlDB == nil || len(p.mysqlDB) == 0 {
		return nil, ErrNotFind
	}
	return p.mysqlDB, nil
}

func (p *parser) GetRedisDbMap() (map[string]*redis.Pool, error) {
	if p.redisDB == nil || len(p.redisDB) == 0 {
		return nil, ErrNotFind
	}
	return p.redisDB, nil
}

func (p *parser) GetEnv() string {
	return p.env
}

//func (p *parser) GetErrorMsg(code string) string {
//	//TODO implement me
//	panic("implement me")
//}

func (p *parser) GetSystemConf() (*SystemConf, error) {
	if p.systemConf == nil {
		return nil, ErrNotFind
	}
	return p.systemConf, nil
}

func (p *parser) GetLogConf() (*LogConf, error) {
	if p.logConf == nil {
		return nil, ErrNotFind
	}
	return p.logConf, nil
}

func (p *parser) GetParserManager() *ParserManager {
	return _parserManager
}

func newParser() *parser {
	_once.Do(func() {
		_parser = &parser{}
	})
	return _parser
}
func GetParser() Parser {
	return _parser
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

func (p *ParserManager) Initialize() {
	// Before load configs
	if p.beforeInitializeConfigs != nil {
		for _, beforeInitializeConfig := range p.beforeInitializeConfigs {
			if err := beforeInitializeConfig(); err != nil {
				panic(err)
			}
		}
	}

	// Load configs
	Cyan(fmt.Sprintf("load [%s] config start", p.loadType))
	p.parserLoader.Load()
	Cyan(fmt.Sprintf("load [%s] config succeed", p.loadType))

	// After load configs
	if p.afterInitializeConfigs != nil {
		for _, afterInitializeConfig := range p.afterInitializeConfigs {
			if err := afterInitializeConfig(_parser); err != nil {
				panic(err)
			}
		}
	}
}

func (p *ParserManager) Destroy() {
	p.parserLoader.Destroy()
}

func NewParserManager(opts ...Option) *ParserManager {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	pat := path.Join(options.confFilepath.dir, _defaultConfigFile)
	if options.confFilepath.file != "" {
		pat = path.Join(options.confFilepath.dir, options.confFilepath.file)
	}

	v := viper.New()
	v.SetConfigFile(pat)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := v.ReadInConfig(); err != nil {
		panic(errors.New("cannot find config-server config,please check system settings"))
	}

	serverConfig := &ServerConf{}
	if err := v.UnmarshalKey("server", serverConfig); err != nil {
		panic(errors.New("cannot find config-server config,please check system settings"))
	}

	globalParser := newParser()
	globalParser.serverConf = serverConfig
	globalParser.env = serverConfig.Env
	globalParser.options = options

	var parserLoader ParserLoader
	parserLoader = NewDefaultParserLoader(globalParser)

	_parserManager = &ParserManager{
		parserLoader: parserLoader,
	}

	return _parserManager
}
func (p *ParserManager) AfterInitializeConfigs(fn []func(p Parser) error) *ParserManager {
	p.afterInitializeConfigs = fn
	return p
}

func (p *ParserManager) BeforeInitializeConfigs(fn []func() error) *ParserManager {
	p.beforeInitializeConfigs = fn
	return p
}
