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
	zhLocale := zh.New()
	uni := ut.New(zhLocale, zhLocale)
	translator, _ := uni.GetTranslator("zh")
	err := zhtranslations.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		return nil
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

func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		return common.Error(c, common.CodeNotFound, "用户不存在")
	}
	return common.Success(c, user)
}

func (h *Handler) CreateUser(c echo.Context) error {
	var req struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.Bind(&req); err != nil {
		return common.Error(c, common.CodeBadRequest, "请求参数解析失败")
	}

	if err := h.validator.Struct(&req); err != nil {
		return common.Error(c, common.CodeValidationError, h.translateValidationError(err))
	}

	user, err := h.userService.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		return common.Error(c, common.CodeInternalError, "创建用户失败")
	}
	return common.Created(c, user)
}

func (h *Handler) ListUsers(c echo.Context) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return common.Error(c, common.CodeInternalError, "获取用户列表失败")
	}
	return common.Success(c, users)
}

func (h *Handler) UpdateUser(c echo.Context) error {
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

func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.userService.DeleteUser(id)
	if err != nil {
		return common.Error(c, common.CodeNotFound, "用户不存在")
	}
	return common.SuccessWithMessage(c, "删除成功", nil)
}
