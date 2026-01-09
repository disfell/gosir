#!/bin/bash

# Gosir Docker 一键部署脚本
# 整合构建、部署、清理等功能

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${GREEN}$1${NC}"
}

print_warn() {
    echo -e "${YELLOW}$1${NC}"
}

print_error() {
    echo -e "${RED}$1${NC}"
}

print_header() {
    echo ""
    echo "========================================="
    echo -e "${GREEN}$1${NC}"
    echo "========================================="
}

# 显示使用说明
show_usage() {
    cat << EOF
Gosir Docker 一键部署脚本

使用方法:
    $0 [命令]

可用命令:
    build       构建镜像
    deploy      构建并部署
    start       启动容器
    stop        停止容器
    restart     重启容器
    logs        查看日志
    clean       清理 Docker 资源
    status      查看状态
    help        显示帮助

示例:
    $0 deploy     # 一键部署
    $0 clean      # 清理资源

EOF
}

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "错误: Docker 未安装"
        exit 1
    fi
    
    if ! command -v docker compose &> /dev/null; then
        print_error "错误: Docker Compose 未安装"
        exit 1
    fi
}

# 清理 Docker 资源
docker_clean() {
    print_header "清理 Docker 资源"

    docker container prune -f >/dev/null 2>&1 || true
    docker image prune -a -f >/dev/null 2>&1 || true
    docker builder prune -f >/dev/null 2>&1 || true

    print_info "✅ 清理完成！"
}

# 构建镜像
docker_build() {
    print_header "构建 Docker 镜像"

    docker compose down >/dev/null 2>&1 || true
    docker builder prune -f >/dev/null 2>&1 || true

    echo "🔨 开始构建..."
    docker compose --progress=plain build
    print_info "✅ 构建完成！"
}

# 部署容器
docker_deploy() {
    print_header "部署 Docker 容器"

    if docker compose down >/dev/null 2>&1; then
        echo "🛑 已停止旧容器"
    fi

    docker compose up -d
    print_info "✅ 部署完成！"
    echo "   📡 应用地址: http://localhost:1323"
    echo "   📚 Swagger: http://localhost:1323/swagger/index.html"
}

# 一键构建并部署
docker_deploy_with_build() {
    print_header "一键构建并部署"

    if docker compose down >/dev/null 2>&1; then
        echo "🛑 已停止旧容器"
    fi

    docker container prune -f >/dev/null 2>&1 || true
    docker image prune -a -f >/dev/null 2>&1 || true
    docker builder prune -f >/dev/null 2>&1 || true

    echo "🔨 开始构建..."
    docker compose --progress=plain build
    echo "🚀 部署中..."
    docker compose up -d

    print_info "✅ 部署完成！"
    echo "   📡 应用地址: http://localhost:1323"
    echo "   📚 Swagger: http://localhost:1323/swagger/index.html"
}

# 启动容器
docker_start() {
    print_header "启动 Docker 容器"
    docker compose up -d
    print_info "✅ 容器已启动"
    echo "   📡 应用地址: http://localhost:1323"
}

# 停止容器
docker_stop() {
    print_header "停止 Docker 容器"
    docker compose down
    print_info "✅ 容器已停止"
}

# 重启容器
docker_restart() {
    print_header "重启 Docker 容器"

    if docker compose down >/dev/null 2>&1; then
        echo "🛑 已停止容器"
    fi

    docker compose up -d
    print_info "✅ 容器已重启"
    echo "   📡 应用地址: http://localhost:1323"
}

# 查看日志
docker_logs() {
    print_header "查看容器日志"
    docker compose logs -f
}

# 查看状态
docker_status() {
    print_header "Docker 状态"
    docker compose ps
}

# 主逻辑
main() {
    check_docker
    
    cd "$(dirname "$0")"  # 进入脚本所在目录（docker/）
    
    case "${1:-help}" in
        build)
            docker_build
            ;;
        deploy)
            docker_deploy
            ;;
        deploy-with-build)
            docker_deploy_with_build
            ;;
        start)
            docker_start
            ;;
        stop)
            docker_stop
            ;;
        restart)
            docker_restart
            ;;
        logs)
            docker_logs
            ;;
        clean)
            docker_clean
            ;;
        status)
            docker_status
            ;;
        help|--help|-h)
            show_usage
            ;;
        *)
            print_error "未知命令: $1"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

# 执行主逻辑
main "$@"
