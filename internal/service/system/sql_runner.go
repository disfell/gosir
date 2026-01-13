package system

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.uber.org/zap"

	"gosir/internal/database"
	"gosir/internal/logger"
	"gosir/internal/repository"
)

// ExecuteSQLScripts 从指定文件夹按顺序执行所有 SQL 脚本
// 用于执行数据初始化、索引创建、视图等非核心结构的迁移
func ExecuteSQLScripts(folderPath string) error {
	logger.Info("Executing SQL scripts from folder", zap.String("folder", folderPath))

	// 获取数据库连接
	db := database.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 检查文件夹是否存在
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		logger.Warn("SQL scripts folder not found, skipping", zap.String("folder", folderPath))
		return nil
	}

	// 读取文件夹中的所有 .sql 文件
	var sqlFiles []string
	err = filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			sqlFiles = append(sqlFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read sql files: %w", err)
	}

	// 按文件名排序（确保按顺序执行）
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		logger.Info("No SQL scripts found in folder", zap.String("folder", folderPath))
		return nil
	}

	logger.Info("Found SQL scripts", zap.Int("count", len(sqlFiles)))

	// 初始化迁移记录仓储
	migrationRepo := repository.NewMigrationRepository()

	// 按顺序执行每个 SQL 文件
	for _, sqlFile := range sqlFiles {
		if err := executeSQLFile(sqlDB, migrationRepo, sqlFile); err != nil {
			return fmt.Errorf("failed to execute %s: %w", filepath.Base(sqlFile), err)
		}
	}

	logger.Info("All SQL scripts executed successfully")
	return nil
}

// executeSQLFile 执行单个 SQL 文件
func executeSQLFile(db *sql.DB, migrationRepo *repository.MigrationRepository, filePath string) error {
	// 读取 SQL 文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// 检查是否已执行
	version := filepath.Base(filePath)
	isExecuted, err := migrationRepo.IsExecuted(version)
	if err != nil {
		logger.Warn("Failed to check migration status",
			zap.String("file", filepath.Base(filePath)),
			zap.Error(err))
	}
	if isExecuted {
		logger.Info("SQL script already executed, skipping",
			zap.String("file", filepath.Base(filePath)))
		return nil
	}

	// 执行 SQL 语句
	// 分割 SQL 语句（支持多条语句）
	statements := splitSQLStatements(string(content))

	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		logger.Debug("Executing SQL statement",
			zap.String("file", filepath.Base(filePath)),
			zap.Int("statement", i+1),
			zap.Int("total", len(statements)),
			zap.String("sql", stmt))

		if _, err := db.Exec(stmt); err != nil {
			// 如果表已存在，忽略错误（幂等性）
			if !isTableExistsError(err) {
				logger.Error("Failed to execute SQL statement",
					zap.String("file", filepath.Base(filePath)),
					zap.Int("statement", i+1),
					zap.String("sql", stmt),
					zap.Error(err))
				return fmt.Errorf("failed to execute statement: %w", err)
			}
			logger.Info("SQL statement skipped (table already exists)",
				zap.String("file", filepath.Base(filePath)),
				zap.Int("statement", i+1))
		}
	}

	// 记录已执行的迁移
	if err := migrationRepo.RecordMigration(version); err != nil {
		logger.Warn("Failed to record migration",
			zap.String("file", filepath.Base(filePath)),
			zap.Error(err))
	}

	logger.Info("SQL script executed successfully",
		zap.String("file", filepath.Base(filePath)))
	return nil
}

// splitSQLStatements 分割 SQL 语句
func splitSQLStatements(content string) []string {
	var statements []string
	var current strings.Builder

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过注释和空行
		if strings.HasPrefix(line, "--") || line == "" {
			continue
		}

		current.WriteString(line)
		current.WriteString("\n")

		// 检查语句结束
		if strings.HasSuffix(line, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	// 处理最后没有分号的语句
	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}

// isTableExistsError 检查是否是表已存在的错误
func isTableExistsError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "already exists") ||
		strings.Contains(errMsg, "duplicate table")
}
