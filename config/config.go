package config

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/hyzx-go/common-b2c/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var config Config

type Config struct {
	System SystemConf `mapstructure:"system" json:"system" yaml:"system"`
	Log    LogConf    `mapstructure:"log" json:"log" yaml:"log"`
	Mysql  MysqlList  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis  RedisConf  `mapstructure:"redis" json:"redis" yaml:"redis"`
	Oss    OssConf    `mapstructure:"oss" json:"oss" yaml:"oss"`
}

type ServerConf struct {
	Env string `mapstructure:"env" json:"env" yaml:"env"`
}

type SystemConf struct {
	ServiceName string `mapstructure:"service_name" json:"serviceName" yaml:"service_name"`
	Version     string `mapstructure:"version" json:"version" yaml:"version"`
	ServePort   string `mapstructure:"serve_port" json:"servePort" yaml:"serve_port"`
	ProjectId   string `mapstructure:"project_id" json:"projectId" yaml:"project_id"`
	HostName    string `mapstructure:"host_name" json:"hostName" yaml:"host_name"`
	Local       string `mapstructure:"local" json:"local" yaml:"local"`
	Lang        string `mapstructure:"lang" json:"lang" yaml:"lang"`
	TimeZone    string `mapstructure:"time_zone" json:"timeZone" yaml:"time_zone"`
}

type LogConf struct {
	LogDir               string `mapstructure:"service_name" json:"logDir" yaml:"log_dir"`
	LogFile              string `mapstructure:"service_name" json:"logFile" yaml:"log_file"`
	StackFilter          string `mapstructure:"stack_filter" json:"stackFilter" yaml:"stack_filter"`
	EnableTerminalOutput bool   `mapstructure:"enable_terminal_output" json:"enableTerminalOutput" yaml:"enable_terminal_output"`
	EnableFileOutput     bool   `mapstructure:"enable_file_output" json:"enableFileOutput" yaml:"enable_file_output"`
	EnableGormOutput     bool   `mapstructure:"enable_gorm.output" json:"enableGormOutput" yaml:"enable_gorm.output"`
}
type Mysql struct {
	InsName     string `mapstructure:"ins_name" json:"insName" yaml:"ins_name"`
	Address     string `mapstructure:"address" json:"address" yaml:"address"`
	DbName      string `mapstructure:"db_name" json:"dbName" yaml:"db_name"`
	Username    string `mapstructure:"username" json:"username" yaml:"username"`
	Password    string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConn int    `mapstructure:"max_idle_conn" json:"maxIdleConn" yaml:"max_idle_conn"`
	MaxOpenConn int    `mapstructure:"max_open_conn" json:"maxOpenConn" yaml:"max_open_conn"`
}
type MysqlList struct {
	List []Mysql `mapstructure:"list" json:"list" yaml:"list"`
}

type RedisList struct {
	List []RedisConf `mapstructure:"list" json:"list" yaml:"list"`
}

type RedisConf struct {
	InsName string `mapstructure:"ins_name" json:"insName" yaml:"ins_name"`

	Address      string `mapstructure:"address" json:"address" yaml:"address"`
	Auth         string `mapstructure:"auth" json:"auth" yaml:"auth"`
	Db           int    `mapstructure:"db" json:"db" yaml:"db"`
	ConnTimeout  int    `mapstructure:"conn_timeout" json:"connTimeout" yaml:"conn_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout" json:"readTimeout" yaml:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" json:"writeTimeout" yaml:"write_timeout"`
	MaxIdle      int    `mapstructure:"max_idle" json:"maxIdle" yaml:"max_idle"`
	MaxActive    int    `mapstructure:"max_active" json:"maxActive" yaml:"max_active"`
	IsWait       bool   `mapstructure:"is_wait" json:"isWait" yaml:"is_wait"`
	IdleTimeout  int    `mapstructure:"idle_timeout" json:"idleTimeout" yaml:"idle_timeout"`
}

type OssConf struct {
	AccessKey    string `mapstructure:"access_key" json:"accessKey" yaml:"access_key"`
	AccessSecret string `mapstructure:"access_secret" json:"accessSecret" yaml:"access_secret"`
	RegionId     string `mapstructure:"region_id" json:"regionId" yaml:"region_id"`
	Endpoint     string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	RoleArn      string `mapstructure:"role_arn" json:"roleArn" yaml:"role_arn"`
	Bucket       string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
}

