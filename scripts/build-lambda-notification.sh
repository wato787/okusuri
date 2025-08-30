#!/bin/bash
set -e

echo "🔔 通知Lambda用ZIPパッケージをビルド中..."

cd notification

# ビルド用ディレクトリ作成
mkdir -p ../infra/dist/lambda

# Goバイナリビルド
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w" \
  -o ../infra/dist/lambda/notification \
  .

# ZIPパッケージ作成
cd ../infra/dist/lambda
zip -j notification.zip notification

# クリーンアップ
rm -f notification

echo "✅ 通知Lambda ZIPパッケージ完了: infra/dist/lambda/notification.zip"