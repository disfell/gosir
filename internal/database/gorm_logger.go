package database

import (
	"context"
	"fmt"
	"time"

	"gosir/internal/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

// ZapLogger 自定义 GORM Logger，使用 zap 输出日志
type ZapLogger struct {
	LogLevel gormlogger.LogLevel
}

var GormLevelToZap = map[gormlogger.LogLevel]zapcore.Level{
	gormlogger.Silent: zapcore.DPanicLevel,
	gormlogger.Error:  zapcore.ErrorLevel,
	gormlogger.Warn:   zapcore.WarnLevel,
	gormlogger.Info:   zapcore.InfoLevel,
}

// LogMode 设置日志级别
func (l *ZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 记录信息日志
func (l *ZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		logger.Logger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn 记录警告日志
func (l *ZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		logger.Logger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error 记录错误日志
func (l *ZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		logger.Logger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace 记录 SQL 查询日志
func (l *ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error:
		logger.Logger.Error("SQL error",
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
			zap.Error(err),
		)
	case l.LogLevel == gormlogger.Warn && elapsed > 200*time.Millisecond:
		logger.Logger.Warn("Slow SQL",
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.LogLevel >= gormlogger.Info:
		logger.Logger.Debug("SQL query",
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}
