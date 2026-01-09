#!/bin/bash

# Docker 清理脚本
# 定期清理 Docker 的无用资源，释放磁盘空间

set -e

echo "========================================="
echo "Docker 清理脚本"
echo "========================================="
echo ""

# 1. 清理停止的容器
echo "1. 清理停止的容器..."
docker container prune -f
echo "✅ 完成"
echo ""

# 2. 清理未使用的镜像
echo "2. 清理未使用的镜像..."
docker image prune -a -f
echo "✅ 完成"
echo ""

# 3. 清理构建缓存（这是占用空间最大的）
echo "3. 清理构建缓存..."
docker builder prune -f
echo "✅ 完成"
echo ""

# 4. 清理未使用的卷（谨慎使用）
echo "4. 清理未使用的卷（跳过，避免误删数据）..."
# docker volume prune -f  # 注释掉，避免删除数据
echo "⚠️  已跳过（包含数据）"
echo ""

# 5. 清理所有未使用资源
echo "5. 清理所有未使用资源..."
docker system prune -a -f
echo "✅ 完成"
echo ""

# 显示清理后的空间使用情况
echo "========================================="
echo "清理后的空间使用情况："
echo "========================================="
docker system df
echo ""

echo "========================================="
echo "✅ 清理完成！"
echo "========================================="
