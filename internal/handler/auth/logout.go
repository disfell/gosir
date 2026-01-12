package auth

import (
	"gosir/internal/common"

	"github.com/labstack/echo/v4"
)

// Logout 登出
// @Summary      用户登出
// @Description  将当前 token 加入黑名单，强制失效
// @Tags         认证
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} common.Response
// @Failure      401 {object} common.Response
// @Router       /auth/logout [post]
func (h *Handler) Logout(c echo.Context) error {
	// 从 context 获取 claims
	claims, ok := c.Get("claims").(*common.JWTClaims)
	if !ok {
		return common.Error(c, common.CodeUnauthorized, "无效的认证信息")
	}

	// 将 token 加入黑名单
	jwtManager := common.GetJWTManager()
	if jwtManager == nil {
		return common.Error(c, common.CodeInternalError, "JWT 管理器未初始化")
	}

	jwtManager.AddToBlacklist(claims.JTI, claims.ExpiresAt.Time)

	return common.Success(c, nil)
}
