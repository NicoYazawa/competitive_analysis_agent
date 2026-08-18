#!/bin/bash
# stop.sh — 一键关闭所有服务
# 用法: ./scripts/stop.sh

set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "==> Stopping Frontend..."
pkill -f "vite.*3000" 2>/dev/null || true

echo "==> Stopping API server..."
pkill -f "bin/api" 2>/dev/null || true

echo "==> Stopping Worker..."
pkill -f "bin/worker" 2>/dev/null || true

echo "==> Stopping Docker services..."
docker compose down 2>/dev/null || true

echo "==> Cleaning up PID files..."
rm -f "$ROOT/bin/api" "$ROOT/bin/worker" 2>/dev/null || true

echo ""
echo "=== All services stopped ==="
