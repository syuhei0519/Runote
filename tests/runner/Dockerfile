FROM node:22

# 作業ディレクトリを設定
WORKDIR /app

# package.jsonとpackage-lock.jsonを先にコピーして依存関係だけ先にインストール（キャッシュ活用）
COPY package*.json ./
RUN npm install

# 残りのファイルをコピー
COPY . .

# エントリーポイントを指定（docker-composeからはCMDで上書き可能）
CMD ["node", "--loader", "ts-node/esm", "index.ts"]