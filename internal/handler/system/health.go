package system

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string    `json:"status" example:"ok"`                      // 状态
	Timestamp time.Time `json:"timestamp" example:"2026-01-08T10:00:00Z"` // 时间戳
}

// HealthCheck 健康检查
// @Summary      健康检查
// @Description  检查服务健康状态
// @Tags         系统
// @Accept       json
// @Produce      json
// @Success      200 {object} HealthResponse
// @Router       /health [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	})
}
