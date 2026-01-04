package database

import (
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"gosir/internal/logger"
)

var DB *gorm.DB

func InitDB(dsn, logLevel string, zapLogger *zap.Logger) error {
	var err error

	// 解析日志级别
	level := parseLogLevel(logLevel)

	// 创建 Zap 日志适配器
	zapLoggerAdapter := NewZapLogger(zapLogger, level)

	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: zapLoggerAdapter,
	})

	if err != nil {
		return err
	}

	logger.Info("Database connected successfully")
	return nil
}

// parseLogLevel 解析日志级别字符串
func parseLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	default:
		return gormlogger.Info
	}
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
