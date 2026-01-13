package config

import (
	"fmt"
	"strings"

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
func Load(path string) (Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to load config file: %w", err)
	}
	v.SetEnvPrefix("GOSIR")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// PrintConfig 打印配置信息（JWT Secret 已脱敏）
func (c Config) PrintConfig() {
	maskSecret := func(secret string) string {
		if len(secret) <= 8 {
			return strings.Repeat("*", len(secret))
		}
		return secret[:4] + strings.Repeat("*", len(secret)-8) + secret[len(secret)-4:]
	}

	fmt.Println("========================================")
	fmt.Println("         Current Configuration")
	fmt.Println("========================================")
	fmt.Printf("Server:\n")
	fmt.Printf("  Port: %d\n", c.Server.Port)
	fmt.Printf("  Mode: %s\n", c.Server.Mode)
	fmt.Println()
	fmt.Printf("Database:\n")
	fmt.Printf("  Path: %s\n", c.Database.Path)
	fmt.Printf("  LogLevel: %s\n", c.Database.LogLevel)
	fmt.Println()
	fmt.Printf("JWT:\n")
	fmt.Printf("  Secret: %s\n", maskSecret(c.JWT.Secret))
	fmt.Printf("  ExpireHours: %d\n", c.JWT.ExpireHours)
	fmt.Println()
	fmt.Printf("Log:\n")
	fmt.Printf("  Level: %s\n", c.Log.Level)
	fmt.Printf("  Path: %s\n", c.Log.Path)
	fmt.Printf("  Format: %s\n", c.Log.Format)
	fmt.Println("========================================")
}
