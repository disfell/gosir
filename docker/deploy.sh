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
    deploy      部署容器
    build_and_deploy 构建并部署
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

    # 输出 Docker 版本信息
    echo "🐳 Docker 版本: $(docker --version)"
    echo "📦 Docker Compose 版本: $(docker compose version)"
    echo ""
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

    # 使用时间戳作为版本号
    VERSION=$(date +%Y%m%d-%H%M%S)

    echo "🔨 开始构建..."
    echo "📦 版本号: $VERSION"

    docker compose --progress=plain build \
        --build-arg "VERSION=$VERSION"

    # 获取构建后的镜像名称（从 docker-compose.yml 读取）
    IMAGE_NAME=$(docker compose config | grep -A1 "image:" | head -n2 | tail -n1 | awk '{print $2}')

    if [ -z "$IMAGE_NAME" ]; then
        IMAGE_NAME="gosir:latest"
    fi

    # 给镜像打 tag
    docker tag "$IMAGE_NAME" "gosir:$VERSION"

    print_info "✅ 构建完成！"
    echo "   🏷️  镜像标签: gosir:$VERSION"
    echo "   🏷️  镜像标签: $IMAGE_NAME"
}

# 部署容器
docker_deploy() {
    print_header "部署 Docker 容器"
    docker compose up -d
    print_info "✅ 部署完成！"
    echo "   📡 应用地址: http://localhost:1323"
    echo "   📚 Swagger: http://localhost:1323/swagger/index.html"
}

# 一键构建并部署
docker_build_and_deploy() {
    docker_stop
    docker_build
    docker_deploy
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
        build_and_deploy)
            docker_build_and_deploy
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
        clean)
            docker_clean
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
