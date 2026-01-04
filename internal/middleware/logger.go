package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gosir/internal/logger"
)

// ZapLoggerMiddleware Echo 请求日志中间件
func ZapLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			// 处理请求
			err := next(c)

			// 记录日志
			latency := time.Since(start)
			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Int("status", res.Status),
				zap.Duration("latency", latency),
				zap.String("ip", c.RealIP()),
				zap.String("user_agent", req.UserAgent()),
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
				logger.Error("request completed with error", fields...)
			} else {
				logger.Info("request completed", fields...)
			}

			return err
		}
	}
}
