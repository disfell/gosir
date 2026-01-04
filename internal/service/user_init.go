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
	createReq := &CreateUserRequest{
		Name:     "管理员",
		Email:    "admin@example.com",
		Password: "admin123",
		Status:   nil, // 使用默认状态
	}
	_, err := userService.CreateUser(createReq)
	return err
}
