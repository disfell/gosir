package cron

import (
	"fmt"
	"gosir/internal/common"
	"gosir/internal/logger"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Manager 定时任务管理器
type Manager struct {
	cron *cron.Cron
}

var manager *Manager

// Init 初始化定时任务管理器
func Init() {
	manager = &Manager{
		cron: cron.New(cron.WithSeconds()), // 支持秒级精度
	}

	// 注册所有定时任务
	manager.registerTasks()

	// 启动定时任务
	manager.cron.Start()

	logger.Info("Cron manager started")
}

// Stop 停止定时任务管理器
func Stop() {
	if manager != nil && manager.cron != nil {
		manager.cron.Stop()
		logger.Info("Cron manager stopped")
	}
}

// registerTasks 注册所有定时任务
func (cm *Manager) registerTasks() {
	// JWT 黑名单清理任务 - 每小时执行一次
	cm.addJob("0 0 * * * *", "清理过期的 token 黑名单", cm.cleanupExpiredBlacklistTask)

	// 示例1: 每5秒执行一次
	cm.addJob("*/5 * * * * *", "每5秒执行的任务", cm.everyFiveSecondsTask)

	// 示例2: 每分钟执行一次
	cm.addJob("0 * * * * *", "每分钟执行的任务", cm.everyMinuteTask)

	// 示例3: 每小时执行一次
	cm.addJob("0 0 * * * *", "每小时执行的任务", cm.everyHourTask)

	// 示例4: 每天凌晨2点执行
	cm.addJob("0 0 2 * * *", "每天凌晨2点执行的任务", cm.dailyTask)

	// 示例5: 使用 Cron 表达式 (每天上午10点30分)
	cm.addJob("0 30 10 * * *", "每天上午10点30分", cm.scheduledTask)
}

// addJob 添加任务
func (cm *Manager) addJob(spec string, name string, job func()) {
	_, err := cm.cron.AddFunc(spec, job)
	if err != nil {
		logger.Error("Failed to add cron job",
			zap.String("name", name),
			zap.String("spec", spec),
			zap.Error(err),
		)
		return
	}

	logger.Info("Cron job registered",
		zap.String("name", name),
		zap.String("spec", spec),
	)
}

// ==================== 示例任务函数 ====================

// cleanupExpiredBlacklistTask 清理过期的黑名单 token
func (cm *Manager) cleanupExpiredBlacklistTask() {
	jwtManager := common.GetJWTManager()
	if jwtManager == nil {
		logger.Error("JWT manager not initialized for cleanup")
		return
	}

	sizeBefore := jwtManager.GetBlacklistSize()
	jwtManager.CleanupExpiredBlacklist()
	sizeAfter := jwtManager.GetBlacklistSize()

	logger.Info("Blacklist cleanup completed",
		zap.Int("before", sizeBefore),
		zap.Int("after", sizeAfter),
		zap.Int("cleaned", sizeBefore-sizeAfter),
	)
}

// everyFiveSecondsTask 每5秒执行一次的任务
func (cm *Manager) everyFiveSecondsTask() {
	logger.Debug("执行每5秒任务", zap.String("task", "everyFiveSeconds"))
	// 在这里添加你的业务逻辑
	// 例如: 清理缓存、检查状态、发送心跳等
}

// everyMinuteTask 每分钟执行一次的任务
func (cm *Manager) everyMinuteTask() {
	logger.Debug("执行每分钟任务", zap.String("task", "everyMinute"))
	// 例如: 定期统计数据、同步信息等
}

// everyHourTask 每小时执行一次的任务
func (cm *Manager) everyHourTask() {
	logger.Debug("执行每小时任务", zap.String("task", "everyHour"))
	// 例如: 生成报表、备份数据等
}

// dailyTask 每天凌晨2点执行的任务
func (cm *Manager) dailyTask() {
	logger.Debug("执行每日任务", zap.String("task", "daily"))
	// 例如: 数据归档、日志清理、定期维护等
}

// scheduledTask 定时执行的任务
func (cm *Manager) scheduledTask() {
	logger.Debug("执行定时任务", zap.String("task", "scheduled"))
	// 在这里添加你的具体业务逻辑
}

// ==================== 高级用法示例 ====================

// Cron 表达式格式: 秒 分 时 日 月 周
// 字段说明:
// 秒: 0-59
// 分: 0-59
// 时: 0-23
// 日: 1-31
// 月: 1-12
// 周: 0-7 (0和7都代表周日)

// 常用示例:
// "*/5 * * * *"        # 每5分钟
// "0 * * * *"          # 每小时
// "0 0 * * *"          # 每天凌晨
// "0 0 2 * * *"        # 每天凌晨2点
// "0 0 2 * * 1"        # 每周一凌晨2点
// "0 0 0 1 * *"        # 每月1号凌晨
// "0 */30 9-17 * * *"  # 工作时间每30分钟 (9点到17点)
// "0 30 10 * * *"      # 每天上午10点30分
// "0 0 12 * * MON-FRI" # 周一到周五中午12点

// AddCustomJob 添加自定义任务 (供外部调用)
func AddCustomJob(spec string, name string, job func()) error {
	if manager == nil {
		return fmt.Errorf("cron manager not initialized")
	}

	id, err := manager.cron.AddFunc(spec, job)
	if err != nil {
		return err
	}

	logger.Info("Custom cron job added",
		zap.String("name", name),
		zap.String("spec", spec),
		zap.Int64("id", int64(id)),
	)

	return nil
}
