package auth

import (
	"gosir/internal/common"
	"strings"

	"github.com/labstack/echo/v4"
)

// RefreshTokenRequest 刷新 token 请求
type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // 当前 token
}

// RefreshTokenResponse 刷新 token 响应
type RefreshTokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // 新 token
}

// RefreshToken 刷新 token
// @Summary      刷新 token
// @Description  使用当前 token 获取新的 token，延长登录时间
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "刷新 token 信息"
// @Success      200 {object} common.Response{data=RefreshTokenResponse}
// @Failure      400 {object} common.Response
// @Failure      401 {object} common.Response
// @Router       /auth/refresh [post]
func (h *Handler) RefreshToken(c echo.Context) error {
	var req RefreshTokenRequest

	// 解析请求
	if err := c.Bind(&req); err != nil {
		return common.Error(c, common.CodeBadRequest, "请求参数解析失败")
	}

	// 验证参数
	if err := h.validator.Struct(&req); err != nil {
		return common.Error(c, common.CodeValidationError, h.translateValidationError(err))
	}

	// 刷新 token
	jwtManager := common.GetJWTManager()
	if jwtManager == nil {
		return common.Error(c, common.CodeInternalError, "JWT 管理器未初始化")
	}

	newToken, err := jwtManager.RefreshToken(req.Token)
	if err != nil {
		if strings.Contains(err.Error(), "已失效") {
			return common.Error(c, common.CodeUnauthorized, "token 已失效，请重新登录")
		}
		return common.Error(c, common.CodeUnauthorized, "刷新 token 失败: "+err.Error())
	}

	return common.Success(c, RefreshTokenResponse{
		Token: newToken,
	})
}
