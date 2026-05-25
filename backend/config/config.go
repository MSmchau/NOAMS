package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Netmiko  NetmikoConfig
	Oxidized OxidizedConfig
	Monitor  MonitorConfig
	Log      LogConfig
}

type MonitorConfig struct {
	PingInterval int    // 检测间隔（秒）
	PingTimeout  int    // 单次超时（秒）
	PingRetry    int    // 重试次数
	PingMethod   string // 检测方式：tcp/icmp
	Concurrency  int    // 最大并发检测数
}

type ServerConfig struct {
	Port         string
	RunMode      string
	CorsOrigins  []string
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	DSN      string // SQLite: file path; MySQL: auto-built from other fields
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	ExpireHour int
}

type NetmikoConfig struct {
	WorkerURL string
}

type OxidizedConfig struct {
	APIURL string
}

type LogConfig struct {
	Level     string
	Filename  string
	MaxSize   int
	MaxBackup int
	MaxAge    int
}

var AppConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("NOAMS")

	viper.BindEnv("server.port", "NOAMS_SERVER_PORT")
	viper.BindEnv("database.driver", "NOAMS_DB_DRIVER")
	viper.BindEnv("database.host", "NOAMS_DB_HOST")
	viper.BindEnv("database.port", "NOAMS_DB_PORT")
	viper.BindEnv("database.user", "NOAMS_DB_USER")
	viper.BindEnv("database.password", "NOAMS_DB_PASSWORD")
	viper.BindEnv("database.dbname", "NOAMS_DB_NAME")
	viper.BindEnv("database.dsn", "NOAMS_DB_DSN")
	viper.BindEnv("redis.addr", "NOAMS_REDIS_ADDR")
	viper.BindEnv("redis.password", "NOAMS_REDIS_PASSWORD")
	viper.BindEnv("jwt.secret", "NOAMS_JWT_SECRET")
	viper.BindEnv("netmiko.worker_url", "NOAMS_NETMIKO_WORKER_URL")
	viper.BindEnv("oxidized.api_url", "NOAMS_OXIDIZED_API_URL")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.run_mode", "debug")
	viper.SetDefault("server.cors_origins", []string{"*"})
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.dsn", "./data/noams.db")
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.user", "netops")
	viper.SetDefault("database.dbname", "netops")
	viper.SetDefault("redis.addr", "127.0.0.1:6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("jwt.expire_hour", 24)
	viper.SetDefault("monitor.ping_interval", 300)
	viper.SetDefault("monitor.ping_timeout", 3)
	viper.SetDefault("monitor.ping_retry", 1)
	viper.SetDefault("monitor.ping_method", "tcp")
	viper.SetDefault("monitor.concurrency", 50)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.filename", "logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backup", 30)
	viper.SetDefault("log.max_age", 30)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			zap.L().Warn("config file not found, using env and defaults")
		}
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port:         viper.GetString("server.port"),
			RunMode:      viper.GetString("server.run_mode"),
			CorsOrigins:  viper.GetStringSlice("server.cors_origins"),
		},
		Database: DatabaseConfig{
			Driver:   viper.GetString("database.driver"),
			Host:     viper.GetString("database.host"),
			Port:     viper.GetString("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.dbname"),
			DSN:      viper.GetString("database.dsn"),
		},
		Redis: RedisConfig{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("jwt.secret"),
			ExpireHour: viper.GetInt("jwt.expire_hour"),
		},
		Netmiko: NetmikoConfig{
			WorkerURL: viper.GetString("netmiko.worker_url"),
		},
		Oxidized: OxidizedConfig{
			APIURL: viper.GetString("oxidized.api_url"),
		},
		Monitor: MonitorConfig{
			PingInterval: viper.GetInt("monitor.ping_interval"),
			PingTimeout:  viper.GetInt("monitor.ping_timeout"),
			PingRetry:    viper.GetInt("monitor.ping_retry"),
			PingMethod:   viper.GetString("monitor.ping_method"),
			Concurrency:  viper.GetInt("monitor.concurrency"),
		},
		Log: LogConfig{
			Level:     viper.GetString("log.level"),
			Filename:  viper.GetString("log.filename"),
			MaxSize:   viper.GetInt("log.max_size"),
			MaxBackup: viper.GetInt("log.max_backup"),
			MaxAge:    viper.GetInt("log.max_age"),
		},
	}
}

func (c *Config) DSN() string {
	if c.Database.Driver == "sqlite" {
		return c.Database.DSN
	}
	return c.Database.User + ":" + c.Database.Password +
		"@tcp(" + c.Database.Host + ":" + c.Database.Port + ")/" +
		c.Database.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func (c *Config) DBDriver() string {
	if c.Database.Driver == "sqlite" {
		return "sqlite"
	}
	return "mysql"
}
