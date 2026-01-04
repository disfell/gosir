package service

import (
	"gosir/internal/database"
	"gosir/internal/model"

	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword 验证密码
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// LoginByPassword 通过邮箱密码登录
func LoginByPassword(email, password string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	// 验证密码
	if err := VerifyPassword(user.Password, password); err != nil {
		return nil, err
	}

	return &user, nil
}

// LoginByAccount 通过账号（邮箱或手机号）密码登录
func LoginByAccount(account, password string) (*model.User, error) {
	var user model.User
	// 查找邮箱或手机号匹配的用户
	err := database.DB.Where("email = ? OR phone = ?", account, account).First(&user).Error
	if err != nil {
		return nil, err
	}

	// 验证密码
	if err := VerifyPassword(user.Password, password); err != nil {
		return nil, err
	}

	return &user, nil
}
