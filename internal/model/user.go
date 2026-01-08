package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"` // 用户ID
	Name      string         `gorm:"type:varchar(255);not null" json:"name" example:"张三"`                                  // 姓名
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" example:"zhangsan@example.com"`   // 邮箱
	Password  string         `gorm:"type:varchar(255);not null" json:"-" example:"hashedpassword"`                         // 密码（不返回）
	Phone     string         `gorm:"type:varchar(20)" json:"phone" example:"13800138000"`                                  // 手机号
	Avatar    string         `gorm:"type:varchar(500)" json:"avatar" example:"http://example.com/avatar.jpg"`              // 头像
	Status    int            `gorm:"type:tinyint;default:1;comment:1:正常 2:禁用" json:"status" example:"1"`                   // 状态：1-正常 2-禁用
	LastLogin *time.Time     `json:"last_login" example:"2026-01-08T10:00:00Z"`                                            // 最后登录时间
	CreatedAt time.Time      `json:"created_at" example:"2026-01-08T10:00:00Z"`                                            // 创建时间
	UpdatedAt time.Time      `json:"updated_at" example:"2026-01-08T10:00:00Z"`                                            // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                                                       // 删除时间（不返回）
}

// UserResponse 用户响应（用于 API 返回，不包含敏感字段）
type UserResponse struct {
	ID        string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string     `json:"name" example:"张三"`
	Email     string     `json:"email" example:"zhangsan@example.com"`
	Phone     string     `json:"phone" example:"13800138000"`
	Avatar    string     `json:"avatar" example:"http://example.com/avatar.jpg"`
	Status    int        `json:"status" example:"1"`
	LastLogin *time.Time `json:"last_login" example:"2026-01-08T10:00:00Z"`
	CreatedAt time.Time  `json:"created_at" example:"2026-01-08T10:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2026-01-08T10:00:00Z"`
}

func (User) TableName() string {
	return "users"
}
