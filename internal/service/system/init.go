package system

import (
	"gosir/internal/database"
	usermodel "gosir/internal/model/user"
	"gosir/internal/service/user"
)

// InitAdminUser 初始化管理员账号（如果不存在）
func InitAdminUser() error {
	var count int64
	database.DB.Model(&usermodel.User{}).Where("email = ?", "admin@gosir.com").Count(&count)

	if count > 0 {
		// 管理员已存在
		return nil
	}

	userService := user.NewUserService()
	createReq := &user.CreateUserRequest{
		Name:     "管理员",
		Email:    "admin@gosir.com",
		Password: "admin123",
		Phone:    "15578007781",
		Status:   nil, // 使用默认状态
	}
	_, err := userService.CreateUser(createReq)
	return err
}
