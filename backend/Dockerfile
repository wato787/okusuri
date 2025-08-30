# ビルドステージ
FROM public.ecr.aws/docker/library/golang:1.24-alpine AS builder

# ビルド依存関係のインストール
RUN apk add --no-cache git ca-certificates tzdata

# 作業ディレクトリの設定
WORKDIR /app

# Go modファイルのコピー
COPY go.mod go.sum ./

# 依存関係のダウンロード
RUN go mod download

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# 最終ステージ
FROM public.ecr.aws/docker/library/alpine:latest

# ランタイム依存関係のインストール
RUN apk --no-cache add ca-certificates tzdata

# Lambda Web Adapterのコピー
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.9.1 /lambda-adapter /opt/extensions/lambda-adapter

# 作業ディレクトリの設定
WORKDIR /app

# ビルドステージからバイナリをコピー
COPY --from=builder /app/main .

# 追加ファイルのコピー（テンプレート、静的ファイルなど）
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/pkg ./pkg

# 環境変数の設定
ENV PORT=8080
ENV GIN_MODE=release

# ポートの公開
EXPOSE 8080

# Lambda Web Adapterをエントリーポイントに設定
CMD ["/opt/extensions/lambda-adapter"]
