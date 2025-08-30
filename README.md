# Okusuri - お薬管理アプリケーション

このリポジトリは、お薬管理アプリケーション「Okusuri」のプロジェクトです。

## プロジェクト構成

- **`backend`** - Go 製のバックエンド API
- **`frontend`** - Next.js 製のフロントエンド（Next.js 15）
- **`okusuri-v2`** - Vite + React 製のフロントエンド（V2）
- **`notification`** - 通知関連のドキュメント
- **`infra`** - インフラ関連のドキュメント
- **`docs`** - 設計・Terraformドキュメント

## 開発環境のセットアップ

### 前提条件

- Node.js 18 以上
- pnpm 8 以上
- Go 1.24 以上
- PostgreSQL

### インストール

```bash
# フロントエンドの依存関係をインストール
cd frontend && pnpm install
cd ../okusuri-v2 && pnpm install

# バックエンドの依存関係をインストール
cd ../backend && go mod download
```

### 開発サーバーの起動

```bash
# フロントエンド（Next.js）
pnpm dev:frontend

# V2（Vite + React）
pnpm dev:v2

# バックエンド（Go）
cd backend && go run cmd/server/main.go
# または
cd backend && make dev
```

### ビルド

```bash
# フロントエンド
pnpm build:frontend

# V2
pnpm build:v2

# バックエンド
cd backend && go build cmd/server/main.go
```

### リント・フォーマット

```bash
# フロントエンド
pnpm lint:frontend
pnpm lint:fix:frontend

# V2
pnpm lint:v2
pnpm lint:fix:v2

# バックエンド
cd backend && golangci-lint run
```

## 開発ワークフロー

1. 各プロジェクトは独立したディレクトリで開発
2. ルートの `package.json` でフロントエンドプロジェクトの操作を簡素化
3. バックエンドはGoの標準的な開発フローを使用

## 技術スタック

### バックエンド

- Go 1.24
- Gin (Web フレームワーク)
- GORM (ORM)
- PostgreSQL

### フロントエンド

- Next.js 15 (frontend)
- Vite + React (okusuri-v2)
- TypeScript
- Tailwind CSS
- Radix UI

## ライセンス

Private
