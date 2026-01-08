package common

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code" example:"0"`          // 业务状态码
	Message string      `json:"message" example:"success"` // 响应消息
	Data    interface{} `json:"data"`                      // 响应数据
}

// Pagination 分页数据
type Pagination struct {
	Page     int         `json:"page" example:"1"`       // 当前页
	PageSize int         `json:"page_size" example:"10"` // 每页数量
	Total    int64       `json:"total" example:"100"`    // 总数
	Items    interface{} `json:"items"`                  // 数据列表
}

// 常用业务状态码
const (
	CodeSuccess         = 0   // 成功
	CodeBadRequest      = 400 // 请求参数错误
	CodeUnauthorized    = 401 // 未授权
	CodeForbidden       = 403 // 禁止访问
	CodeNotFound        = 404 // 资源不存在
	CodeInternalError   = 500 // 服务器内部错误
	CodeValidationError = 422 // 参数验证错误
)

// 响应消息映射
var codeMessages = map[int]string{
	CodeSuccess:         "success",
	CodeBadRequest:      "请求参数错误",
	CodeUnauthorized:    "未授权，请先登录",
	CodeForbidden:       "禁止访问",
	CodeNotFound:        "资源不存在",
	CodeInternalError:   "服务器内部错误",
	CodeValidationError: "参数验证失败",
}

// Success 成功响应
// @Success      200 {object} Response
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: codeMessages[CodeSuccess],
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Created 创建成功响应
func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Code:    CodeSuccess,
		Message: "创建成功",
		Data:    data,
	})
}

// Error 错误响应
// @Failure      200 {object} Response
func Error(c echo.Context, code int, message string) error {
	if message == "" {
		message = codeMessages[code]
	}
	return c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// Paginate 分页响应
func Paginate(c echo.Context, page, pageSize int, total int64, items interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: codeMessages[CodeSuccess],
		Data: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Items:    items,
		},
	})
}
