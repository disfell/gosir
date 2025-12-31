package auth

import (
	"gosir/internal/common"
	"gosir/internal/middleware"
	"gosir/internal/service"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	userService *service.UserService
}

func New(userService *service.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// Login 登录
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return common.Error(c, common.CodeBadRequest, "请求参数解析失败")
	}

	if req.Email == "" || req.Password == "" {
		return common.Error(c, common.CodeValidationError, "邮箱和密码不能为空")
	}

	// 验证邮箱密码
	user, err := service.LoginByPassword(req.Email, req.Password)
	if err != nil {
		return common.Error(c, common.CodeUnauthorized, "邮箱或密码错误")
	}

	// 生成 token
	token, err := middleware.GenerateToken(user.Email)
	if err != nil {
		return common.Error(c, common.CodeInternalError, "生成 token 失败")
	}

	// 不返回密码
	user.Password = ""

	return common.Success(c, LoginResponse{
		Token: token,
		User:  user,
	})
}
