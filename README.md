# Okusuri - お薬管理アプリケーション

このリポジトリは、お薬管理アプリケーション「Okusuri」のプロジェクトです。

## プロジェクト構成

- **`backend`** - Go 製のバックエンド API
- **`frontend`** - Next.js 製のフロントエンド（Next.js 15）
- **`okusuri-v2`** - Vite + React 製のフロントエンド（V2）
- **`notification`** - 通知関連のドキュメント
- **`infra`** - インフラ関連のドキュメント
- **`docs`** - 設計・Terraformドキュメント

## 前提条件

- Node.js 18 以上
- pnpm 8 以上
- Go 1.24 以上
- PostgreSQL
- [Task](https://taskfile.dev/) (プロジェクト管理ツール)

## セットアップ

### Taskのインストール

```bash
# macOS (Homebrew)
brew install go-task

# その他のプラットフォーム
# https://taskfile.dev/installation/
```

### 依存関係のインストール

```bash
# 全プロジェクトの依存関係をインストール
task install:all

# 個別にインストール
task install:frontend    # Next.js
task install:v2          # Vite + React
```

## 開発

### 開発サーバーの起動

```bash
# フロントエンド（Next.js）
task dev:frontend

# V2（Vite + React）
task dev:v2

# バックエンド（Go）
task dev:backend
```

### ビルド

```bash
# フロントエンド
task build:frontend

# V2
task build:v2

# バックエンド
task build:backend
```

### リント・フォーマット

```bash
# フロントエンド
task lint:frontend
task lint:fix:frontend

# V2
task lint:v2
task lint:fix:v2
```

### テスト

```bash
# バックエンド
task test:backend
```

## その他のコマンド

```bash
# プロジェクトの状態確認
task status

# ビルド成果物のクリーンアップ
task clean

# 本番サーバー起動
task start:frontend

# Viteプレビュー
task preview:v2
```

## 利用可能なタスク

```bash
# 全タスクを表示
task --list

# タスクの詳細を表示
task --list-all
```

## 開発ワークフロー

1. 各プロジェクトは独立したディレクトリで開発
2. Taskfileでプロジェクト全体の操作を統一
3. 各プロジェクトの依存関係は個別に管理

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
