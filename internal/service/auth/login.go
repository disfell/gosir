package auth

import (
	"gosir/internal/model"
	"gosir/internal/service/user"

	"golang.org/x/crypto/bcrypt"
)

// LoginService 登录服务
type LoginService struct {
	userService *user.UserService
}

// NewLoginService 创建登录服务
func NewLoginService() *LoginService {
	return &LoginService{
		userService: user.NewUserService(),
	}
}

// VerifyPassword 验证密码
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// LoginByPassword 通过邮箱密码登录
func (s *LoginService) LoginByPassword(email, password string) (*model.User, error) {
	userData, err := s.userService.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if err := VerifyPassword(userData.Password, password); err != nil {
		return nil, err
	}

	return userData, nil
}

// LoginByAccount 通过账号（邮箱或手机号）密码登录
func (s *LoginService) LoginByAccount(account, password string) (*model.User, error) {
	userData, err := s.userService.GetUserByEmailOrPhone(account)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if err := VerifyPassword(userData.Password, password); err != nil {
		return nil, err
	}

	return userData, nil
}
