# Okusuri V2

服薬記録と管理のためのアプリケーション（Vite + React + TypeScript 版）

## セットアップ

### 1. 依存関係のインストール

```bash
pnpm install
```

### 2. 環境変数の設定

```bash
# 環境変数ファイルをコピー
cp env.example .env

# VAPIDキーを生成（Web Push通知用）
pnpm run generate-vapid
```

生成された VAPID キーを`.env`ファイルに設定してください：

```env
VITE_API_URL=http://localhost:8080
VITE_VAPID_PUBLIC_KEY=生成された公開キー
VITE_GOOGLE_CLIENT_ID=your_google_client_id_here
```

### 3. 開発サーバーの起動

```bash
pnpm dev
```

## 利用可能なスクリプト

- `pnpm dev` - 開発サーバーを起動
- `pnpm build` - プロダクションビルド
- `pnpm preview` - ビルド結果をプレビュー
- `pnpm lint` - リンターを実行
- `pnpm lint:fix` - リンターを実行して自動修正
- `pnpm format` - コードフォーマット
- `pnpm check` - コードチェック
- `pnpm generate-vapid` - VAPID キーを生成

## 技術スタック

- **フレームワーク**: Vite + React + TypeScript
- **UI**: Tailwind CSS + Radix UI
- **状態管理**: TanStack Query
- **認証**: Better Auth
- **ルーティング**: React Router DOM
- **通知**: Web Push API
- **PWA**: Vite PWA Plugin

## 機能

- 服薬記録の管理
- カレンダー表示
- 通知設定
- PWA 対応
- レスポンシブデザイン

## 開発

### ディレクトリ構造

```
src/
├── components/     # UIコンポーネント
├── contexts/       # React Context
├── hooks/          # カスタムフック
├── lib/            # ユーティリティライブラリ
├── pages/          # ページコンポーネント
├── services/       # APIサービス
├── types/          # TypeScript型定義
└── utils/          # ユーティリティ関数
```

### 環境変数

| 変数名                  | 説明                             | 必須 |
| ----------------------- | -------------------------------- | ---- |
| `VITE_API_URL`          | バックエンド API の URL          | ✅   |
| `VITE_VAPID_PUBLIC_KEY` | Web Push 通知用の VAPID 公開キー | ✅   |
| `VITE_GOOGLE_CLIENT_ID` | Google OAuth クライアント ID     | ✅   |

## ライセンス

MIT
