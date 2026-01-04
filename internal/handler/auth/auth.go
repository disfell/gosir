package auth

import (
	"fmt"
	"gosir/internal/common"
	"gosir/internal/middleware"
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

// New 创建认证处理器
func New(userService *service.UserService) *Handler {
	validate := validator.New()

	// 获取中文翻译器
	zhLocale := zh.New()
	uni := ut.New(zhLocale, zhLocale)
	translator, ok := uni.GetTranslator("zh")
	if !ok {
		panic("failed to get chinese translator")
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

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string `json:"account" validate:"required"` // 账号（邮箱或手机号）
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

	// 解析请求
	if err := c.Bind(&req); err != nil {
		return common.Error(c, common.CodeBadRequest, "请求参数解析失败")
	}

	// 使用 validator 验证
	if err := h.validator.Struct(&req); err != nil {
		return common.Error(c, common.CodeValidationError, h.translateValidationError(err))
	}

	// 验证账号密码
	user, err := service.LoginByAccount(req.Account, req.Password)
	if err != nil {
		return common.Error(c, common.CodeUnauthorized, "账号或密码错误")
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

// translateValidationError 翻译验证错误
func (h *Handler) translateValidationError(err error) string {
	var fieldMessages []string

	for _, e := range err.(validator.ValidationErrors) {
		// 获取字段的中文名称
		fieldName := e.Field()
		var chineseField string
		switch fieldName {
		case "Account":
			chineseField = "账号"
		case "Password":
			chineseField = "密码"
		default:
			chineseField = fieldName
		}

		// 获取翻译后的错误消息
		msg := e.Translate(h.translator)
		if msg == "" {
			msg = fmt.Sprintf("%s验证失败", chineseField)
		}

		// 组合字段名和错误消息
		fieldMessages = append(fieldMessages, fmt.Sprintf("%s%s", chineseField, msg))
	}

	if len(fieldMessages) == 0 {
		return "参数验证失败"
	}

	return strings.Join(fieldMessages, "；")
}
