#!/bin/bash
set -e

echo "🐳 API Lambda用Dockerイメージをビルド中..."

cd backend
docker build -t okusuri-backend:latest .

echo "✅ Dockerイメージビルド完了: okusuri-backend:latest"