var dbMap = make(map[string]*gorm.DB)

func GetMysqlIns(insName string) *gorm.DB {

	if ins, ok := dbMap[insName]; ok {
		return ins
	}

	errInfo := "ERR: no such database：" + insName + ", please check conf name"
	panic(errInfo)
}

func (c *MysqlList) InitMysql() {
	for _, mysqlConfig := range c.List {
		dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Address, mysqlConfig.DbName)
		conf := mysql.New(mysql.Config{
			DSN:                       dsn,   // mysql dsn
			DefaultStringSize:         256,   // string type default length
			DisableDatetimePrecision:  true,  // disable datetime precision (Databases earlier than MySQL 5.6 are not supported )
			DontSupportRenameIndex:    true,  // The index is reconstructed after deletion (Databases prior to MySQL 5.7 and MariaDB do not support renamed indexes)
			DontSupportRenameColumn:   true,  // Rename columns with 'change'. Databases prior to MySQL 8 and MariaDB do not support renaming columns
			SkipInitializeWithVersion: false, // This parameter is automatically configured based on the current MySQL version
		})

		opts := &gorm.Config{}
		if config.Log.EnableGormOutput {
			opts = &gorm.Config{Logger: log.NewGormLogger()}
		}

		client, err := gorm.Open(conf, opts)

		if err != nil {
			panic("mysqlErr-" + mysqlConfig.Address + "-err:" + err.Error())
		}

		// Get the common database object sql.DB and use the functionality it provides
		db, err := client.DB()

		if err != nil {
			panic("mysqlErr-" + mysqlConfig.Address + "-err:" + err.Error())
		}

		err = db.Ping()

		if err != nil {
			panic("mysqlErr-" + mysqlConfig.Address + "-err:" + err.Error())
		}

		db.SetMaxIdleConns(mysqlConfig.MaxIdleConn)

		db.SetMaxOpenConns(mysqlConfig.MaxOpenConn)

		db.SetConnMaxLifetime(time.Hour)

		dbMap[mysqlConfig.InsName] = client
	}
}

// 全局连接池
var connPool = map[string]*redis.Pool{}

// 注册 Redis 连接池
func (r *RedisList) InitRedis() error {
	for _, redisConf := range r.List {
		// 初始化连接池
		pool := redisConf.newRedisPool()
		connPool[redisConf.InsName] = pool

		// 验证 Redis 连接是否正常
		if err := redisConf.validateRedisConnection(); err != nil {
			return fmt.Errorf("failed to register Redis instance [%s]: %w", redisConf.InsName, err)
		}
	}
	return nil
}

// 验证 Redis 连接是否正常
func (conf *RedisConf) validateRedisConnection() error {
	conn, err := GetRedisIns(conf.InsName)
	defer conn.Close()
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	// 测试连接是否正常
	if _, err := conn.Do("PING"); err != nil {
		return fmt.Errorf("ping error: %w", err)
	}
	return nil
}

// 获取 Redis 连接
func GetRedisIns(options ...string) (redis.Conn, error) {
	if len(options) == 0 {
		return nil, errors.New("instance name is required")
	}

	instanceName := options[0]
	pool, ok := connPool[instanceName]
	if !ok {
		return nil, fmt.Errorf("redis client not exist for [%s]", instanceName)
	}

	conn := pool.Get()
	if len(options) > 1 {
		db := options[1]
		if _, err := conn.Do("SELECT", db); err != nil {
			conn.Close()
			return nil, fmt.Errorf("redis select db fail: %w", err)
		}
	}

	return conn, nil
}

// 创建 Redis 连接池
func (conf *RedisConf) newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Wait:        conf.IsWait,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp", conf.Address,
				redis.DialConnectTimeout(time.Duration(conf.ConnTimeout)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(conf.ReadTimeout)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(conf.WriteTimeout)*time.Millisecond),
			)
			if err != nil {
				return nil, fmt.Errorf("dial error: %w", err)
			}

			// 认证
			if conf.Auth != "" {
				if _, err := conn.Do("AUTH", conf.Auth); err != nil {
					conn.Close()
					return nil, fmt.Errorf("auth error: %w", err)
				}
			}

			// 选择数据库
			if _, err := conn.Do("SELECT", conf.Db); err != nil {
				conn.Close()
				return nil, fmt.Errorf("select db error: %w", err)
			}

			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err // 不再 panic，而是返回错误
		},
	}
}
