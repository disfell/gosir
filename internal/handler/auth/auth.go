package auth

import (
	"myapp/internal/common"
	"myapp/internal/middleware"

	"github.com/labstack/echo/v4"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
}

// Login 登录
func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return common.Error(c, common.CodeBadRequest, "请求参数解析失败")
	}

	// 简单示例：验证用户名密码
	// 实际项目中应该从数据库查询用户并验证密码哈希
	if req.Username == "" || req.Password == "" {
		return common.Error(c, common.CodeValidationError, "用户名和密码不能为空")
	}

	// 示例：这里应该是查询数据库验证用户
	// user, err := userService.GetByUsername(req.Username)
	// if err != nil || !user.CheckPassword(req.Password) {
	//     return common.Error(c, common.CodeUnauthorized, "用户名或密码错误")
	// }

	// 示例：演示模式，任意非空用户名密码都通过
	// 生成 token
	token, err := middleware.GenerateToken(req.Username)
	if err != nil {
		return common.Error(c, common.CodeInternalError, "生成 token 失败")
	}

	return common.Success(c, LoginResponse{Token: token})
}
