FROM golang:1.22

WORKDIR /app

# 依存ファイルを先にコピーしてキャッシュ有効化
COPY go.mod ./
COPY go.sum ./

RUN go mod download

# ソースコードを追加
COPY . .

CMD ["go", "run", "main.go"]