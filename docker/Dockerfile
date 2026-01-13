# 多阶段构建 - 构建阶段
FROM golang:1.25.5-alpine AS builder

# 设置工作目录
WORKDIR /app

# 配置国内 Go 模块镜像源（加速下载）
ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on \
    CGO_ENABLED=1

# 安装必要的依赖（包括 GCC 用于 CGO）
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    build-base \
    sqlite-dev

# 安装 swag 工具用于生成 Swagger 文档
RUN go install github.com/swaggo/swag/cmd/swag@latest

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 生成 Swagger 文档（如果代码中有 swagger 注解）
ENV PATH="${PATH}:/root/go/bin"
RUN swag init -g cmd/server/main.go --parseDependency --parseInternal -o docs || echo "Swagger generation skipped or failed, continuing..."

# 构建应用（启用 CGO 以支持 SQLite）
RUN go build \
    -ldflags="-w -s" \
    -o gosir \
    cmd/server/main.go

# 运行阶段
FROM alpine:3.19

# 安装必要的依赖
RUN apk add --no-cache ca-certificates tzdata sqlite-libs

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/gosir .

# 创建必要的目录
RUN mkdir -p logs data config

# 设置权限
RUN chmod +x gosir

# 暴露端口
EXPOSE 1323

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:1323/health || exit 1

# 运行应用
CMD ["./gosir"]
