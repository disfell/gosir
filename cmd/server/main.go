package main

import (
	"log"
	"strconv"

	"myapp/config"
	"myapp/internal/database"
	"myapp/internal/handler"
	"myapp/internal/middleware"
	"myapp/internal/service"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer database.CloseDB()

	// 自动迁移
	if err := service.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建 Echo 实例
	e := echo.New()

	// 初始化 JWT
	middleware.InitJWT(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// 设置统一错误处理
	e.HTTPErrorHandler = middleware.ErrorHandler()

	// 全局中间件
	e.Use(echomiddleware.RequestLogger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// 公开路由（无需鉴权）
	handler.SetupPublicRoutes(e)

	// 受保护路由组（需要鉴权）
	protected := e.Group("/api", middleware.AuthMiddleware())
	handler.SetupRoutes(protected)

	// 启动服务
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	e.Logger.Fatal(e.Start(addr))
}
