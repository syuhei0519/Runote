FROM node:22

# 作業ディレクトリを設定
WORKDIR /app

# 依存パッケージを先にインストール（キャッシュ効率化）
COPY package*.json ./
RUN npm install

# wait-port をグローバルインストール（必要に応じて npx にも切り替え可）
RUN npm install -g wait-port

# YAMLテストケースを明示的にコピー
COPY yaml/ yaml/

# その他のコード・設定ファイルをコピー
COPY . .

# テスト対象サービスが起動するのを待ってからテスト実行
CMD wait-port post-service:3000 && \
    wait-port emotion-service:8080 && \
    wait-port tag-service:4000 && \
    wait-port auth-service:8000 && \
    node --loader ts-node/esm src/index.ts