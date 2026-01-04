package database

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// ZapLogger 自定义 GORM Logger，使用 Zap 输出日志
type ZapLogger struct {
	logger *zap.Logger
	level  logger.LogLevel
}

// NewZapLogger 创建 Zap 日志适配器
func NewZapLogger(zapLogger *zap.Logger, level logger.LogLevel) *ZapLogger {
	return &ZapLogger{
		logger: zapLogger,
		level:  level,
	}
}

// LogMode 设置日志级别
func (l *ZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

// Info 记录信息日志
func (l *ZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		l.logger.Sugar().Infof(msg, data...)
	}
}

// Warn 记录警告日志
func (l *ZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Warn {
		l.logger.Sugar().Warnf(msg, data...)
	}
}

// Error 记录错误日志
func (l *ZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Error {
		l.logger.Sugar().Errorf(msg, data...)
	}
}

// Trace 记录 SQL 查询日志
func (l *ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.level >= logger.Error:
		l.logger.Error("SQL error",
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
			zap.Error(err),
		)
	case l.level == logger.Warn && elapsed > 200*time.Millisecond:
		l.logger.Warn("Slow SQL",
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.level >= logger.Info:
		l.logger.Info("SQL query",
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}
