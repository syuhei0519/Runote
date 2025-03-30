# 開発環境用 Dockerfile
FROM node:22

# 作業ディレクトリ
WORKDIR /app

# 依存ファイルコピー
COPY package*.json ./
COPY tsconfig.json ./

# ts-node-dev を入れるために先に依存だけインストール
RUN npm install

# アプリケーションコード
COPY . .

# Prisma Client を生成
RUN npx prisma generate

# 開発用のコマンド（ts-node-dev）
CMD ["npx", "ts-node-dev", "--respawn", "src/index.ts"]
