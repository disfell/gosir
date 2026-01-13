# Model 目录结构说明

## 概述

`internal/model` 目录采用**按业务模块分类**的文件组织方式，每个模块都有独立的子文件夹。

## 目录结构

```
internal/model/
├── migration/               # 数据库迁移相关
│   └── migration.go         # SchemaMigration Model
└── user/                   # 用户模块
    ├── user.go              # 用户 Model 和 UserResponse
    └── user_status.go       # 用户状态常量和方法
```

## 各模块说明

### 1. migration/

**文件：** `migration.go`

**内容：**
```go
type SchemaMigration struct {
    Version   string    `gorm:"primaryKey;type:varchar(255)"`
    AppliedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
}
```

**用途：**
- 记录 SQL 脚本执行状态
- 由 `AutoMigrate` 自动创建和维护
- 保证 SQL 脚本的幂等性

### 2. user/

**文件：** `user.go`

**内容：**
```go
type User struct {
    ID        string         `gorm:"primaryKey;type:varchar(36)"`
    Name      string         `gorm:"type:varchar(255);not null"`
    Email     string         `gorm:"type:varchar(255);uniqueIndex;not null"`
    Password  string         `gorm:"type:varchar(255);not null"`
    Phone     string         `gorm:"type:varchar(20)"`
    Avatar    string         `gorm:"type:varchar(500)"`
    Status    int            `gorm:"type:tinyint;default:1"`
    LastLogin *time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserResponse struct {
    // 用于 API 返回，不包含敏感字段（如 Password）
}
```

**文件：** `user_status.go`

**内容：**
```go
type UserStatus int

const (
    UserStatusNormal  UserStatus = 1 // 正常
    UserStatusDisabled UserStatus = 2 // 禁用
)

func (s UserStatus) String() string {
    // 返回状态的中文描述
}
```

**用途：**
- 用户数据模型
- 用户状态枚举和转换方法

## 导入方式

### 在其他包中使用

由于 package name 可能与变量名冲突，建议使用 **import alias**：

```go
// 推荐：使用 alias 避免冲突
import usermodel "gosir/internal/model/user"

// 使用
user := &usermodel.User{
    Name:  "张三",
    Email: "zhangsan@example.com",
}
status := usermodel.UserStatusNormal
```

### 不推荐的方式（会有冲突）

```go
// 不推荐：package name 与局部变量冲突
import "gosir/internal/model/user"

// 编译错误
func CreateUser() {
    user := &user.User{}  // ❌ 变量名与 package name 冲突
}
```

## 新增模块指南

### 1. 创建模块文件夹

```bash
mkdir -p internal/model/order
```

### 2. 创建 Model 文件

**`internal/model/order/order.go`**

```go
package order

import (
    "time"
    "gorm.io/gorm"
)

type Order struct {
    ID        string    `gorm:"primaryKey;type:varchar(36)"`
    UserID    string    `gorm:"type:varchar(36);not null"`
    Total     float64   `gorm:"type:decimal(10,2);not null"`
    Status    int       `gorm:"type:tinyint;default:1"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type OrderResponse struct {
    ID        string  `json:"id"`
    UserID    string  `json:"user_id"`
    Total     float64 `json:"total"`
    Status    int     `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**`internal/model/order/order_status.go`**

```go
package order

type OrderStatus int

const (
    OrderStatusPending   OrderStatus = 1 // 待支付
    OrderStatusPaid      OrderStatus = 2 // 已支付
    OrderStatusCompleted OrderStatus = 3 // 已完成
    OrderStatusCancelled OrderStatus = 4 // 已取消
)

func (s OrderStatus) String() string {
    switch s {
    case OrderStatusPending:
        return "待支付"
    case OrderStatusPaid:
        return "已支付"
    case OrderStatusCompleted:
        return "已完成"
    case OrderStatusCancelled:
        return "已取消"
    default:
        return "未知"
    }
}
```

### 3. 创建数据库表

**`migrations/002_create_orders.sql`**

```sql
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
```

### 4. 在其他包中使用

```go
import ordermodel "gosir/internal/model/order"

func CreateOrder() {
    order := &ordermodel.Order{
        ID:     uuid.New().String(),
        UserID: "user-123",
        Total:  99.99,
        Status: int(ordermodel.OrderStatusPending),
    }
    // ...
}
```

## 注意事项

1. **Package Name 与模块名保持一致**
   - 文件夹：`user/`
   - Package name：`package user`

2. **使用 import alias 避免冲突**
   ```go
   import usermodel "gosir/internal/model/user"
   ```

3. **Model 分离**
   - `User`：数据模型（对应数据库表）
   - `UserResponse`：API 响应（不包含敏感字段）

4. **常量和枚举单独文件**
   - 如 `user_status.go`、`order_status.go`
   - 便于管理状态定义

5. **不使用 AutoMigrate 创建业务表**
   - 业务表通过 SQL 脚本创建
   - Model 只用于代码层面的类型安全

## 文件命名规范

| 文件名 | 用途 | 示例 |
|---------|------|------|
| `{module}.go` | 主 Model 和 Response | `user.go`、`order.go` |
| `{module}_status.go` | 状态常量和枚举 | `user_status.go` |
| `{module}_type.go` | 其他类型定义 | `order_type.go` |

## 总结

这种目录结构的优势：

1. **清晰的模块划分** - 每个业务模块独立
2. **易于维护** - 相关文件集中管理
3. **避免冲突** - 使用 import alias
4. **易于扩展** - 新增模块简单直观
5. **符合 Go 惯例** - 小而专注的 package
