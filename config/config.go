package config

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
	LogLevel    string `mapstructure:"logLevel" json:"logLevel" yaml:"logLevel"`
	LogSoftLink string `mapstructure:"logSoftLink" json:"logSoftLink" yaml:"logSoftLink"`
	LogDir      string `mapstructure:"logDir" json:"logDir" yaml:"logDir"`
	MaxAge      int    `mapstructure:"maxAge" json:"maxAge" yaml:"maxAge"`
}
type MySQLConf struct {
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Addr         string `mapstructure:"addr"`
	InsName      string `mapstructure:"ins_name"`
	Dbname       string `mapstructure:"dbName" json:"dbName" yaml:"dbName"`
	MaxIdleConns int    `mapstructure:"maxIdleConns" json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"maxOpenConns" json:"maxOpenConns" yaml:"maxOpenConns"`
}

type RedisConfig struct {
	InsName      string `mapstructure:"ins_name"`
	Db           int    `mapstructure:"db"`
	Addr         string `mapstructure:"addr"`
	Auth         string `mapstructure:"auth"`
	MaxIdle      int    `mapstructure:"max_idle"`
	MaxActive    int    `mapstructure:"max_active"`
	MaxWait      bool   `mapstructure:"max_wait"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
	ConnTimeout  int    `mapstructure:"conn_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type AliyunOssConfig struct {
	AccessKey    string `mapstructure:"access_key"`
	AccessSecret string `mapstructure:"access_secret"`
	RegionId     string `mapstructure:"region_id"`
	Endpoint     string `mapstructure:"endpoint"`
	RoleArn      string `mapstructure:"rolearn"`
	Bucket       string `mapstructure:"bucket"`
}
