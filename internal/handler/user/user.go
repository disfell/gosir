package user

import (
	"myapp/internal/common"
	"myapp/internal/service"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *service.UserService
}

func New(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		return common.Error(c, common.CodeNotFound, "用户不存在")
	}
	return common.Success(c, user)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min:6"`
	}
	if err := c.Bind(&req); err != nil {
		return common.Error(c, common.CodeBadRequest, "请求参数解析失败")
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return common.Error(c, common.CodeValidationError, "姓名、邮箱和密码不能为空")
	}
	user, err := h.userService.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		return common.Error(c, common.CodeInternalError, "创建用户失败")
	}
	return common.Created(c, user)
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return common.Error(c, common.CodeInternalError, "获取用户列表失败")
	}
	return common.Success(c, users)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var req struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Avatar string `json:"avatar"`
		Status *int   `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return common.Error(c, common.CodeBadRequest, "请求参数解析失败")
	}
	user, err := h.userService.UpdateUser(id, req.Name, req.Email, req.Phone, req.Avatar, req.Status)
	if err != nil {
		return common.Error(c, common.CodeNotFound, "用户不存在")
	}
	return common.Success(c, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.userService.DeleteUser(id)
	if err != nil {
		return common.Error(c, common.CodeNotFound, "用户不存在")
	}
	return common.SuccessWithMessage(c, "删除成功", nil)
}
