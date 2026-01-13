package repository

import (
	"gosir/internal/database"
	migrationmodel "gosir/internal/model/migration"

	"gorm.io/gorm"
)

// MigrationRepository 迁移记录仓储层
type MigrationRepository struct {
	db *gorm.DB
}

// NewMigrationRepository 创建迁移记录仓储实例
func NewMigrationRepository() *MigrationRepository {
	return &MigrationRepository{
		db: database.DB,
	}
}

// IsExecuted 检查迁移是否已执行
func (r *MigrationRepository) IsExecuted(version string) (bool, error) {
	var count int64
	err := r.db.Model(&migrationmodel.SchemaMigration{}).
		Where("version = ?", version).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// RecordMigration 记录已执行的迁移
func (r *MigrationRepository) RecordMigration(version string) error {
	migrationRecord := &migrationmodel.SchemaMigration{
		Version: version,
	}
	return r.db.Create(migrationRecord).Error
}

// GetAll 获取所有迁移记录
func (r *MigrationRepository) GetAll() ([]*migrationmodel.SchemaMigration, error) {
	var records []*migrationmodel.SchemaMigration
	err := r.db.Order("applied_at ASC").Find(&records).Error
	return records, err
}

// Delete 删除指定版本的迁移记录
func (r *MigrationRepository) Delete(version string) error {
	return r.db.Where("version = ?", version).
		Delete(&migrationmodel.SchemaMigration{}).Error
}
