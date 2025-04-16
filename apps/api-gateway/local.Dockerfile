FROM node:22

# 作業ディレクトリ作成
WORKDIR /app

# package.json と lockファイルをコピーして依存インストール
COPY package*.json ./
RUN npm install

# ソースコードをコピー
COPY . .

# 開発用の環境変数ファイルを明示（必要に応じて）
ENV NODE_ENV=development

# ts-node-dev で開発起動
CMD ["npx", "ts-node-dev", "--respawn", "--transpile-only", "src/index.ts"]