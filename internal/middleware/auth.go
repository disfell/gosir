package middleware

import (
	"myapp/internal/common"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// jwtManager JWT 管理器实例
var jwtManager *common.JWTManager

// InitJWT 初始化 JWT 管理器
func InitJWT(secretKey string, expiryHours int) {
	jwtManager = common.NewJWTManager(secretKey, time.Duration(expiryHours)*time.Hour)
}

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			// 检查 Authorization 请求头
			if authHeader == "" {
				return common.Error(c, common.CodeUnauthorized, "缺少 Authorization 请求头")
			}

			// 检查 Bearer 前缀
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return common.Error(c, common.CodeUnauthorized, "无效的 Authorization 格式，应为: Bearer {token}")
			}

			tokenString := parts[1]

			// 验证 token
			claims, err := jwtManager.ValidateToken(tokenString)
			if err != nil {
				return common.Error(c, common.CodeUnauthorized, "无效的 token: "+err.Error())
			}

			// 将用户信息存入 context
			c.Set("user_id", claims.UserID)
			c.Set("claims", claims)

			return next(c)
		}
	}
}

// RequireAuth 保留此函数用于特定路由的二次验证（可选）
func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("user_id")
			if userID == nil {
				return common.Error(c, common.CodeUnauthorized, "unauthorized")
			}
			return next(c)
		}
	}
}

// GenerateToken 生成 JWT token（供 handler 使用）
func GenerateToken(userID string) (string, error) {
	return jwtManager.GenerateToken(userID)
}

// ValidateToken 验证 JWT token（供 handler 使用）
func ValidateToken(tokenString string) (*common.JWTClaims, error) {
	return jwtManager.ValidateToken(tokenString)
}

// GetUserID 从 context 获取用户 ID
func GetUserID(c echo.Context) string {
	if userID, ok := c.Get("user_id").(string); ok {
		return userID
	}
	return ""
}



// GetClaims 从 context 获取 JWT 声明
func GetClaims(c echo.Context) *common.JWTClaims {
	if claims, ok := c.Get("claims").(*common.JWTClaims); ok {
		return claims
	}
	return nil
}
