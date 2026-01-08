# Gosir

一个基于 Go 语言开发的 REST API 服务，使用 Echo 框架构建，提供用户认证、JWT 授权等基础功能。

## 技术栈

- **Go** 1.25.5
- **Echo** - Web 框架
- **GORM** - ORM 框架
- **SQLite** - 数据库
- **Zap** - 日志库
- **JWT** - 身份认证

## 项目结构

```
gosir/
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── config/
│   ├── config.go            # 配置加载
│   └── config.yaml          # 配置文件
├── internal/
│   ├── common/              # 公共组件
│   │   ├── error.go         # 错误处理
│   │   ├── jwt.go           # JWT 工具
│   │   └── response.go      # 响应封装
│   ├── database/            # 数据库初始化
│   ├── handler/             # HTTP 处理器
│   │   ├── auth/            # 认证相关
│   │   ├── system/          # 系统相关
│   │   └── user/            # 用户管理
│   ├── logger/              # 日志系统
│   ├── middleware/          # 中间件
│   │   ├── auth.go          # 认证中间件
│   │   ├── error_handler.go # 错误处理
│   │   └── echo_logger.go   # 日志中间件
│   ├── model/               # 数据模型
│   ├── repository/          # 数据访问层
│   └── service/             # 业务逻辑层
├── migrations/              # 数据库迁移
└── http/                    # HTTP 测试文件
```

## 快速开始

### 前置要求

- Go 1.25.5 或更高版本

### 安装依赖

```bash
go mod download
```

### 配置

复制并修改配置文件：

```bash
cp config/config.yaml config/local.yaml
```

编辑 `config/local.yaml`，修改以下配置：

```yaml
server:
  port: 1323                # 服务端口
  mode: debug               # 运行模式: debug, release, test

database:
  path: data.db             # SQLite 数据库文件路径
  log_level: info           # 数据库日志级别

jwt:
  secret: your-secret-key   # JWT 密钥（生产环境请修改）
  expire_hours: 24          # Token 过期时间（小时）

log:
  level: info               # 日志级别
  path: logs/app.log        # 日志文件路径
  format: json              # 日志格式: json, text
```

### 运行服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:1323` 启动。

### 构建

```bash
go build -o gosir cmd/server/main.go
```

## API 文档

### 公开接口

#### 健康检查
```
GET /health
```

#### 用户登录
```
POST /auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "your-password"
}
```

默认管理员账号：
- 用户名: `admin`
- 密码: 首次运行时查看日志中的随机密码

### 受保护接口

所有 `/api/*` 接口都需要在请求头中携带 JWT Token：

```
Authorization: Bearer <your-jwt-token>
```

#### 获取用户列表
```
GET /api/users
```

#### 创建用户
```
POST /api/users
Content-Type: application/json

{
  "name": "张三",
  "email": "zhangsan@example.com"
}
```

#### 获取用户详情
```
GET /api/users/:id
```

#### 更新用户
```
PUT /api/users/:id
Content-Type: application/json

{
  "name": "张三",
  "email": "zhangsan@example.com"
}
```

#### 删除用户
```
DELETE /api/users/:id
```

## 开发

### 运行测试

```bash
go test ./...
```

### HTTP 测试

项目包含 `http/test.http` 文件，可以使用 REST Client 插件进行 API 测试。

## 数据库迁移

数据库表会在服务启动时自动创建。首次启动时会：

1. 创建 `users` 表
2. 初始化管理员账号（admin）

## 日志

日志文件位于 `logs/app.log`，使用 JSON 格式记录。

## 安全注意事项

- 生产环境请修改 `config.yaml` 中的 JWT secret
- 配置文件包含敏感信息，已添加到 `.gitignore`
- 建议使用环境变量覆盖敏感配置

## Swagger 文档

项目集成了 Swagger API 文档，启动服务后访问：

```
http://localhost:1323/swagger/index.html
```

更新 Swagger 文档：

```bash
# 使用 swag 生成文档
swag init -g cmd/server/main.go -o docs

# 或使用 Makefile
make swagger
```

## Makefile 命令

项目提供了 Makefile 简化常用操作：

```bash
make build       # 构建应用
make run         # 运行应用
make test        # 运行测试
make clean       # 清理构建文件
make swagger     # 生成 Swagger 文档
docker-build     # 构建 Docker 镜像
docker-run       # 运行 Docker 容器
docker-down      # 停止 Docker 容器
```

## Docker 部署

项目支持 Docker 容器化部署，详细文档请查看 `DEPLOYMENT.md`。

快速启动：

```bash
cd docker
docker compose up -d
```

## License

MIT
