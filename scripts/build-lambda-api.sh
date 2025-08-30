#!/bin/bash
set -e

echo "🐳 API Lambda用Dockerイメージをビルド中..."

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

# Dockerイメージをビルド
cd backend
docker build -t okusuri-api:latest .

echo "✅ Dockerイメージビルド完了: okusuri-api:latest"
echo "📝 次のステップ: task deploy:backend でECRにプッシュ"