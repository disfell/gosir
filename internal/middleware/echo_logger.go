package middleware

import (
	"gosir/internal/logger"
	"io"
	"os"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

// EchoLogger 实现了 Echo 的 Logger 接口
type EchoLogger struct{}

// Output 实现 Echo.Logger 接口，将日志重定向到 zap
func (l *EchoLogger) Output() io.Writer {
	return os.Stdout
}

// SetOutput 实现 Echo.Logger 接口
func (l *EchoLogger) SetOutput(w io.Writer) {}

// Prefix 实现 Echo.Logger 接口
func (l *EchoLogger) Prefix() string {
	return ""
}

// SetPrefix 实现 Echo.Logger 接口
func (l *EchoLogger) SetPrefix(p string) {}

// Level 实现 Echo.Logger 接口
func (l *EchoLogger) Level() log.Lvl {
	return log.INFO
}

// SetLevel 实现 Echo.Logger 接口
func (l *EchoLogger) SetLevel(v log.Lvl) {}

// SetHeader 实现 Echo.Logger 接口
func (l *EchoLogger) SetHeader(h string) {}

// Print 实现 Echo.Logger 接口
func (l *EchoLogger) Print(i ...interface{}) {
	logger.Logger.Info("Echo log", zap.Any("msg", i))
}

// Printf 实现 Echo.Logger 接口
func (l *EchoLogger) Printf(format string, i ...interface{}) {
	logger.Logger.Info("Echo log", zap.String("msg", format), zap.Any("args", i))
}

// Printj 实现 Echo.Logger 接口
func (l *EchoLogger) Printj(j log.JSON) {
	logger.Logger.Info("Echo log", zap.Any("json", j))
}

// Debug 实现 Echo.Logger 接口
func (l *EchoLogger) Debug(i ...interface{}) {
	logger.Logger.Debug("Echo debug", zap.Any("msg", i))
}

// Debugf 实现 Echo.Logger 接口
func (l *EchoLogger) Debugf(format string, i ...interface{}) {
	logger.Logger.Debug("Echo debug", zap.String("msg", format), zap.Any("args", i))
}

// Debugj 实现 Echo.Logger 接口
func (l *EchoLogger) Debugj(j log.JSON) {
	logger.Logger.Debug("Echo debug", zap.Any("json", j))
}

// Info 实现 Echo.Logger 接口
func (l *EchoLogger) Info(i ...interface{}) {
	logger.Logger.Info("Echo info", zap.Any("msg", i))
}

// Infof 实现 Echo.Logger 接口
func (l *EchoLogger) Infof(format string, i ...interface{}) {
	logger.Logger.Info("Echo info", zap.String("msg", format), zap.Any("args", i))
}

// Infoj 实现 Echo.Logger 接口
func (l *EchoLogger) Infoj(j log.JSON) {
	logger.Logger.Info("Echo info", zap.Any("json", j))
}

// Warn 实现 Echo.Logger 接口
func (l *EchoLogger) Warn(i ...interface{}) {
	logger.Logger.Warn("Echo warn", zap.Any("msg", i))
}

// Warnf 实现 Echo.Logger 接口
func (l *EchoLogger) Warnf(format string, i ...interface{}) {
	logger.Logger.Warn("Echo warn", zap.String("msg", format), zap.Any("args", i))
}

// Warnj 实现 Echo.Logger 接口
func (l *EchoLogger) Warnj(j log.JSON) {
	logger.Logger.Warn("Echo warn", zap.Any("json", j))
}

// Error 实现 Echo.Logger 接口
func (l *EchoLogger) Error(i ...interface{}) {
	logger.Logger.Error("Echo error", zap.Any("msg", i))
}

// Errorf 实现 Echo.Logger 接口
func (l *EchoLogger) Errorf(format string, i ...interface{}) {
	logger.Logger.Error("Echo error", zap.String("msg", format), zap.Any("args", i))
}

// Errorj 实现 Echo.Logger 接口
func (l *EchoLogger) Errorj(j log.JSON) {
	logger.Logger.Error("Echo error", zap.Any("json", j))
}

// Fatal 实现 Echo.Logger 接口
func (l *EchoLogger) Fatal(i ...interface{}) {
	logger.Logger.Fatal("Echo fatal", zap.Any("msg", i))
}

// Fatalf 实现 Echo.Logger 接口
func (l *EchoLogger) Fatalf(format string, i ...interface{}) {
	logger.Logger.Fatal("Echo fatal", zap.String("msg", format), zap.Any("args", i))
}

// Fatalj 实现 Echo.Logger 接口
func (l *EchoLogger) Fatalj(j log.JSON) {
	logger.Logger.Fatal("Echo fatal", zap.Any("json", j))
}

// Panic 实现 Echo.Logger 接口
func (l *EchoLogger) Panic(i ...interface{}) {
	logger.Logger.Panic("Echo panic", zap.Any("msg", i))
}

// Panicf 实现 Echo.Logger 接口
func (l *EchoLogger) Panicf(format string, i ...interface{}) {
	logger.Logger.Panic("Echo panic", zap.String("msg", format), zap.Any("args", i))
}

// Panicj 实现 Echo.Logger 接口
func (l *EchoLogger) Panicj(j log.JSON) {
	logger.Logger.Panic("Echo panic", zap.Any("json", j))
}
