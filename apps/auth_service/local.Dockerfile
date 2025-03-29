# 開発用ベースイメージ（Python 3.10）
FROM python:3.10-slim

# Poetry インストール
ENV POETRY_VERSION=1.8.2
RUN apt-get update && apt-get install -y curl make && \
    curl -sSL https://install.python-poetry.org | python3 - && \
    ln -s ~/.local/bin/poetry /usr/local/bin/poetry

# 作業ディレクトリ
WORKDIR /app

# Poetry 設定（仮想環境を作らずグローバルにインストール）
ENV POETRY_VIRTUALENVS_CREATE=false
COPY pyproject.toml ./
RUN poetry install --no-root

# ソースコードコピー
COPY . .

# Uvicorn + FastAPI 起動（ホットリロードあり）
CMD ["poetry", "run", "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]