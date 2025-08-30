# Okusuri - お薬管理アプリケーション

このリポジトリは、お薬管理アプリケーション「Okusuri」のモノレポです。

## プロジェクト構成

- **`okusuri-backend`** - Go 製のバックエンド API
- **`okusuri-frontend`** - Next.js 製のフロントエンド（Next.js 15）
- **`okusuri-v2`** - Vite + React 製のフロントエンド（V2）

## 開発環境のセットアップ

### 前提条件

- Node.js 18 以上
- pnpm 8 以上
- Go 1.24 以上
- PostgreSQL

### インストール

```bash
# 依存関係のインストール
pnpm install

# 全プロジェクトの開発サーバーを起動
pnpm dev

# 個別のプロジェクトを起動
pnpm dev:backend    # バックエンドのみ
pnpm dev:frontend   # フロントエンドのみ
pnpm dev:v2         # V2のみ
```

### ビルド

```bash
# 全プロジェクトをビルド
pnpm build

# 個別のプロジェクトをビルド
pnpm build:backend
pnpm build:frontend
pnpm build:v2
```

### リント・フォーマット

```bash
# 全プロジェクトのリント
pnpm lint

# 全プロジェクトのリント修正
pnpm lint:fix
```

## 開発ワークフロー

1. 各アプリケーションは独立したディレクトリで開発
2. ルートの `package.json` でワークスペース全体の管理
3. `pnpm --filter` コマンドで個別プロジェクトの操作

## 技術スタック

### バックエンド

- Go 1.24
- Gin (Web フレームワーク)
- GORM (ORM)
- PostgreSQL

### フロントエンド

- Next.js 15 (okusuri-frontend)
- Vite + React (okusuri-v2)
- TypeScript
- Tailwind CSS
- Radix UI

## ライセンス

Private
