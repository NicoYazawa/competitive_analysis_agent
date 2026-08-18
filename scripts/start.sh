#!/bin/bash
# start.sh — 一键启动所有服务
# 用法: ./scripts/start.sh

set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "==> Starting Docker services (PostgreSQL + Redis)..."
docker compose up -d

echo "==> Waiting for PostgreSQL to be ready..."
for i in {1..30}; do
  if docker exec competitive_analysis_pg pg_isready -U postgres -p 5432 > /dev/null 2>&1; then
    echo "    PostgreSQL is ready"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "ERROR: PostgreSQL did not become ready in time"
    exit 1
  fi
  sleep 1
done

echo "==> Waiting for Redis to be ready..."
for i in {1..15}; do
  if docker exec competitive_analysis_redis redis-cli ping > /dev/null 2>&1; then
    echo "    Redis is ready"
    break
  fi
  if [ $i -eq 15 ]; then
    echo "ERROR: Redis did not become ready in time"
    exit 1
  fi
  sleep 1
done

echo "==> Building Go backend..."
go build -o bin/api ./cmd/api
go build -o bin/worker ./cmd/worker

echo "==> Starting API server (port 8080)..."
mkdir -p "$ROOT/logs"
nohup ./bin/api --config configs/development.yaml > logs/api.log 2>&1 &
API_PID=$!
echo "    API PID: $API_PID"

echo "==> Starting Worker..."
nohup ./bin/worker --config configs/development.yaml > logs/worker.log 2>&1 &
WORKER_PID=$!
echo "    Worker PID: $WORKER_PID"

echo "==> Starting Frontend (port 3000)..."
cd "$ROOT/frontend"
npm run dev -- --host 0.0.0.0 --port 3000 > ../logs/frontend.log 2>&1 &
FRONTEND_PID=$!
echo "    Frontend PID: $FRONTEND_PID"
cd "$ROOT"

echo ""
echo "=== All services started ==="
echo "  API:      http://localhost:8080"
echo "  Frontend: http://localhost:3000"
echo "  PostgreSQL: localhost:54320"
echo "  Redis:    localhost:6379"
echo ""
echo "  PIDs: api=$API_PID worker=$WORKER_PID frontend=$FRONTEND_PID"
echo ""
echo "  Logs:"
echo "    bin/api logs/logs/api.log"
echo "    bin/worker logs/logs/worker.log"
echo "    logs/frontend.log"
echo ""
echo "  To stop: ./scripts/stop.sh"
