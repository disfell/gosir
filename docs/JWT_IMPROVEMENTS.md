# JWT 改进文档

## 改进概述

已基于本地缓存（`sync.Map`）实现了 JWT Token 黑名单机制，无需 Redis 即可实现 Token 撤销功能。

## 新增功能

### 1. JWT 管理器增强 (`internal/common/jwt.go`)

#### 新增字段
- `JTI` (JWT ID): 每个 Token 的唯一标识符
- `Issuer`: Token 签发者
- `blacklist`: 本地黑名单缓存（`sync.Map`）

#### 新增方法

| 方法 | 说明 |
|------|------|
| `AddToBlacklist(jti, expiredAt)` | 将 Token 加入黑名单 |
| `IsTokenBlacklisted(jti)` | 检查 Token 是否在黑名单中 |
| `CleanupExpiredBlacklist()` | 清理已过期的黑名单条目 |
| `GetBlacklistSize()` | 获取黑名单大小 |

### 2. 认证中间件增强 (`internal/middleware/auth.go`)

- 在验证 Token 时自动检查黑名单
- 黑名单中的 Token 返回 "token 已失效，请重新登录"
- 将 JTI 存入 Context 供后续使用

### 3. 新增 API 接口

#### 登出接口 (`POST /auth/logout`)
```go
// 将当前 Token 加入黑名单
// 强制失效，用户需重新登录
```

#### 刷新 Token 接口 (`POST /auth/refresh`)
```go
// 使用当前 Token 获取新的 Token
// 延长登录时间，无需重新登录
// 支持随时刷新，不限制时间
```

### 4. 定时清理任务

- **任务**: 每小时清理一次过期黑名单条目
- **频率**: `0 0 * * * *` (每小时整点)
- **日志**: 记录清理前后的黑名单大小

## 工作原理

### Token 生成流程
```
1. 用户登录成功
2. 生成 JTI (UUID)
3. 创建 JWT Claims 包含 JTI
4. 签发 Token 返回给客户端
```

### Token 验证流程
```
1. 从 Authorization Header 提取 Token
2. 解析并验证签名
3. 检查 Token 是否在黑名单中
4. 如果在黑名单 → 拒绝访问
5. 如果不在黑名单 → 通过验证
```

### Token 撤销流程（登出）
```
1. 用户调用 /auth/logout
2. 从 Context 获取 JTI
3. 将 JTI 加入黑名单
4. 下次请求时验证失败
```

### 黑名单自动清理
```
1. 定时任务每小时执行一次
2. 遍历黑名单条目
3. 删除已过期的 Token (ExpiredAt < 当前时间)
4. 释放内存空间
```

## API 使用示例

### 登录
```bash
POST /auth/login
{
  "account": "admin@example.com",
  "password": "password123"
}

Response:
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {...}
  }
}
```

### 登出
```bash
POST /auth/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...

Response:
{
  "code": 0,
  "message": "success",
  "data": null
}
```

### 刷新 Token
```bash
POST /auth/refresh
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}

Response:
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs... (新token)"
  }
}
```

## 配置说明

### JWT 配置 (`config/config.yaml`)
```yaml
jwt:
  secret: "your-secret-key-here"  # 签名密钥
  expire_hours: 24                # Token 有效期（小时）
```

### 环境变量
```bash
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24
```

## 安全性改进

| 改进项 | 说明 |
|--------|------|
| ✅ JTI | 每个 Token 唯一标识，支持精确撤销 |
| ✅ 黑名单 | 登出后 Token 立即失效 |
| ✅ 自动清理 | 定时清理过期黑名单，避免内存泄漏 |
| ✅ 签发者 | 添加 Issuer 字段，符合最佳实践 |
| ✅ 随时刷新 | 移除 Refresh 时间限制，提升体验 |

## 性能说明

### 本地缓存 vs Redis

| 特性 | 本地缓存 (sync.Map) | Redis |
|------|-------------------|-------|
| 响应速度 | 极快 (内存) | 快 (网络) |
| 部署复杂度 | 低 | 高 |
| 分布式支持 | ❌ 单机 | ✅ 支持 |
| 持久化 | ❌ 无 | ✅ 可选 |
| 适合场景 | 单机、测试环境 | 生产、多实例 |

### 内存占用估算
```
单条黑名单条目: ~200 bytes
1000 个用户登出: ~200 KB
10000 个用户登出: ~2 MB
```

**结论**: 对于中小规模应用，本地缓存完全够用。

## 迁移到 Redis

如果需要升级到 Redis，只需替换 `sync.Map` 为 Redis 客户端：

```go
// 替换黑名单实现
import "github.com/redis/go-redis/v9"

type JWTManager struct {
    secretKey string
    redis     *redis.Client
}

func (m *JWTManager) AddToBlacklist(jti string, expiredAt time.Duration) {
    ctx := context.Background()
    m.redis.Set(ctx, "blacklist:"+jti, "1", expiredAt)
}

func (m *JWTManager) IsTokenBlacklisted(jti string) bool {
    ctx := context.Background()
    exists, _ := m.redis.Exists(ctx, "blacklist:"+jti)
    return exists > 0
}
```

## 常见问题

### Q1: 重启服务后黑名单会丢失吗？
**A**: 是的，本地缓存在重启后会清空。如果需要持久化，请使用 Redis。

### Q2: 黑名单内存占用会持续增长吗？
**A**: 不会，定时任务每小时会自动清理过期的 Token。

### Q3: 多实例部署时黑名单如何同步？
**A**: 本地缓存不支持跨实例同步。多实例部署必须使用 Redis。

### Q4: 用户修改密码后如何让旧 Token 失效？
**A**: 在修改密码逻辑中调用 `AddToBlacklist` 即可。

### Q5: Token 刷新会生成新的 JTI 吗？
**A**: 是的，每次刷新都会生成新的 Token 和 JTI，旧 Token 自动失效。

## 测试建议

### 测试场景
1. 正常登录 → 获取 Token
2. 使用 Token 访问受保护资源 → 成功
3. 调用登出接口 → Token 加入黑名单
4. 使用原 Token 访问 → 失败，提示 "token 已失效"
5. 调用刷新 Token → 获取新 Token
6. 使用新 Token 访问 → 成功

### 测试黑名单清理
1. 登出多个 Token
2. 查看日志中的黑名单大小
3. 等待定时任务执行
4. 观察黑名单大小变化
