package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

func CustomLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			err := next(c)
			
			req := c.Request()
			res := c.Response()
			
			duration := time.Since(start)
			
			// 自定义日志格式
			c.Logger().Infof(
				"method=%s path=%s status=%d duration=%v",
				req.Method,
				req.URL.Path,
				res.Status,
				duration,
			)
			
			return err
		}
	}
}
