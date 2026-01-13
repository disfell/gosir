package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        string         `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string         `json:"name" example:"张三"`
	Email     string         `json:"email" example:"zhangsan@example.com"`
	Password  string         `json:"-"`
	Phone     string         `json:"phone" example:"13800138000"`
	Avatar    string         `json:"avatar" example:"http://example.com/avatar.jpg"`
	Status    int            `json:"status" example:"1"`
	LastLogin *time.Time     `json:"last_login" example:"2026-01-08T10:00:00Z"`
	CreatedAt time.Time      `json:"created_at" example:"2026-01-08T10:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2026-01-08T10:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (User) TableName() string {
	return "users"
}
