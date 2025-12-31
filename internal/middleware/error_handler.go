package middleware

import (
	"errors"
	"gosir/internal/common"
	"gosir/internal/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		req := c.Request()

		// 检查是否是 AppError
		var appErr *common.AppError
		if errors.As(err, &appErr) {
			logger.Logger.Error("AppError occurred",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("code", appErr.Code),
				zap.String("message", appErr.Message),
				zap.Error(appErr.Err),
			)

			_ = c.JSON(http.StatusOK, common.Response{
				Code:    appErr.Code,
				Message: appErr.Message,
				Data:    nil,
			})
			return
		}

		// 检查是否是 HTTP 错误
		var httpErr *echo.HTTPError
		if errors.As(err, &httpErr) {
			code := common.CodeInternalError
			message := httpErr.Message.(string)

			switch httpErr.Code {
			case http.StatusBadRequest:
				code = common.CodeBadRequest
			case http.StatusUnauthorized:
				code = common.CodeUnauthorized
			case http.StatusForbidden:
				code = common.CodeForbidden
			case http.StatusNotFound:
				code = common.CodeNotFound
			}

			logger.Logger.Error("HTTPError occurred",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("code", code),
				zap.String("message", message),
			)

			_ = c.JSON(http.StatusOK, common.Response{
				Code:    code,
				Message: message,
				Data:    nil,
			})
			return
		}

		// 其他未知错误
		logger.Logger.Error("Unknown error occurred",
			zap.String("method", req.Method),
			zap.String("path", req.URL.Path),
			zap.Error(err),
		)

		_ = c.JSON(http.StatusOK, common.Response{
			Code:    common.CodeInternalError,
			Message: "服务器内部错误",
			Data:    nil,
		})
	}
}
