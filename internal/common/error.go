package common

import (
	"fmt"
)

// AppError 自定义错误类型
type AppError struct {
	Code    int    // 业务错误码
	Message string // 错误消息
	Err     error  // 原始错误
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 解包错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建应用错误
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Wrap 包装错误
func Wrap(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
