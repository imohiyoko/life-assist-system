# --- ステージ1: ビルド用ステージ ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

# go.modとgo.sumを先にコピーして依存関係をキャッシュ
COPY go.mod* ./
RUN go mod download

# プロジェクトのソースコードを全てコピー
COPY . .

# アプリケーションをビルド (mainパッケージのパスを正確に指定)
RUN CGO_ENABLED=0 GOOS=linux go build -o /server-binary ./cmd/api

# --- ステージ2: 実行用ステージ ---
FROM alpine:latest

# --- この部分を修正 ---
# 作業ディレクトリを /app に変更し、設定ファイルの場所と一致させる
WORKDIR /app

# ビルド用ステージから、コンパイル済みの実行ファイルのみをコピー
COPY --from=builder /server-binary .

# コンテナがリッスンするポートを8080に指定
EXPOSE 8080

# コンテナ起動時に実行するコマンド
CMD ["./server-binary"]
