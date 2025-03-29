# 🛠️ Runote 技術スタックまとめ（学習向け・マイクロサービス構成）

## ✅ アーキテクチャ
- マイクロサービス構成
- Kubernetes（Minikube/kind → AKS）
- GitOps（Argo CD）による自動デプロイ
- Jenkins による CI/CD パイプライン構築

---

## 📦 サービス別技術スタック

### 1. 投稿管理サービス（post-service）
- 言語: Node.js（NestJS）
- DB: PostgreSQL
- 役割: 投稿の作成・取得・編集・削除

### 2. 感情タグサービス（emotion-service）
- 言語: Go
- DB: SQLite または Redis
- 役割: 感情タグの管理（プリセット / ラベル）

### 3. タグ管理サービス（tag-service）
- 言語: Python（FastAPI）
- DB: MongoDB
- 役割: 自由タグの登録・検索・関連づけ

### 4. 認証サービス（auth-service）
- 言語: Laravel または Flask
- DB: MySQL
- 役割: ユーザー認証（JWT対応）

### 5. 画像アップロードサービス（image-service）
- 言語: Rust または Node.js
- Storage: S3互換（MinIO）
- 役割: ユーザー投稿画像の保存・取得

### 6. フロントエンド（frontend-app）
- フレームワーク: React（Next.js）または Vue 3（Vite）
- 状態管理: Zustand / Pinia / Redux Toolkit
- 機能: 投稿作成・一覧・詳細表示などのSPA

---

## ⚙️ インフラ / CI / セキュリティ

### Kubernetes
- ローカル: Minikube or kind
- 本番: Azure Kubernetes Service（AKS）

### IaC（Infrastructure as Code）
- Terraform（Azureリソース管理）
- Ansible（初期構成・設定管理）

### CI/CD
- Jenkins（CIパイプライン）
- Argo CD（GitOpsでのK8sデプロイ）

---

## 🧪 テスト・監視・品質管理

### パフォーマンステスト
- Locust（Pythonベースのユーザーシナリオテスト）

### 脆弱性スキャン
- Trivy（Dockerイメージと依存ライブラリのセキュリティチェック）
  - CI/CDと連携して脆弱性検出時にSlack通知を実装予定

### ログ・監視
- Prometheus（メトリクス収集）
- Grafana（ダッシュボード表示）
- Loki（ログ収集）

---

## 🔐 その他
- API Gateway: Traefik または NGINX
- Service Discovery: KubernetesのServiceによる名前解決
- 通知: Slack Webhook（Trivyのスキャン結果通知に使用）