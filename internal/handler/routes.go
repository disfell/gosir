package handler

import (
	"gosir/internal/handler/auth"
	"gosir/internal/handler/system"
	userhandler "gosir/internal/handler/user"
	"gosir/internal/service/user"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// SetupPublicRoutes 设置公开路由（无需鉴权）
func SetupPublicRoutes(e *echo.Echo) {
	userService := user.NewUserService()
	authHandler := auth.New(userService)

	// 系统路由
	e.GET("/health", system.HealthCheck)

	// Swagger 文档路由
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// 认证路由
	e.POST("/auth/login", authHandler.Login)
}

// SetupRoutes 设置受保护路由（需要鉴权）
func SetupRoutes(e *echo.Group) {
	userService := user.NewUserService()
	authHandler := auth.New(userService)
	userHandler := userhandler.New(userService)

	// 认证路由
	e.POST("/auth/logout", authHandler.Logout)
	e.POST("/auth/refresh", authHandler.RefreshToken)

	// 用户路由
	e.GET("/users", userHandler.ListUsers)
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.GetUser)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)
}
