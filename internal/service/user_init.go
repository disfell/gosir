package service

import (
	"gosir/internal/database"
	"gosir/internal/model"
)

// InitAdminUser 初始化管理员账号（如果不存在）
func InitAdminUser() error {
	var count int64
	database.DB.Model(&model.User{}).Where("email = ?", "admin@example.com").Count(&count)

	if count > 0 {
		// 管理员已存在
		return nil
	}

	userService := NewUserService()
	_, err := userService.CreateUser("管理员", "admin@example.com", "admin123")
	return err
}
