package model

import "time"

// SchemaMigration 数据库迁移记录模型
type SchemaMigration struct {
	Version   string    `gorm:"primaryKey;type:varchar(255)" json:"version"`               // 迁移版本号（文件名）
	AppliedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"applied_at"` // 应用时间
}

func (SchemaMigration) TableName() string {
	return "schema_migrations"
}
