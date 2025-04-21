# 開発用 Dockerfile

FROM node:20-slim

# 作業ディレクトリ作成
WORKDIR /app

# package.json / lock ファイルをコピー
COPY package*.json ./

# パッケージインストール
RUN npm install

# ソースコードコピー
COPY . .

# ポート指定
EXPOSE 5173

# devサーバ起動（ホットリロード付き）
CMD ["npm", "run", "dev"]