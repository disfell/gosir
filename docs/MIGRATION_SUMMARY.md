# 数据库迁移方案总结

## 设计理念

本项目采用**混合迁移策略**，结合了 GORM AutoMigrate 的便利性和 SQL 脚本的可控性。

### 核心原则

1. **AutoMigrate 仅用于 schema_migrations 表**
   - 这是唯一通过 AutoMigrate 创建的表
   - 用于追踪 SQL 脚本的执行状态
   - 表结构简单，不易出错

2. **所有业务表通过 SQL 脚本管理**
   - 完全控制表结构和索引策略
   - 支持版本控制和审计
   - 幂等性保障，可重复执行

## 架构对比

| 方面 | 纯 AutoMigrate | 纯 SQL 脚本 | 混合方案（当前） |
|------|---------------|-------------|----------------|
| 开发效率 | ✅ 高 | ⚠️ 中 | ✅ 高 |
| 生产安全 | ❌ 低 | ✅ 高 | ✅ 高 |
| 版本控制 | ❌ 无 | ✅ 有 | ✅ 有 |
| 索引优化 | ⚠️ 基础 | ✅ 完全可控 | ✅ 完全可控 |
| 团队协作 | ❌ 冲突风险 | ✅ 可审计 | ✅ 可审计 |
| 复杂操作 | ❌ 不支持 | ✅ 支持 | ✅ 支持 |

## 执行流程

```
应用启动
    ↓
1. AutoMigrate (创建 schema_migrations 表)
    ↓
2. ExecuteSQLScripts (执行 migrations/*.sql)
    ↓
3. InitAdminUser (初始化管理员)
    ↓
4. 启动服务
```

## 核心代码

### 1. AutoMigrate (internal/service/system/migrate.go)

```go
func AutoMigrate() error {
    logger.Info("Running AutoMigrate for schema_migrations table")

    if err := database.DB.AutoMigrate(
        &model.SchemaMigration{},  // 仅此表
    ); err != nil {
        logger.Error("AutoMigrate failed", zap.Error(err))
        return err
    }

    logger.Info("AutoMigrate completed successfully")
    return nil
}
```

### 2. SchemaMigration Model (internal/model/migration.go)

```go
type SchemaMigration struct {
    Version   string    `gorm:"primaryKey;type:varchar(255)" json:"version"`
    AppliedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"applied_at"`
}
```

### 3. SQL 脚本执行器 (internal/service/system/sql_runner.go)

核心功能：
- 扫描 migrations 文件夹
- 按文件名顺序执行
- 检查已执行记录（幂等性）
- 记录执行状态

## SQL 脚本示例

### 001_init_users.sql

```sql
-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar VARCHAR(500),
    status INTEGER DEFAULT 1,
    last_login DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- 复合索引
CREATE INDEX IF NOT EXISTS idx_users_status_created ON users(status, created_at);
```

## 使用场景

### ✅ 适合场景
- 中小型项目
- 需要 SQL 控制但不想引入复杂迁移工具
- 团队规模 2-10 人
- 需要快速迭代但保证生产安全

### ❌ 不适合场景
- 大型分布式系统（建议使用 golang-migrate）
- 需要复杂回滚逻辑
- 多个数据库类型支持
- 严格的变更审批流程

## 实际验证

### 第一次启动
```
✅ AutoMigrate 创建 schema_migrations 表
✅ SQL 脚本创建 users 表和索引
✅ 记录执行状态到 schema_migrations
✅ 初始化管理员账号
```

### 第二次启动
```
✅ AutoMigrate 检查 schema_migrations 表（无需更新）
✅ SQL 脚本跳过已执行的文件
✅ 快速启动
```

## 优势总结

1. **开发效率** - SQL 脚本比 golang-migrate 简单
2. **生产安全** - 幂等性保障，可重复执行
3. **可控性强** - 完全控制表结构和索引
4. **版本追踪** - 通过文件名和 schema_migrations 表
5. **易于理解** - 逻辑清晰，新人快速上手
6. **易于扩展** - 可随时切换到专业迁移工具

## 文件清单

```
gosir/
├── internal/
│   ├── model/
│   │   └── migration.go       # SchemaMigration Model
│   └── service/system/
│       ├── migrate.go         # AutoMigrate (schema_migrations)
│       └── sql_runner.go      # SQL 脚本执行器
├── migrations/
│   └── 001_init_users.sql     # 业务表 SQL 脚本
├── docs/
│   ├── MIGRATION.md           # 详细使用文档
│   └── MIGRATION_SUMMARY.md   # 本文档
└── cmd/server/main.go         # 启动入口
```

## 下一步优化

如果项目规模扩大，可以考虑：

1. **引入 golang-migrate**
   - 支持回滚操作
   - 更完善的版本控制
   - 支持多种数据库

2. **添加 up/down 脚本**
   - 001_xxx.up.sql
   - 001_xxx.down.sql

3. **CI/CD 集成**
   - 自动迁移验证
   - 预发布环境测试

4. **迁移验证工具**
   - 检查 SQL 语法
   - 性能影响评估

## 总结

当前方案是**实用主义**的体现：
- 不过度工程化
- 满足实际需求
- 易于维护和扩展
- 生产环境安全可靠

对于大多数 Go 项目，这是一个**性价比极高**的迁移方案。
