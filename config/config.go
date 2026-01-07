package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Path     string
	LogLevel string // silent, error, warn, info
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type LogConfig struct {
	Level  string
	Path   string
	Format string
}

// Load 加载配置（支持配置文件和环境变量，环境变量优先级更高）
func Load(path string) (*Config, error) {
	v := viper.New()

	// 设置配置文件路径
	v.SetConfigFile(path)

	// 设置默认值
	setDefaults(v)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// 支持环境变量
	v.AutomaticEnv()

	// 绑定特定的环境变量到配置项
	bindEnvVariables(v)

	// 解析配置到结构体
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 1323)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("database.log_level", "info")
	v.SetDefault("jwt.expire_hours", 24)
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
}

// bindEnvVariables 绑定环境变量
func bindEnvVariables(v *viper.Viper) {
	// Server 配置
	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("server.mode", "SERVER_MODE")

	// Database 配置
	v.BindEnv("database.path", "DATABASE_PATH")
	v.BindEnv("database.log_level", "DATABASE_LOG_LEVEL")

	// JWT 配置
	v.BindEnv("jwt.secret", "JWT_SECRET")
	v.BindEnv("jwt.expire_hours", "JWT_EXPIRE_HOURS")

	// Log 配置
	v.BindEnv("log.level", "LOG_LEVEL")
	v.BindEnv("log.path", "LOG_PATH")
	v.BindEnv("log.format", "LOG_FORMAT")
}
