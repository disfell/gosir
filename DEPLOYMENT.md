# Gosir Docker 部署指南

## 前置要求

- Docker 20.10+
- Docker Compose 2.0+

## 部署步骤

### 1. 准备配置文件

编辑 `config/config.yaml`，根据生产环境调整配置：

```yaml
server:
  port: 1323
  mode: release  # 生产模式

database:
  path: data.db
  log_level: warn  # 生产环境建议使用 warn

jwt:
  secret: your-production-secret-key  # 必须修改为强密钥
  expire_hours: 24

log:
  level: info
  path: logs/app.log
  format: json
```

### 2. 使用 Docker Compose 部署（推荐）

```bash
# 进入 docker 目录
cd docker

# 构建并启动
docker compose up -d

# 查看日志
docker compose logs -f

# 停止服务
docker compose down

# 停止并删除数据
docker compose down -v
```

### 3. 手动构建和运行

```bash
# 构建镜像
docker build -f docker/Dockerfile -t gosir:latest .

# 运行容器
docker run -d \
  --name gosir-app \
  -p 1323:1323 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  -v $(pwd)/config/config.yaml:/app/config/config.yaml:ro \
  -e SERVER_MODE=release \
  --restart unless-stopped \
  gosir:latest
```

## 环境变量配置

支持通过环境变量覆盖配置：

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `SERVER_PORT` | 服务端口 | 1323 |
| `SERVER_MODE` | 运行模式 | release |
| `DATABASE_PATH` | 数据库路径 | /app/data/data.db |
| `DATABASE_LOG_LEVEL` | 数据库日志级别 | warn |
| `JWT_SECRET` | JWT密钥 | - |
| `JWT_EXPIRE_HOURS` | Token过期时间(小时) | 24 |
| `LOG_LEVEL` | 日志级别 | info |
| `LOG_PATH` | 日志文件路径 | /app/logs/app.log |
| `LOG_FORMAT` | 日志格式 | json |

## 数据持久化

项目使用 SQLite 数据库，数据存储在：

- 数据库文件：`./data/data.db`
- 日志文件：`./logs/app.log`

通过 Docker volume 挂载实现数据持久化。

## 健康检查

容器内置健康检查，访问 `http://localhost:1323/health` 检查服务状态。

```bash
# 查看容器健康状态
docker inspect --format='{{.State.Health.Status}}' gosir-app
```

## 生产环境建议

### 1. 安全配置

- **必须修改 JWT_SECRET**：使用强随机密钥
- **设置 SERVER_MODE=release**：关闭调试模式
- **限制日志级别**：生产环境使用 `warn` 或 `error`

### 2. 资源限制

在 `docker/docker-compose.yml` 中添加资源限制：

```yaml
gosir:
  # ... 其他配置
  deploy:
    resources:
      limits:
        cpus: '1'
        memory: 512M
      reservations:
        cpus: '0.5'
        memory: 256M
```

### 3. 日志管理

配置日志轮转，避免日志文件过大：

```yaml
gosir:
  # ... 其他配置
  logging:
    driver: "json-file"
    options:
      max-size: "10m"
      max-file: "3"
```

### 4. 备份策略

定期备份数据库文件：

```bash
# 备份数据库
docker exec gosir-app sh -c "cp /app/data/data.db /app/data/data.db.bak"
docker cp gosir-app:/app/data/data.db.bak ./backup/data.db.$(date +%Y%m%d)
```

### 5. 监控和告警

- 监控容器健康状态
- 监控日志文件大小
- 监控 API 响应时间

## 故障排查

### 查看日志

```bash
# Docker Compose
docker compose logs -f

# Docker
docker logs -f gosir-app
```

### 进入容器

```bash
docker exec -it gosir-app sh
```

### 检查端口

```bash
# 检查端口占用
netstat -tuln | grep 1323

# 或使用 lsof
lsof -i :1323
```

### 重建容器

```bash
# Docker Compose
docker compose up -d --build --force-recreate

# Docker
docker stop gosir-app
docker rm gosir-app
docker run -d --name gosir-app ...
```

## 性能优化

1. **构建优化**：使用多阶段构建，减小镜像体积
2. **缓存优化**：利用 Docker 层缓存，加速构建
3. **网络优化**：使用自定义网络，提高容器间通信效率

## 更新部署

```bash
# 进入 docker 目录
cd docker

# 拉取最新代码
git pull origin main

# 重新构建并启动
docker compose up -d --build

# 查看日志确认启动成功
docker compose logs -f
```

## 常用命令速查

```bash
# 进入 docker 目录
cd docker

# 构建并启动
docker compose up -d

# 重新构建并启动
docker compose up -d --build

# 强制重建
docker compose up -d --force-recreate

# 查看日志
docker compose logs -f

# 查看服务状态
docker compose ps

# 停止服务
docker compose stop

# 停止并删除容器
docker compose down

# 停止并删除容器及数据卷
docker compose down -v

# 查看资源使用
docker stats
```
