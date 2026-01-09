package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	koanf "github.com/knadh/koanf/v2"
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
	k := koanf.New(".")

	// 1. 设置默认值
	_ = k.Set("server.port", 1323)
	_ = k.Set("server.mode", "debug")
	_ = k.Set("database.log_level", "info")
	_ = k.Set("jwt.expire_hours", 24)
	_ = k.Set("log.level", "info")
	_ = k.Set("log.format", "json")

	// 2. 加载配置文件
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// 3. 加载环境变量（覆盖配置文件）
	// 自动读取所有环境变量，格式：SECTION_KEY -> section.key
	if err := k.Load(env.Provider("", "_", func(s string) string {
		return strings.ToLower(strings.ReplaceAll(s, "_", "."))
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// 4. 构建配置结构体（直接从 Koanf 读取）
	cfg := Config{
		Server: ServerConfig{
			Port: k.Int("server.port"),
			Mode: k.String("server.mode"),
		},
		Database: DatabaseConfig{
			Path:     k.String("database.path"),
			LogLevel: k.String("database.log_level"),
		},
		JWT: JWTConfig{
			Secret:      k.String("jwt.secret"),
			ExpireHours: k.Int("jwt.expire_hours"),
		},
		Log: LogConfig{
			Level:  k.String("log.level"),
			Path:   k.String("log.path"),
			Format: k.String("log.format"),
		},
	}

	return &cfg, nil
}
