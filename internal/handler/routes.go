package handler

import (
	"myapp/internal/handler/auth"
	"myapp/internal/handler/system"
	"myapp/internal/handler/user"
	"myapp/internal/service"

	"github.com/labstack/echo/v4"
)

// SetupPublicRoutes 设置公开路由（无需鉴权）
func SetupPublicRoutes(e *echo.Echo) {
	// 系统路由
	e.GET("/health", system.HealthCheck)

	// 认证路由
	e.POST("/auth/login", auth.Login)

	// 测试 panic 接口
	e.GET("/panic", func(c echo.Context) error {
		// 故意触发 panic
		var arr []int
		_ = arr[10] // 数组越界
		return nil
	})
}

// SetupRoutes 设置受保护路由（需要鉴权）
func SetupRoutes(e *echo.Group) {
	userService := service.NewUserService()
	userHandler := user.New(userService)

	// 用户路由
	e.GET("/users", userHandler.GetUser)
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.GetUser)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)
}
