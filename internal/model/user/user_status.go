package model

// UserStatus 用户状态
type UserStatus int

const (
	UserStatusNormal   UserStatus = 1 // 正常
	UserStatusDisabled UserStatus = 2 // 禁用
)

func (s UserStatus) String() string {
	switch s {
	case UserStatusNormal:
		return "正常"
	case UserStatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}
