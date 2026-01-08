# Docker 构建常见问题

## 问题 1: Go 模块下载超时

### 错误信息
```
> [builder 5/7] RUN go mod download:
go: github.com/KyleBanks/depth@v1.2.1: Get "https://proxy.golang.org/...": dial tcp: i/o timeout
```

### 原因
国内访问 Google 的 Go 模块代理 `proxy.golang.org` 不稳定或被屏蔽。

### 解决方案 1: 使用国内镜像源（推荐）

在 Dockerfile 中添加环境变量：
```dockerfile
# 使用七牛云 Go 模块镜像
ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on
```

### 解决方案 2: 构建时传递环境变量

```bash
docker compose build \
  --build-arg GOPROXY=https://goproxy.cn,direct
```

### 解决方案 3: 使用阿里云镜像

```dockerfile
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/,direct \
    GO111MODULE=on
```

### 其他可用镜像源

| 镜像源 | 地址 |
|--------|------|
| 七牛云 | https://goproxy.cn |
| 阿里云 | https://mirrors.aliyun.com/goproxy/ |
| 腾讯云 | https://mirrors.cloud.tencent.com/go |
| 中国科大 | https://goproxy.io |
| 字节跳动 | https://goproxy.bolt.moe |

## 问题 2: DNS 解析失败

### 错误信息
```
dial tcp: lookup github.com: no such host
```

### 解决方案

配置 Docker 使用国内 DNS：
```yaml
# docker-compose.yml
services:
  gosir:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    dns:
      - 8.8.8.8
      - 114.114.114.114
```

或者在 Dockerfile 中：
```dockerfile
RUN echo "nameserver 8.8.8.8" > /etc/resolv.conf
```

## 问题 3: Git 拉取代码失败

### 错误信息
```
go: git@github.com:xxx/xxx.git: git@github.com:xxx/xxx.git is using git protocol
```

### 解决方案

配置 Git 使用 HTTPS 协议：
```dockerfile
RUN git config --global url."https://github.com/".insteadOf git@github.com:
```

## 问题 4: 缓存失效导致重复下载

### 解决方案

优化 Dockerfile 利用缓存层：
```dockerfile
# 先复制依赖文件，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 再复制源代码，代码变化不影响依赖下载
COPY . .
```

## 问题 5: 构建镜像过大

### 优化方案

1. 使用多阶段构建（已实现）
2. 减少不必要的依赖
3. 清理构建缓存

```dockerfile
# 构建后清理
RUN apk del --purge git && \
    rm -rf /var/cache/apk/*
```

## 完整的优化 Dockerfile 示例

```dockerfile
# 多阶段构建 - 构建阶段
FROM golang:1.25.5-alpine AS builder

# 设置工作目录
WORKDIR /app

# 配置国内镜像源和 DNS
ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on

# 安装必要的依赖
RUN apk add --no-cache git ca-certificates tzdata

# 配置 Git 使用 HTTPS
RUN git config --global url."https://github.com/".insteadOf git@github.com:

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖（这一层会被缓存）
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o gosir \
    cmd/server/main.go

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata sqlite-libs

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/gosir .

# 创建必要的目录
RUN mkdir -p logs data

# 设置权限
RUN chmod +x gosir

# 暴露端口
EXPOSE 1323

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:1323/health || exit 1

# 运行应用
CMD ["./gosir"]
```

## 构建命令

### 使用国内镜像源构建
```bash
cd docker
docker compose build
```

### 强制重新构建
```bash
docker compose build --no-cache
```

### 清理缓存后构建
```bash
docker compose build --no-cache --pull
```

## Docker 镜像加速配置

### 配置 Docker 镜像加速器

编辑 `/etc/docker/daemon.json`：
```json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.ccs.tencentyun.com"
  ]
}
```

重启 Docker：
```bash
sudo systemctl restart docker
```

### MacOS 配置

在 Docker Desktop 设置中：
1. 打开 Docker Desktop
2. Settings → Docker Engine
3. 添加镜像配置
4. Apply & Restart

## 常用命令

```bash
# 查看构建日志
docker compose build --progress=plain

# 查看镜像大小
docker images gosir

# 清理未使用的镜像
docker image prune -a

# 清理构建缓存
docker builder prune
```

## 网络问题排查

```bash
# 测试 Docker 网络连接
docker run --rm alpine ping -c 3 8.8.8.8

# 测试 DNS 解析
docker run --rm alpine nslookup github.com

# 测试代理连接
docker run --rm alpine wget -qO- https://goproxy.cn
```
