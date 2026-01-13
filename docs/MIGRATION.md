# 数据库迁移说明

本项目采用**混合迁移策略**，结合 GORM AutoMigrate 和 SQL 脚本的优势。

## 迁移策略

### 1. AutoMigrate - 仅用于 schema_migrations 表
**负责：** 只创建和维护 `schema_migrations` 表

**文件位置：** `internal/model/migration.go`

**优势：**
- 通过 GORM 管理，与代码同步
- 表结构简单，不易出错
- 为 SQL 脚本提供执行状态追踪

**Model 定义：**
```go
type SchemaMigration struct {
    Version   string    `gorm:"primaryKey;type:varchar(255)"`
    AppliedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
}
```

### 2. SQL 脚本 - 所有业务表和优化
**负责：** 所有业务表的创建、索引、数据初始化、视图、触发器等

**文件位置：** `migrations/*.sql`

**优势：**
- 完全控制表结构和索引策略
- 支持复杂的 SQL 操作
- 可追踪执行状态（幂等性）
- 适合团队协作和版本控制

**示例：**
```sql
-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    -- ...
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status_created ON users(status, created_at);

-- 初始化数据
INSERT INTO config (key, value) VALUES ('app_name', 'Gosir');
```

## 迁移执行顺序

应用启动时按以下顺序执行：

```
1. AutoMigrate (创建 schema_migrations 表)
   ↓
2. ExecuteSQLScripts (执行 migrations 文件夹中的所有 SQL 脚本)
   ↓
3. InitAdminUser (初始化管理员账号)
   ↓
4. 启动服务
```

## 使用指南

### 添加新表

**方式一：使用 SQL 脚本（推荐）**
```sql
-- migrations/002_create_orders.sql
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
```

**方式二：定义 Model（可选，用于代码层面）**
```go
// internal/model/order.go
type Order struct {
    ID        string    `gorm:"primaryKey;type:varchar(36)"`
    UserID    string    `gorm:"type:varchar(36);not null"`
    Total     float64   `gorm:"type:decimal(10,2)"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```
> 注意：Model 定义只用于代码层面的类型安全，表结构由 SQL 脚本创建。

### 添加新字段

**使用 SQL 脚本**
```sql
-- migrations/003_add_age_column.sql
ALTER TABLE users ADD COLUMN age INTEGER DEFAULT 0;

-- 如果新字段需要索引
CREATE INDEX IF NOT EXISTS idx_users_age ON users(age);
```

### 初始化数据

```sql
-- migrations/004_seed_data.sql
INSERT INTO config (key, value) VALUES 
    ('app_name', 'Gosir'),
    ('max_users', '1000');
```

### SQL 脚本命名规范

```
<序列号>_<描述>.sql

示例：
001_init_users.sql
002_order_indexes.sql
003_seed_data.sql
```

## 幂等性保障

- **AutoMigrate**：GORM 自动处理，不会重复创建
- **SQL 脚本**：使用 `IF NOT EXISTS` 确保幂等
- **执行记录**：`schema_migrations` 表记录已执行的脚本

## 生产环境建议

### 开发/测试环境
- 使用 AutoMigrate 快速迭代
- SQL 脚本用于优化测试

### 生产环境
- **首次部署**：AutoMigrate + SQL 脚本
- **后续更新**：使用 golang-migrate 等专业工具
- **备份**：执行迁移前备份数据库

## 最佳实践

1. **schema_migrations 表** → 使用 AutoMigrate（唯一例外）
2. **业务表结构** → 使用 SQL 脚本
3. **索引优化** → 使用 SQL 脚本
4. **种子数据** → 使用 SQL 脚本
5. **视图/存储过程** → 使用 SQL 脚本
6. **复杂迁移** → 使用专业迁移工具（golang-migrate）

### 为什么这样设计？

| 方案 | AutoMigrate | SQL 脚本 |
|------|-------------|----------|
| 表结构控制 | 自动生成，依赖 GORM | 完全手动控制 |
| 索引策略 | 基础索引 | 支持复合索引、覆盖索引 |
| 性能优化 | 有限 | 完全可控 |
| 版本控制 | ❌ 无 | ✅ 通过文件名 |
| 团队协作 | ⚠️ 冲突风险 | ✅ 可审计 |
| 生产安全 | ⚠️ 风险高 | ✅ 幂等性保障 |

## 项目结构

```
gosir/
├── internal/
│   ├── model/
│   │   ├── migration.go       # SchemaMigration Model (AutoMigrate)
│   │   ├── user.go            # 用户 Model (代码层面)
│   │   └── ...                # 其他 Model
│   └── service/
│       └── system/
│           ├── migrate.go     # AutoMigrate (仅 schema_migrations)
│           └── sql_runner.go  # SQL 脚本执行器
├── migrations/                 # 所有业务表的 SQL 脚本
│   ├── 001_init_users.sql     # 用户表 + 索引
│   ├── 002_create_orders.sql  # 订单表 + 索引
│   └── ...
└── cmd/server/main.go         # 启动入口
```

## FAQ

### Q: 为什么不在 AutoMigrate 里维护业务表？
A: 
- AutoMigrate 缺乏版本控制，无法追踪变更历史
- 无法精细控制索引策略（复合索引、覆盖索引等）
- 生产环境风险高，可能导致数据丢失
- 团队协作时容易出现冲突

### Q: SQL 脚本如何保证幂等性？
A: 使用 `IF NOT EXISTS` 和 `CREATE INDEX IF NOT EXISTS` 等语法

### Q: 如何回滚迁移？
A: 
- 开发环境：删除数据库重新运行
- 生产环境：使用 `schema_migrations` 表记录版本，配合备份恢复

### Q: Model 定义是否还需要？
A: 需要！Model 用于代码层面的类型安全和 ORM 操作，但表结构由 SQL 脚本创建
