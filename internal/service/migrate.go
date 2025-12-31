package service

import (
	"gosir/internal/database"
	"gosir/internal/model"
)

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	return database.DB.AutoMigrate(
		&model.User{},
	)
}
