package util

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func Error(c echo.Context, code int, message string) error {
	return c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

func BadRequest(c echo.Context, message string) error {
	return Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, message)
}

func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, message)
}

func InternalError(c echo.Context, message string) error {
	return Error(c, http.StatusInternalServerError, message)
}
