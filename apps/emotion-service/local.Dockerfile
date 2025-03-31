FROM golang:1.23

WORKDIR /app

RUN go install github.com/air-verse/air@latest

# 依存ファイルを先にコピーしてキャッシュ有効化
COPY go.mod ./
COPY go.sum ./

RUN go mod download

# ソースコードを追加
COPY . .

CMD ["air"]