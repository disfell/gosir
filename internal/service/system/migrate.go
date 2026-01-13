package system

import (
	"go.uber.org/zap"

	"gosir/internal/database"
	"gosir/internal/logger"
	migrationmodel "gosir/internal/model/migration"
)

// AutoMigrate 自动迁移 schema_migrations 表
// 这是唯一通过 AutoMigrate 维护的表，用于记录 SQL 脚本执行状态
func AutoMigrate() error {
	logger.Info("Running AutoMigrate for schema_migrations table")

	if err := database.DB.AutoMigrate(
		&migrationmodel.SchemaMigration{},
	); err != nil {
		logger.Error("AutoMigrate failed", zap.Error(err))
		return err
	}

	logger.Info("AutoMigrate completed successfully")
	return nil
}
