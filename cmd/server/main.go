package main

import (
	"os"
	"path/filepath"
	"strconv"

	"gosir/config"
	_ "gosir/docs" // 导入 swagger 文档
	"gosir/internal/common"
	"gosir/internal/database"
	"gosir/internal/handler"
	"gosir/internal/logger"
	"gosir/internal/middleware"
	"gosir/internal/service"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// @title           Gosir API
// @version         1.0
// @description     一个基于 Go 语言开发的 REST API 服务，使用 Echo 框架构建
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:1323
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 请输入 JWT token，格式：Bearer <token>
func main() {
	// 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// 创建日志目录
	logDir := filepath.Dir(cfg.Log.Path)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	// 初始化日志系统
	if err := logger.InitWithConfig(&logger.LogConfig{
		Path:   cfg.Log.Path,
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	}); err != nil {
		panic("Failed to init logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting application...")

	// 初始化数据库
	if err := database.InitDB(cfg.Database.Path, cfg.Database.LogLevel, logger.Log); err != nil {
		logger.Fatal("Failed to connect database",
			zap.Error(err),
		)
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			logger.Error("Failed to close database", zap.Error(err))
		}
	}()

	// 自动迁移
	if err := service.AutoMigrate(); err != nil {
		logger.Fatal("Failed to migrate database",
			zap.Error(err),
		)
	}

	// 初始化管理员账号
	if err := service.InitAdminUser(); err != nil {
		logger.Fatal("Failed to init admin user",
			zap.Error(err),
		)
	}

	// 创建 Echo 实例
	e := echo.New()
	e.HideBanner = true

	// 初始化 JWT
	common.InitJWT(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// 设置统一错误处理
	e.HTTPErrorHandler = middleware.ErrorHandler()

	// 全局中间件
	e.Use(middleware.ZapLoggerMiddleware())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	// 公开路由（无需鉴权）
	handler.SetupPublicRoutes(e)

	// 受保护路由组（需要鉴权）
	protected := e.Group("/api", middleware.AuthMiddleware())
	handler.SetupRoutes(protected)

	// 启动服务
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	logger.Info("Server starting",
		zap.String("addr", addr),
		zap.String("mode", cfg.Server.Mode),
	)
	if err := e.Start(addr); err != nil {
		logger.Fatal("Failed to start server",
			zap.Error(err),
		)
	}
}
