package middleware

import (
	"gosir/internal/common"
	"strings"

	"github.com/labstack/echo/v4"
)

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
			jwtManager := common.GetJWTManager()
			if jwtManager == nil {
				return common.Error(c, common.CodeUnauthorized, "JWT 管理器未初始化")
			}

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
