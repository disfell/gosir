package user

import (
	"fmt"
	"gosir/internal/common"
	"gosir/internal/service"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	userService *service.UserService
	validator   *validator.Validate
	translator  ut.Translator
}

func New(userService *service.UserService) *Handler {
	validate := validator.New()

	// 获取中文翻译器
	zhLocale := zh.New()
	uni := ut.New(zhLocale, zhLocale)
	translator, ok := uni.GetTranslator("zh")
	if !ok {
		panic("failed to get Chinese translator")
	}

	// 注册默认翻译
	if err := zhtranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		panic(fmt.Sprintf("failed to register translations: %v", err))
	}

	return &Handler{
		userService: userService,
		validator:   validate,
		translator:  translator,
	}
}

// 中文错误消息映射
var fieldNames = map[string]string{
	"Name":     "姓名",
	"Email":    "邮箱",
	"Password": "密码",
	"Phone":    "手机号",
	"Avatar":   "头像",
	"Status":   "状态",
}

func (h *Handler) translateValidationError(err error) string {
	var fieldMessages []string

	for _, e := range err.(validator.ValidationErrors) {
		fieldName := e.Field()
		chineseField := fieldNames[fieldName]
		if chineseField == "" {
			chineseField = fieldName
		}

		var errorMsg string
		switch e.Tag() {
		case "required":
			errorMsg = fmt.Sprintf("%s不能为空", chineseField)
		case "email":
			errorMsg = fmt.Sprintf("%s格式不正确", chineseField)
		case "min":
			errorMsg = fmt.Sprintf("%s长度不能少于%s个字符", chineseField, e.Param())
		case "max":
			errorMsg = fmt.Sprintf("%s长度不能超过%s个字符", chineseField, e.Param())
		default:
			errorMsg = fmt.Sprintf("%s验证失败: %s", chineseField, e.Tag())
		}
		fieldMessages = append(fieldMessages, errorMsg)
	}

	return strings.Join(fieldMessages, "；")
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required" example:"张三"`                                       // 姓名
	Email    string `json:"email" validate:"required,email" example:"zhangsan@example.com"`              // 邮箱
	Password string `json:"password" validate:"required,min=6" example:"password123"`                    // 密码
	Phone    string `json:"phone" validate:"omitempty,max=20" example:"13800138000"`                     // 手机号
	Avatar   string `json:"avatar" validate:"omitempty,max=500" example:"http://example.com/avatar.jpg"` // 头像
	Status   *int   `json:"status" validate:"omitempty,oneof=1 2" example:"1"`                           // 状态：1-正常 2-禁用
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name   string `json:"name" example:"李四"`                                  // 姓名
	Email  string `json:"email" example:"lisi@example.com"`                   // 邮箱
	Phone  string `json:"phone" example:"13900139000"`                        // 手机号
	Avatar string `json:"avatar" example:"http://example.com/new-avatar.jpg"` // 头像
	Status *int   `json:"status" example:"2"`                                 // 状态：1-正常 2-禁用
}

// GetUser 获取用户详情
// @Summary      获取用户详情
// @Description  根据用户ID获取用户详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path string true "用户ID"
// @Success      200 {object} common.Response{data=model.UserResponse}
// @Failure      404 {object} common.Response
// @Router       /api/users/{id} [get]
func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		return common.Error(c, common.CodeNotFound, "用户不存在")
	}
	return common.Success(c, user)
}

// CreateUser 创建用户
// @Summary      创建用户
// @Description  创建新用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body CreateUserRequest true "用户信息"
// @Success      201 {object} common.Response{data=model.UserResponse}
// @Failure      400 {object} common.Response
// @Failure      500 {object} common.Response
// @Router       /api/users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var req CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return common.Error(c, common.CodeBadRequest, "请求参数解析失败")
	}

	if err := h.validator.Struct(&req); err != nil {
		return common.Error(c, common.CodeValidationError, h.translateValidationError(err))
	}

	createReq := &service.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   req.Status,
	}

	user, err := h.userService.CreateUser(createReq)
	if err != nil {
		return common.Error(c, common.CodeInternalError, "创建用户失败")
	}
	return common.Created(c, user)
}

// ListUsers 获取用户列表
// @Summary      获取用户列表
// @Description  获取所有用户列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} common.Response{data=[]model.UserResponse}
// @Failure      500 {object} common.Response
// @Router       /api/users [get]
func (h *Handler) ListUsers(c echo.Context) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return common.Error(c, common.CodeInternalError, "获取用户列表失败")
	}
	return common.Success(c, users)
}

// UpdateUser 更新用户
// @Summary      更新用户
// @Description  更新用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path string true "用户ID"
// @Param        request body UpdateUserRequest true "用户信息"
// @Success      200 {object} common.Response{data=model.UserResponse}
// @Failure      400 {object} common.Response
// @Failure      404 {object} common.Response
// @Router       /api/users/{id} [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var req UpdateUserRequest

	if err := c.Bind(&req); err != nil {
		return common.Error(c, common.CodeBadRequest, "请求参数解析失败")
	}

	user, err := h.userService.UpdateUser(id, req.Name, req.Email, req.Phone, req.Avatar, req.Status)
	if err != nil {
		return common.Error(c, common.CodeNotFound, "用户不存在")
	}
	return common.Success(c, user)
}

// DeleteUser 删除用户
// @Summary      删除用户
// @Description  根据用户ID删除用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path string true "用户ID"
// @Success      200 {object} common.Response
// @Failure      404 {object} common.Response
// @Router       /api/users/{id} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.userService.DeleteUser(id)
	if err != nil {
		return common.Error(c, common.CodeNotFound, "用户不存在")
	}
	return common.SuccessWithMessage(c, "删除成功", nil)
}
