package middleware

import (
	"errors"
	"myapp/internal/common"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		req := c.Request()

		// 检查是否是 AppError
		var appErr *common.AppError
		if errors.As(err, &appErr) {
			c.Logger().Errorf("[%s %s] code=%d message=%s err=%v",
				req.Method, req.URL.Path, appErr.Code, appErr.Message, appErr.Err)

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

			c.Logger().Errorf("[%s %s] code=%d message=%s",
				req.Method, req.URL.Path, code, message)

			_ = c.JSON(http.StatusOK, common.Response{
				Code:    code,
				Message: message,
				Data:    nil,
			})
			return
		}

		// 其他未知错误
		c.Logger().Errorf("[%s %s] unknown error: %v", req.Method, req.URL.Path, err)

		_ = c.JSON(http.StatusOK, common.Response{
			Code:    common.CodeInternalError,
			Message: "服务器内部错误",
			Data:    nil,
		})
	}
}
