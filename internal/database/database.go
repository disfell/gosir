package database

import (
	"gosir/internal/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dsn, logLevel string) error {
	var err error

	// 解析日志级别
	level := parseLogLevel(logLevel)

	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: &ZapLogger{
			LogLevel: level,
		},
	})

	if err != nil {
		return err
	}

	logger.Logger.Info("Database connected successfully")
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
